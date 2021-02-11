/*
Package paginate provides a basic Paginator object to paginate records of a single database table.
Its primary job is to generate a raw sql command with the corresponding arguments that can be
executed against a sequel database with an sql driver of your preference.

Paginator also provides some utility functions like GetRowPtrArgs, NextData, and Scan, to make
easy to retrieve and read the paginated data.

Paginator also handles basic filtering of records with the parameters coming from a request url.

And Paginator also returns a PaginationResponse which contains useful information for clients to do proper
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

Similarly, when parameters with the notequal sign (<>) in the request url are repeated, Paginator will
interpret this as a NOT IN sql clause. So for example, given a database table ``Employees``
and a request url like:

	http://localhost/employees?name<>julia&name<>mark

Paginator will produce an sql query similar to this:

	SELECT id, name FROM employees WHERE name NOT IN($1,$2) ORDER BY id LIMIT 30 OFFSET 0

Example of the table struct field tags and their meanings
(use a ; to specify multiple tags at the same time):

	// Paginator will infer the database column name from the tag "col",
	// so in this case the column name would be "person_name".
	// If no "col" tag is given Paginator will infer the database
	// column name from the name of the struct field in snake lowercase
	// form. So given a struct field "MyName", Paginator will infer the
	// the database column name as "my_name".
	Name string `paginate:"col=person_name"`

	// By default all the columns of the given table cannot be filtered
	// with the parameters coming from a request url. Nevertheless, a user can
	// explicitly tell Paginator to filter a column by specifying the "filter"
	// tag in the table struct fields.
	LastName string `paginate:"filter"`

	// By default, once a database column has been defined as filterable
	// with the tag "filter", Paginator will map the request parameters with
	// the column names of the given table and it will filter the table
	// accordingly. However, in some cases you might want to have different
	// names for your request parameters in order to be more expressive.
	// In those cases use the tag "param" to define a custom request parameter
	// name that can be mapped to a column of the database table.
	// So, for example, in this case a request parameter "person_id" will be
	// used to filter the records of the given table based on the values of
	// the column "id".
	ID int `paginate:"col=id;param=person_id"`

	// The tag "id" is required. If it is not given, Paginator cannot be instantiated
	// and it will return an error. The tag "id" allows Paginator to keep the same order
	// between pages and results. In simple words, it will make the pagination deterministic.
	// Usually, the column used with this tag should be the primary key or unique identifier
	// of the given table.
	ID int `paginate:"id"`

NOTES:

Paginator does not take into consideration performance since it uses the OFFSET sql argument
which reads and counts all rows from the beginning until it reaches the requested page. For
not too big datasets Paginator will just work fine. If you care about performance because you
are dealing with heavy data you might want to write a custom solution for that.

USAGE EXAMPLES:

Check the examples folder in the repository of the package to learn more about how use paginate.
*/
package paginate
