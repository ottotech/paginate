package paginate

import "fmt"

type parameters []parameter

func (params *parameters) getParameter(name string) (parameter, bool) {
	for _, p := range *params {
		if p.name == name {
			return p, true
		}
	}
	return parameter{}, false
}

// parameter holds information about a parameter passed in the url.Values from an http.Request.
type parameter struct {
	name  string
	sign  string
	value string
}

// paginationRequest holds information about the pagination operation.
// The helper function getRequestData() helps this package to get this
// information from the url.Values of an http.Request.
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

// whereClause holds information about an sql where clause.
type whereClause struct {
	clause string
	args   []interface{}
	exists bool
}

// mappers holds a collection of mapper objects.
type mappers []mapper

// mapper maps a database column name with a request parameter.
// We use this helper type because there might be cases where the user
// of the package wants to have a different name for the request parameter
// instead of using a real column name from the given table.
type mapper struct {
	col, param string
}

// isColumnMapped checks whether the given column name is mapped with a
// request parameter or not.
func (m mappers) isColumnMapped(columnName string) (columnIsMapped bool, customParameterName string) {
	for _, x := range m {
		if x.col == columnName {
			return true, x.param
		}
	}
	return false, ""
}

// Add adds a mapper in mappers. If an equal
// mapper with the same ``col`` and ``param`` values
// exists, Add will not add the given values to mapper.
func (m *mappers) Add(col, param string) {
	for _, x := range *m {
		if x.col == col && x.param == param {
			return
		}
	}
	*m = append(*m, mapper{
		col:   col,
		param: param,
	})
}

type __dialectPlaceholder map[string]string

func (d __dialectPlaceholder) GetPlaceHolder(dialect string) string {
	return d[dialect]
}

func (d __dialectPlaceholder) CheckIfDialectIsSupported(dialect string) error {
	for k, _ := range d {
		if k == dialect {
			return nil
		}
	}
	return fmt.Errorf("paginate: given dialect %q is not supported by this package", dialect)
}

var dialectPlaceholder = __dialectPlaceholder{
	"mysql":    "?",
	"postgres": "$%v", // This can become later in $1 see: Paginate() implementation for more.
}

type orderBy []string

func (o *orderBy) UniqueValues(skipId string) {
	unique := make([]string, 0)
	for _, s := range *o {
		if s == skipId {
			continue
		}
		in := isStringIn(s, unique)
		if !in {
			unique = append(unique, s)
		}
	}
	*o = nil
	for _, x := range unique {
		*o = append(*o, x)
	}
}
