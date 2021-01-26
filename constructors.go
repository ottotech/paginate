package paginate

import (
	"net/url"
	"reflect"
	"strings"
)

// NewPaginator creates a Paginator object ready to paginate data from a database table.
//
// NewPaginator will try to get the page size for the paginator result from the request url
// If it fails to get the parameter from the request url, it will use the constant PageSize.
func NewPaginator(table interface{}, tableName string, u url.URL) (Paginator, error) {
	p := &paginator{table: table, rv: reflect.ValueOf(table)}

	// Order matters. Validation should happen before getting
	// all the data to initialize the Paginator.
	if err := p.validateTable(); err != nil {
		return p, err
	}
	p.getCols()
	p.getFieldNames()
	p.getFilters()
	p.getID()
	if tableName == "" {
		p.getTableName()
	} else {
		p.name = strings.TrimSpace(tableName)
	}

	c := make(chan parameters)
	go getParameters(p.cols, u, c)
	v := u.Query()
	p.request = getRequestData(v)
	p.parameters = <-c
	return p, nil
}

// NewPaginatorWithLimit creates a Paginator object ready to paginate data from a database table.
//
// NewPaginatorWithLimit specifies explicitly the page size we want to use for the pagination results.
func NewPaginatorWithLimit(pageSize int, table interface{}, tableName string, u url.URL) (Paginator, error) {
	p := &paginator{table: table, rv: reflect.ValueOf(table)}

	// Order matters. Validation should happen before getting
	// all the data to initialize the Paginator.
	if err := p.validateTable(); err != nil {
		return p, err
	}
	p.getCols()
	p.getFieldNames()
	p.getFilters()
	p.getID()
	if tableName == "" {
		p.getTableName()
	} else {
		p.name = strings.TrimSpace(tableName)
	}
	if pageSize <= 0 {
		pageSize = PageSize
	}
	c := make(chan parameters)
	go getParameters(p.cols, u, c)
	v := u.Query()
	p.request = getRequestData(v)
	p.request.pageSize = pageSize // here we override the pageSize
	p.parameters = <-c
	return p, nil
}
