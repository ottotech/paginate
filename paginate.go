/*
Package paginate provides a basic Paginator interface to do pagination of database records.
Its primary job is to generate a raw sql command with the corresponding values
which can be executed in a database.

This package will also handle basic filtering of records by creating raw sql commands
based on the ulr.Values and options coming from a request. And –if used correctly-
Paginator can also return a PaginationResponse which contains useful information
for clients to do proper pagination.

In order to do pagination with Paginator the url.Values should contain the
`page` parameter, if not the PageSize const or pageSize parameter specified by
the client of package will be used.

For ordering records based on column names use the following syntax in the url:

	localhost/some-url?name=otto&sort=+name,-age

	Note: use the sign + for ascending ordering or the sign - for descending ordering
*/
package paginate

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
)

// When initializing a Paginator with NewPaginator this would be the default page size for the pagination.
const (
	PageSize = 30
)

const (
	gt  = ">"
	lt  = "<"
	gte = ">="
	lte = "<="
)

// Paginator is the interface that wraps pagination behaviors.
//
// Paginate will use the data consumed by NewPaginator or NewPaginatorWithLimit
// and it will return an sql command with its corresponding values. The operation
// will run concurrently to improve computation speed.
//
// SetPageCount can be used to set the PageCount of the PaginationResponse. If the
// client of the package pretends to use the PaginationResponse information for
// further pagination it is necessary that he/she sets this variable before calling Response.
// When should SetPageCount be call? Just after executing the sql command generated by Paginate.
// Only after executing the sql command it is possible to know the number of records retrieved from
// the database, in that case that would be the value of the PageCount.
//
// SetTotalResult works similar to its counterpart SetPageCount, but in this case SetTotalResult
// will set the TotalSize of records matching the url.Values parameters and ordering, without
// taking into account the PageSize of Paginator.
// The `count(*) over()` part of the returning sql will serve as the field we can use to set
// the TotalSize value.
// When should SetPageCount be call? Just after executing the sql command generated by Paginate.
//
// Response should be executed after calling SetPageCount and SetTotalResult.
// Response will return a PaginationResponse struct containing useful information for clients
// of the package to do proper and subsequent pagination operations.
//
//
type Paginator interface {
	Paginate() (sql string, args []interface{}, err error)
	Response() PaginationResponse
	SetPageCount(size int)
	SetTotalResult(size int)
}

type pagination struct {
	request    paginationRequest
	response   PaginationResponse
	pageSize   int
	tableName  string
	parameters url.Values
	filters    []filter
	colNames   []string
}

type paginationRequest struct {
	pageNumber int
}

type PaginationResponse struct {
	PageNumber      int  `json:"page_number"`
	NextPageNumber  int  `json:"next_page_number"`
	HasNextPage     bool `json:"has_next_page"`
	HasPreviousPage bool `json:"has_previous_page"`
	PageCount       int  `json:"page_count"`
	TotalSize       int  `json:"total_size"`
}

type whereClause struct {
	clause string
	args   []interface{}
	exists bool
}

type filterClause struct {
	clause string
	args   []interface{}
	exists bool
}

type filter struct {
	field string
	sign  string
	value string
}

func NewPaginator(tableName string, colNames []string, u url.URL) Paginator {
	var paginator Paginator
	c := make(chan []filter)
	decodedURL, err := url.QueryUnescape(u.String())
	if err != nil {
		log.Println(err)
	}
	go getFilters(decodedURL, colNames, c)
	v := u.Query()
	p := new(pagination)
	p.tableName = tableName
	p.colNames = colNames
	p.parameters = v
	p.pageSize = PageSize
	p.filters = <-c
	p.request = getRequestData(v)
	paginator = p
	return paginator
}

func NewPaginatorWithLimit(pageSize int, tableName string, colNames []string, u url.URL) Paginator {
	var paginator Paginator
	c := make(chan []filter)
	decodedURL, err := url.QueryUnescape(u.String())
	if err != nil {
		log.Println(err)
	}
	go getFilters(decodedURL, colNames, c)
	v := u.Query()
	p := new(pagination)
	p.tableName = tableName
	p.colNames = colNames
	p.parameters = v
	p.pageSize = pageSize
	p.filters = <-c
	p.request = getRequestData(v)
	paginator = p
	return paginator
}

func (p *pagination) Paginate() (sql string, values []interface{}, err error) {
	var s string
	var AND = " AND "
	var WHERE = " WHERE"
	c1 := make(chan whereClause)
	c2 := make(chan string)
	c3 := make(chan string)
	c4 := make(chan filterClause)
	go createWhereClause(p.colNames, p.parameters, c1)
	go createPaginationClause(p.request.pageNumber, p.pageSize, c2)
	go createOrderByClause(p.parameters, p.colNames, c3)
	go createFilterClause(p.filters, c4)
	where := <-c1
	pagination := <-c2
	order := <-c3
	filter := <-c4

	numArgs := len(where.args) + len(filter.args)
	placeholders := make([]interface{}, 0)
	for i := 1; i < numArgs+1; i++ {
		placeholders = append(placeholders, i)
	}
	args := make([]interface{}, 0)
	args = append(args, where.args...)
	args = append(args, filter.args...)

	if where.exists && filter.exists {
		s = "SELECT " + strings.Join(p.colNames, ", ") + ", count(*) over() FROM " + p.tableName + where.clause + AND + filter.clause + order + pagination
		s = fmt.Sprintf(s, placeholders...)
	} else if where.exists {
		s = "SELECT " + strings.Join(p.colNames, ", ") + ", count(*) over() FROM " + p.tableName + where.clause + order + pagination
		s = fmt.Sprintf(s, placeholders...)
	} else if filter.exists {
		s = "SELECT " + strings.Join(p.colNames, ", ") + ", count(*) over() FROM " + p.tableName + WHERE + filter.clause + order + pagination
		s = fmt.Sprintf(s, placeholders...)
	} else {
		s = "SELECT " + strings.Join(p.colNames, ", ") + ", count(*) over() FROM " + p.tableName + order + pagination
	}
	return s, args, nil
}

