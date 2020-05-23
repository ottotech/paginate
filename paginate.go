/*
Package paginate provides a basic Paginator interface to do pagination of database records.
Its primary job is to generate a raw sql command with the corresponding values
which can be executed in a database.

This package will also handle basic filtering of records by creating an sql command
based on the parameters coming from a request. And –if used correctly- Paginator can
also return a PaginationResponse which contains useful information for clients to do
proper pagination.

For ordering records based on column names use the following syntax in the url (-, +):

	localhost/some-url?name=otto&sort=+name,-age

For normal filtering the following operators are available:
	eq  = "="
	gt  = ">"
	lt  = "<"
	gte = ">="
	lte = "<="
	ne  = "<>"

*/
package paginate

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Constant that specifies the page size of the pagination results.
// This value will be used in cases where the page size cannot be determined.
// For example, when using NewPaginator, if NewPaginator fails to retrieve
// the page size from the request url it will fallback to this value.
const (
	PageSize = 30
)

// Constants that specify the available filter operators.
// These operators can be used in the request url to filter records.
const (
	eq  = "="
	gt  = ">"
	lt  = "<"
	gte = ">="
	lte = "<="
	ne  = "<>"
)

// Paginator is the interface that wraps pagination behaviors.
type Paginator interface {
	// Paginate will use the data consumed by NewPaginator or NewPaginatorWithLimit,
	// and it will return an sql command with its corresponding values. The operation
	// will run concurrently to improve computation speed.
	//
	// The sql command generated by Paginate should be run against the target database.
	Paginate() (sql string, args []interface{}, err error)

	// SetPageCount can be used to set the PageCount of the PaginationResponse. If the
	// client of the package pretends to use the PaginationResponse information for
	// further pagination it is necessary that he/she sets this variable before calling Response.
	// When should SetPageCount be called? Just after executing the sql command generated by Paginate.
	// Only after executing the sql command it is possible to know the number of records retrieved from
	// the database, in that case that would be the value of the PageCount.
	SetPageCount(size int)

	// SetTotalResult works similar to its counterpart SetPageCount, but in this case SetTotalResult
	// will set the TotalSize of records. The `count(*) over()` part of the returning sql will serve as
	// the field we can use to set the TotalSize value.
	// When should SetTotalResult be called? Just after executing the sql command generated by Paginate.
	SetTotalResult(size int)

	// Response should be executed after calling SetPageCount and SetTotalResult.
	// Response will return a PaginationResponse struct containing useful information for clients
	// of the package so they can do proper and subsequent pagination operations.
	Response() PaginationResponse
}

// pagination is a concrete type that implements the Paginator interface.
type pagination struct {
	request    paginationRequest
	response   PaginationResponse
	tableName  string
	colNames   []string
	parameters parameters
}

// paginationRequest holds information about the pagination that the client
// wants to execute. The helper function getRequestData() helps this package to
// get this information from the url.Values of an http.Request.
type paginationRequest struct {
	pageNumber int
	pageSize   int
}

// PaginationResponse contains information about the pagination.
//
// Clients of this package can use PaginationResponse to paginate further their data.
type PaginationResponse struct {
	PageNumber      int  `json:"page_number"`
	NextPageNumber  int  `json:"next_page_number"`
	HasNextPage     bool `json:"has_next_page"`
	HasPreviousPage bool `json:"has_previous_page"`
	PageCount       int  `json:"page_count"`
	TotalSize       int  `json:"total_size"`
}

// whereClause holds information about a where clauses.
type whereClause struct {
	clause string
	args   []interface{}
	exists bool
}

type parameters []parameter

// getParameter will try to get a parameter by its name from parameters.
func (params *parameters) getParameter(name string) (parameter, bool) {
	for _, p := range *params {
		if p.name == name {
			return p, true
		}
	}
	return parameter{}, false
}

// parameter holds information about a parameter passed in the url.Values from a http.Request.
type parameter struct {
	name  string
	sign  string
	value string
}

// NewPaginator creates a Paginator object ready to paginate data from a database table.
//
// NewPaginator will try to get the page size for the paginator result from the request url
// If it fails to get the parameter from the request url, it will use the constant PageSize.
func NewPaginator(tableName string, colNames []string, u url.URL) Paginator {
	c := make(chan parameters)
	go getParameters(colNames, u, c)
	v := u.Query()
	p := new(pagination)
	p.tableName = tableName
	p.colNames = colNames
	p.request = getRequestData(v)
	p.parameters = <-c
	return p
}

// NewPaginatorWithLimit creates a Paginator object ready to paginate data from a database table.
//
// NewPaginatorWithLimit specifies explicitly the page size we want to use for the pagination results.
func NewPaginatorWithLimit(pageSize int, tableName string, colNames []string, u url.URL) Paginator {
	c := make(chan parameters)
	go getParameters(colNames, u, c)
	v := u.Query()
	p := new(pagination)
	p.tableName = tableName
	p.colNames = colNames
	p.request = getRequestData(v)
	if pageSize <= 0 {
		pageSize = PageSize
	}
	p.request.pageSize = pageSize // here we override the pageSize
	p.parameters = <-c
	return p
}

