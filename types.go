package paginate

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

// whereClause holds information about a where clauses.
type whereClause struct {
	clause string
	args   []interface{}
	exists bool
}
