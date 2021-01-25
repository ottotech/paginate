package paginate

// Constant that specifies the page size of the pagination results.
// This value will be used in cases where the page size cannot be determined.
// For example, when using NewPaginator, if NewPaginator fails to retrieve
// the page size from the request url it will fallback to this value.
const (
	PageSize = 30
	tagsep   = ";"
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