func (p *pagination) Paginate() (sql string, values []interface{}, err error) {
	var s string
	c1 := make(chan whereClause)
	c3 := make(chan string)
	c4 := make(chan string)
	go createWhereClause(p.colNames, p.parameters, c1)
	go createPaginationClause(p.request.pageNumber, p.request.pageSize, c3)
	go createOrderByClause(p.parameters, p.colNames, c4)
	where := <-c1
	pagination := <-c3
	order := <-c4

	numArgs := len(where.args)
	placeholders := make([]interface{}, 0)
	for i := 1; i < numArgs+1; i++ {
		placeholders = append(placeholders, i)
	}

	if where.exists {
		s = "SELECT " + strings.Join(p.colNames, ", ") + ", count(*) over() FROM " + p.tableName + where.clause + order + pagination
		s = fmt.Sprintf(s, placeholders...)
	} else {
		s = "SELECT " + strings.Join(p.colNames, ", ") + ", count(*) over() FROM " + p.tableName + order + pagination
	}
	return s, where.args, nil
}

func (p *pagination) Response() PaginationResponse {
	p.response.PageNumber = p.request.pageNumber

	if (p.request.pageNumber * p.request.pageSize) < p.response.TotalSize {
		p.response.NextPageNumber = p.request.pageNumber + 1
		p.response.HasNextPage = true
	}
	if (p.request.pageNumber * p.request.pageSize) == p.response.TotalSize {
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

func getParameters(colNames []string, u url.URL, c chan parameters) {
	decodedURL, _ := url.PathUnescape(u.String())
	list := make(parameters, 0)
	i := strings.Index(decodedURL, "?")
	if i == -1 {
		c <- list
		return
	}

	getP := func(key, val, char string) (bool, parameter) {
		p := parameter{}
		if strings.Contains(val, char) {
			if val[:len(char)] == char && len(val) > len(char) {
				p.name = key
				p.sign = char
				p.value = val[len(char):]
				return true, p
			}
		}
		return false, p
	}
	params := strings.Split(decodedURL[i+1:], "&")
	for _, n := range colNames {
		for _, p := range params {
			if len(p) <= len(n) {
				continue
			}
			key, value := p[:len(n)], p[len(n):]
			if key != n {
				continue
			}
			// order matters
			if ok, newP := getP(key, value, gte); ok {
				list = append(list, newP)
				continue
			}
			if ok, newP := getP(key, value, lte); ok {
				list = append(list, newP)
				continue
			}
			if ok, newP := getP(key, value, ne); ok {
				list = append(list, newP)
				continue
			}
			if ok, newP := getP(key, value, gt); ok {
				list = append(list, newP)
				continue
			}
			if ok, newP := getP(key, value, lt); ok {
				list = append(list, newP)
				continue
			}
			if ok, newP := getP(key, value, eq); ok {
				list = append(list, newP)
				continue
			}
		}
	}
	// as an special case we need to also get our custom sort parameter
	sort := "sort"
	for _, p := range params {
		if len(p) <= len(sort) {
			continue
		}
		key, value := p[:len(sort)], p[len(sort):]
		if key != sort {
			continue
		}
		if ok, newP := getP(key, value, eq); ok {
			list = append(list, newP)
			continue
		}
	}
	c <- list
}

func getRequestData(v url.Values) paginationRequest {
	p := paginationRequest{}
	if page := v.Get("page"); page != "" {
		page, err := strconv.Atoi(page)
		if err != nil {
			page = 1
		}
		if page <= 0 {
			page = 1
		}
		p.pageNumber = page
	} else {
		p.pageNumber = 1
	}

	if pageSize := v.Get("page_size"); pageSize != "" {
		pageSize, err := strconv.Atoi(pageSize)
		if err != nil {
			pageSize = PageSize
		}
		if pageSize <= 0 {
			pageSize = PageSize
		}
		p.pageSize = pageSize
	} else {
		p.pageSize = PageSize
	}
	return p
}

func createWhereClause(colNames []string, params parameters, c chan whereClause) {
	w := whereClause{}
	var WHERE = " WHERE "
	var AND = " AND "
	var separator string
	var clauses []string
	var values []interface{}

	// map all db column names with the url parameters
	for _, name := range colNames {
		for _, p := range params {
			if p.name == name {
				values = append(values, p.value)
				clauses = append(clauses, p.name+" "+p.sign+" $%v")
			}
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

	clause += fmt.Sprintf("OFFSET %v", offset)
	c <- clause
}

func createOrderByClause(params parameters, colNames []string, c chan string) {
	var ASC = "ASC"
	var DESC = "DESC"
	clauses := make([]string, 0)
	sort, exists := params.getParameter("sort")
	if !exists {
		c <- " ORDER BY id"
		return
	}
	fields := strings.Split(sort.value, ",")
	for _, v := range fields {
		orderBy := string(v[0])
		field := v[1:]
		for _, f := range colNames {
			if f == "id" {
				// we will always order the records by ID (see below). In order
				// to keep the same order between pages or results
				continue
			}
			if field == f {
				if orderBy == "+" {
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