func (p *pagination) Response() PaginationResponse {
	p.response.PageNumber = p.request.pageNumber

	if (p.request.pageNumber * p.pageSize) < p.response.TotalSize {
		p.response.NextPageNumber = p.request.pageNumber + 1
		p.response.HasNextPage = true
	}
	if (p.request.pageNumber * p.pageSize) == p.response.TotalSize {
		p.response.NextPageNumber = 0
		p.response.HasNextPage = false
	}
	if p.response.TotalSize == 0 {
		p.response.NextPageNumber = 0
		p.response.HasNextPage = false
	}

	if p.response.PageNumber > 1 {
		p.response.HasPreviousPage = true
	}
	return p.response
}

func (p *pagination) SetTotalResult(size int) {
	p.response.TotalSize = size
}

func (p *pagination) SetPageCount(count int) {
	p.response.PageCount = count
}

func getRequestData(v url.Values) paginationRequest {
	p := paginationRequest{}
	if page := v.Get("page"); page != "" {
		page, err := strconv.Atoi(page)
		if err != nil {
			page = 1
		}
		p.pageNumber = page
	} else {
		p.pageNumber = 1
	}
	return p
}

func createWhereClause(colNames []string, v url.Values, c chan whereClause) {
	w := whereClause{}
	var WHERE = " WHERE "
	var AND = " AND "
	var separator string
	var clauses []string
	var values []interface{}

	// map all db column names with the url parameters
	for _, name := range colNames {
		if val := v.Get(name); val != "" {
			values = append(values, val)
			clauses = append(clauses, name+" = $%v")
		}
	}
	// use appropriate `separator` to join the clauses
	if len(clauses) == 1 {
		separator = ""
	} else {
		separator = AND
	}
	// let's map the clause and args to the whereClause struct, and specify if there were some where clauses at all
	w.clause = WHERE + strings.Join(clauses, separator)
	w.args = values
	w.exists = len(clauses) > 0
	c <- w
}

func createFilterClause(filters []filter, c chan filterClause) {
	var AND = " AND "
	var s = ""
	f := filterClause{}
	args := make([]interface{}, 0)
	for i, f := range filters {
		if v, err := strconv.Atoi(f.value); err != nil {
			args = append(args, v)
		} else {
			args = append(args, f.value)
		}

		if i == 0 {
			s += fmt.Sprintf(" %s %s ", f.field, f.sign) + "$%v"
			continue
		}
		if i == len(filters)-1 {
			s += fmt.Sprintf(" %s %s ", f.field, f.sign) + "$%v"
			continue
		}
		s += fmt.Sprintf(" %s %s ", f.field, f.sign) + "$%v" + AND
	}
	f.clause = s
	f.args = args
	f.exists = s != ""
	c <- f
}

func createPaginationClause(pageNumber int, pageSize int, c chan string) {
	var clause string
	var offset int

	if pageSize > PageSize {
		pageSize = PageSize
	} else if pageSize < 0 {
		pageSize = PageSize
	}

	clause += fmt.Sprintf(" LIMIT %v ", pageSize)

	if pageNumber < 0 || pageNumber == 0 || pageNumber == 1 {
		offset = 0
	} else {
		offset = pageSize * (pageNumber - 1)
	}

	clause += fmt.Sprintf(" OFFSET %v", offset)
	c <- clause
}

func createOrderByClause(v url.Values, colNames []string, c chan string) {
	var ASC = "ASC"
	var DESC = "DESC"
	clauses := make([]string, 0)
	sort := v.Get("sort")
	if sort == "" {
		c <- " ORDER BY id "
		return
	}
	fields := strings.Split(sort, ",")
	for _, v := range fields {
		orderBy := string(v[0])
		field := string(v[1:])
		for _, f := range colNames {
			if f == "id" {
				continue
			}
			if field == f {
				if orderBy == " " {
					clauses = append(clauses, field+" "+ASC)
				}
				if orderBy == "-" {
					clauses = append(clauses, field+" "+DESC)
				}
			}
		}
	}
	clauses = append(clauses, "id")
	clauseSTR := strings.Join(clauses, ",")
	c <- " ORDER BY " + clauseSTR
}

func getFilters(decodedPath string, colNames []string, c chan []filter) {
	s := make([]filter, 0)
	i := strings.Index(decodedPath, "?")
	if i == -1 {
		c <- s
		return
	}

	getF := func(key, val, char string) (bool, filter) {
		f := filter{}
		if strings.Contains(val, char) {
			if val[:len(char)] == char && len(val) > len(char) {
				f.field = key
				f.sign = char
				f.value = val[len(char):]
				return true, f
			}
		}
		return false, f
	}

	params := strings.Split(decodedPath[i+1:], "&")
	for _, n := range colNames {
		for _, p := range params {

			if len(p) <= len(n) {
				continue
			}
			key, value := p[:len(n)], p[len(n):]
			if key != n {
				continue
			}
			if ok, f := getF(key, value, gte); ok {
				s = append(s, f)
				continue
			}
			if ok, f := getF(key, value, lte); ok {
				s = append(s, f)
				continue
			}
			if ok, f := getF(key, value, gt); ok {
				s = append(s, f)
				continue
			}
			if ok, f := getF(key, value, lt); ok {
				s = append(s, f)
				continue
			}
		}
	}
	c <- s
}
