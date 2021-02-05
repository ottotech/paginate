package paginate

type parameters []parameter

// getParameter will try to get a parameter with the given name from
// parameters. If there is no parameter with that name getParameter
// will return false, otherwise true.
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

// paginationRequest holds information about the paginator that the client
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
