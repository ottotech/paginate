package paginate

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

type Option func(p *paginator) error

// TableName is an option for NewPaginator which indicates
// the name of the table that we want to paginate in the
// database. Use this option if the name of your database
// table cannot be inferred by the given struct table argument.
func TableName(name string) Option {
	return func(p *paginator) error {
		name = strings.TrimSpace(name)
		if name == "" {
			return fmt.Errorf("table name should not be an empty string")
		}
		p.name = name
		return nil
	}
}

// PageSize is an option for NewPaginator which indicates
// the size in number of records that we want our paginator
// to produce per page. ``size`` should be an uint value
// greater than zero. Use this option if you want finer
// control on the pagination size. Using this option will
// override the ``page_size`` parameter coming from the
// request in the url.URL.
func PageSize(size uint) Option {
	return func(p *paginator) error {
		if size == 0 {
			return fmt.Errorf("page size should be an uint value greater than zero")
		}
		p.pageSize = int(size)
		return nil
	}
}

// NewPaginator creates a Paginator object ready to paginate data from a database table.
// See TableName and PageSize options for NewPaginator.
func NewPaginator(table interface{}, u url.URL, opts ...Option) (Paginator, error) {
	p := &paginator{table: table, rv: reflect.ValueOf(table)}

	// Let's try to set the options if any.
	if opts != nil {
		for _, opt := range opts {
			err := opt(p)
			if err != nil {
				return nil, err
			}
		}
	}

	// If the table name has not been defined, let's try to
	// infer the name from the given table struct.
	if p.name == "" {
		p.getTableName()
	}

	v := u.Query()
	requestParameters := getRequestData(v)

	// Let's try to set the pageSize if it has not been set yet.
	// We will try to get this value from the request.
	if p.pageSize == 0 {
		p.pageSize = requestParameters.pageSize
	}

	p.pageNumber = requestParameters.pageNumber

	// Order matters. Validation should happen before getting
	// all the data to initialize the Paginator.
	if err := p.validateTable(); err != nil {
		return p, err
	}

	p.getCols()
	p.getFieldNames()
	p.getFilters()
	p.getID()

	c := make(chan parameters)
	go getParameters(p.cols, u, c)
	p.parameters = <-c
	return p, nil
}
