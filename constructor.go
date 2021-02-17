package paginate

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

type Option func(p *paginator) error

// TableName is an option for NewPaginator which indicates
// the name of the table that paginator will paginate in the
// database. Use this option if the name of your database
// table cannot be inferred by the given struct table name.
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
// the size of the record set that we want our paginator
// object to produce per page. ``size`` should be an uint value
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

func OrderByAsc(column string) Option {
	return func(p *paginator) error {
		p.orderBy = append(p.orderBy, fmt.Sprintf("%s %s", column, "ASC"))
		return nil
	}
}

func OrderByDesc(column string) Option {
	return func(p *paginator) error {
		p.orderBy = append(p.orderBy, fmt.Sprintf("%s %s", column, "DESC"))
		return nil
	}
}

// NewPaginator creates a Paginator object ready to paginate data from a database table.
//
// The table parameter should be a struct object with fields representing the target
// database table you want to paginate. The dialect parameter should be a string
// representing the sql dialect you are using "postgres" or "mysql", for example.
// For available options you can pass to Paginator check: TableName and PageSize.
//
// When the PageSize option is not given paginator will try to get the page size from the
// request parameter ``page_size``. If there is no ``page_size`` parameter NewPaginator
// will set the Paginator with the default page size which is 30. When the TableName option
// is not given, NewPaginator will infer the database table name from the table argument
// given, so it will extract the name from the struct variable.
func NewPaginator(table interface{}, dialect string, u url.URL, opts ...Option) (Paginator, error) {
	err := dialectPlaceholder.CheckIfDialectIsSupported(dialect)
	if err != nil {
		return nil, err
	}

	p := &paginator{table: table, rv: reflect.ValueOf(table), dialect: dialect}

	// Let's try to set the options if any.
	for _, opt := range opts {
		err := opt(p)
		if err != nil {
			return nil, err
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

	// Now let's get all the data that our paginator requires.
	// Order matters: getCols func should be called before
	// getParameters func.
	p.getID()
	p.getColsAndMapParameters()
	p.getFieldNames()
	p.getFilters()
	p.parameters = getParameters(p.cols, p.filters, p.mappers, u)

	// Let's clean our orderBy slice and get unique values.
	p.orderBy.UniqueValues(p.id)

	return p, nil
}
