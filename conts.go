package paginate

// Group of constants that represent some default values used by the paginate package.
const (
	defaultPageSize   = 30
	defaultPageNumber = 1
	tagsep            = ";"
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

// Constants that represent
const (
	_in    = "IN"
	_notin = "NOT IN"
)
