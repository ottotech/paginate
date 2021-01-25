/*
Package paginate provides a basic Paginator interface to do pagination of database records
of a single database table. Its primary job is to generate a raw sql command with the corresponding
arguments that can be executed in a database.

This package will also handle basic filtering of records with the parameters coming from a request.
And –if used correctly- Paginator can also return a PaginationResponse which contains useful information
for clients to do proper pagination.

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
