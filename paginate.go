package paginate

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

const (
	PageSize = 30
)

type Paginator interface {
	Paginate() (sql string, args []interface{}, err error)
	Response() PaginationResponse
	SetTotalResult(size int)
}

type pagination struct {
	request    paginationRequest
	response   PaginationResponse
	tableName  string
	parameters url.Values
	dbColNames []string
}

type paginationRequest struct {
	pageNumber int
	pageSize   int
}

type PaginationResponse struct {
	NextPageNumber int  `json:"next_page_number"`
	HasNextPage    bool `json:"has_next_page"`
	TotalSize      int  `json:"total_size"`
}

type whereClause struct {
	clause string
	args   []interface{}
}

func NewPaginator(tableName string, dbColNames []string, params url.Values) Paginator {
	var paginator Paginator
	p := new(pagination)
	p.tableName = tableName
	p.dbColNames = dbColNames
	p.parameters = params
	p.request = getRequestData(params)
	paginator = p
	return paginator
}

func (p *pagination) Paginate() (sql string, values []interface{}, err error) {
	c1 := make(chan whereClause)
	c2 := make(chan string)
	c3 := make(chan string)
	go createWhereClause(p.dbColNames, p.parameters, c1)
	go createPaginationClause(p.request, c2)
	go createOrderByClause(p.parameters, p.dbColNames, c3)
	where := <-c1
	pagination := <-c2
	order := <-c3
	s := "SELECT " + strings.Join(p.dbColNames, ", ") + ", count(*) over() FROM " + p.tableName + where.clause + order + pagination
	return s, where.args, nil
}

func (p *pagination) Response() PaginationResponse {
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
	return p.response
}

func (p *pagination) SetTotalResult(size int) {
	p.response.TotalSize = size
}

func getRequestData(v url.Values) paginationRequest {
	p := paginationRequest{}
	if page := v.Get("page"); page != "" {
		page, err := strconv.Atoi(page)
		if err != nil {
			page = 1
		}
		p.pageNumber = page
	}
	if pageSize := v.Get("page_size"); pageSize != "" {
		pageSize, err := strconv.Atoi(pageSize)
		if err != nil {
			pageSize = 0
		}
		p.pageSize = pageSize
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
	var placeHolder = 1

	// map all db column names with the url parameters
	for _, name := range colNames {
		if val := v.Get(name); val != "" {
			values = append(values, val)
			clauses = append(clauses, fmt.Sprintf("%s = $%v", name, placeHolder))
			placeHolder++
		}
	}
	// use appropriate `separator` to join the clauses
	if len(clauses) == 1 {
		separator = ""
	} else {
		separator = AND
	}
	// let's map the clause and args to the whereClause struct
	w.clause = WHERE + strings.Join(clauses, separator)
	w.args = values
	// sends to channel
	c <- w
}

func createPaginationClause(p paginationRequest, c chan string) {
	var clause string
	var pageSize int
	var offset int

	if p.pageSize > 0 {
		pageSize = p.pageSize
	} else if p.pageSize > PageSize {
		pageSize = PageSize
	} else {
		pageSize = PageSize
	}

	clause += fmt.Sprintf(" LIMIT %v ", pageSize)

	if p.pageNumber < 0 || p.pageNumber == 0 || p.pageNumber == 1 {
		offset = 0
	} else {
		offset = pageSize * (p.pageNumber - 1)
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
