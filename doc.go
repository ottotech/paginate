/*
Package paginate provides a basic Paginator interface to do pagination of database records
for a single database table. Its primary job is to generate a raw sql command with the
corresponding arguments that can be executed with an sql database driver. Paginator also provides
some utility functions like GetRowPtrArgs, NextData, and Scan, to make easy to retrieve
and read the paginated data.

Paginator also handles basic filtering of records with the parameters coming from the
request url.

Paginator also returns a PaginationResponse, which contains useful information for clients to do proper
pagination.

For filtering database records the following operators are available.
Use these with the parameters in the request url:
	eq  = "="
	gt  = ">"
	lt  = "<"
	gte = ">="
	lte = "<="
	ne  = "<>"


For ordering records based on column names use the following syntax in the url with the ``sort``
parameter. For sorting in ascending order use the plus (+) sign, and for sorting in descending
order use the minus (-) sign:

	http://localhost/employees?name=rob&sort=+name,-age


When parameters with the equal sign (=) in the request url are repeated, Paginator will
interpret this as an IN sql clause. So for example given a database table ``Employees``
and a request url like:

http://localhost/employees?name=julia&name=mark

Paginator will produce an sql query similar to this:

	SELECT id, name FROM employees WHERE name IN($1,$2) ORDER BY id LIMIT 30 OFFSET 0

When parameters with the notequal sign (<>) in the request url are repeated, Paginator will
interpret this as a NOT IN sql clause. So for example given a database table ``Employees``
and a request url like:

	http://localhost/employees?name<>julia&name<>mark

Paginator will produce an sql query similar to this:

	SELECT id, name FROM employees WHERE name NOT IN($1,$2) ORDER BY id LIMIT 30 OFFSET 0

Example of the table struct field tags and their meanings:

	// Paginator will take the database column name from the tag,
	// so in this case the column name would be "person_name".
	Name string `paginate:"col=person_name"`
*/
package paginate
