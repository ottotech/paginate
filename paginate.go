package paginate

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"
)

// ErrPaginatorIsClosed is an error returned by Scan when trying to Scan
// on a Paginator that is closed. That is, a paginator whose values were
// already scanned.
var ErrPaginatorIsClosed = errors.New("paginate: Paginator is closed")

// Paginator wraps pagination behaviors.
//
// Paginator should be used following the next steps in the same order:
//
// 		1. Initialize a Paginator instance with NewPaginator.
// 		2. Call Paginate to create the query and arguments to be executed with an sql driver.
// 		3. Call GetRowPtrArgs when scanning the rows inside the sql.Rows.Next loop.
// 		4. Call NextData to loop over the paginated data and Scan the data afterwards.
// 		5. Call Scan inside the NextData loop to copy the paginated data to the given destination.
// 		6. Call Response to get useful information about the pagination operation.
//
// For more information, see the examples folder to check how to use Paginator.
type Paginator interface {
	// Paginate will return an sql command with the corresponding arguments,
	// so it can be run against any sql driver.
	Paginate() (sql string, args []interface{}, err error)

	// GetRowPtrArgs will prepare the next pointer arguments that can be scanned
	// by sql.Rows.Scan.
	//
	// Always run GetRowPtrArgs when scanning the queried rows with the sql package,
	// for example:
	//
	//   rows, _ := db.Query(myQuery)
	//	 for rows.Next() {
	//		err = rows.Scan(paginator.GetRowPtrArgs()...)
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//	 }
	//
	// Every time GetRowPtrArgs gets called it will save the previous scanned values
	// internally in the Paginator object so you can read/scan them later.
	GetRowPtrArgs() []interface{}

	// NextData will loop over the saved values created by GetRowPtrArgs until
	// all the paginated data has been scanned by Scan. Always use NextData
	// followed by a call to Scan.
	NextData() bool

	// Scan will copy the next paginated data in the given destination. The given destination
	// should be a pointer instance of the same type of the given ``table``.
	//
	// Scan converts columns read from the database into the following
	// common Go types:
	//
	//    *string
	//	  *[]byte
	//    *int, *int8, *int16, *int32, *int64
	//    *uint, *uint8, *uint16, *uint32, *uint64
	//    *bool
	//    *float32, *float64
	//
	// Scan will also convert nullable fields of type string, int32, int64, float64,
	// bool, time.Time to their default zero values with the following helpers provided
	// by the sql package (As the user of this package you do not have to care about this,
	// paginate will automatically handle this for you):
	//
	// 	  - sql.NullString
	// 	  - sql.NullInt32
	// 	  - sql.NullInt64
	// 	  - sql.NullFloat64
	// 	  - sql.NullBool
	// 	  - sql.NullTime
	//
	// For other nullable fields that you might want Scan to handle, use
	// the nullable types provided by this package.
	Scan(dest interface{}) error

	// Response returns a PaginationResponse containing useful information about
	// the pagination, so that clients can do proper and subsequent pagination
	// operations.
	Response() PaginationResponse

	// AddWhereClause adds a custom raw where clause that paginator can use to
	// filter out the rows of the target table in the database. Usually, you will
	// use this when the backend needs to filter the records based on internal logic,
	// like, for example: permissions, ownership, etc.
	AddWhereClause(clause RawWhereClause) error
}

// paginator is the concrete type that implements the Paginator interface.
type paginator struct {
	// dialect represents the sql dialect that paginator will use to build
	// the sql command.
	dialect string

	// table is a representation of a database table and its columns. The
	// given table should be of type struct and its fields represent the
	// columns of the table in the database.
	table interface{}

	// predicates holds custom where clauses created by a user of this
	// package which will be executed in the helper createWhereClause.
	predicates []RawWhereClause

	// stop is used by NextData and Scan. Scan will set the value of stop
	// to true whenever Scan returns an error. This will allow NextData to
	// stop looping over p.rows.
	stop bool

	// name is the name of the table in the database. This package will
	// infer the name of the table from name of the given table struct
	// if the name is not provided.
	name string

	// id represents the pk or unique identifier of the table in the database.
	// This value should be defined in the given table through the tag "id"
	// (e.g. `paginator:"id"`) in one of the `table` struct `fields`.
	// This value is very important since it will make the sort order
	// deterministic when paginating the data.
	id string

	// cols holds the names of the columns of the database table.
	// This package will infer the column names from the struct `fields`
	// of the given table and it will convert any camel case field name
	// into sneak case lowercase. So "MyAwesomeField" will be "my_awesome_field".
	// If the tag "col" is given in one of the fields of the given struct
	// table, the column `name` will be taken from there.
	cols []string

	// fields holds the raw names of the struct "fields" of the given table.
	fields []string

	// filters holds the names of the columns of the table that the user
	// wants to filter. By default all the fields of the table struct cannot
	// be filtered. A user can explicitly tell paginator to filter a
	// column by specifying the "filter" tag in the table struct fields.
	filters []string

	// rows holds the paginated values scanned by the Scan method from the
	// go sql package from the stdlib. See also addRow and GetRowPtrArgs for a
	// better understanding.
	rows []interface{}

	// tmp holds the values for each row in rows. tmp will hold the values
	// temporarily everytime we run GetRowPtrArgs. We use these values to
	// scan with the Scan method from the go sql package the paginated data.
	tmp []interface{}

	// rv holds the reflection value of the given table.
	rv reflect.Value

	// pageSize represents the size of the records, in number of rows, that we want to
	// show per page. We will get this value from the request url values if given.
	// We can also get this value from the PageSize option. Otherwise, we will fallback
	// to the value defined by defaultPageSize.
	pageSize int

	// pageNumber represents the "number" of the page that the end user wants
	// to see for the paginated data. We will get this value from the request url values
	// if given. Otherwise, we will fallback to value defined by defaultPageNumber.
	pageNumber int

	// totalSize represents the total number of records in the given table in
	// the database.
	totalSize int

	// pageCount represents the total number of records retrieved by paginator
	// from the database.
	pageCount int

	// closed is used by Scan. When closed == true it means that all rows
	// have been scanned to the given destinations and paginator does not have
	// the scanned rows anymore, so further call to Scan will not work at
	// this point.
	closed bool

	// started is used by GetRowPtrArgs and Scan. When started == true it means
	// that the rows consumption has started, so no further calls to GetRowPtrArgs
	// can be done. A call to GetRowPtrArgs at this point will return nil.
	started bool

	// once is used by Scan. It's purpose is to set only ``once``: started, pageSize,
	// and run addRow the first time Scan is used.
	once sync.Once

	// parameters holds the parameters required to create the sql where clause
	// that we will use to paginate and filter the database table.
	parameters parameters

	// mappers holds a collection of mapper objects. See mappers documentation for more.
	mappers mappers

	// response holds useful information about the pagination operation. Clients can use
	// this information to do subsequent pagination calls.
	response PaginationResponse

	// orderBy holds custom "ORDER BY" clauses that will be added to the generated
	// sql command. See createOrderByClause.
	orderBy orderBy
}

func (p *paginator) Paginate() (sql string, values []interface{}, err error) {
	var s string
	c1 := make(chan whereClause)
	c2 := make(chan string)
	c3 := make(chan string)
	go createWhereClause(p.dialect, p.cols, p.parameters, p.predicates, c1)
	go createPaginationClause(p.pageNumber, p.pageSize, c2)
	go createOrderByClause(p.parameters, p.cols, p.orderBy, p.id, c3)
	where := <-c1
	pagination := <-c2
	order := <-c3

	if where.exists {
		s = "SELECT " + strings.Join(p.cols, ", ") + ", count(*) over() FROM " + p.name + where.clause + order + pagination
		// As an special case we need to enumerate the placeholders if users are using
		// postgres. See, for example, the documentation of this postgres driver library:
		// https://pkg.go.dev/github.com/lib/pq#section-documentation
		if p.dialect == "postgres" {
			numArgs := len(where.args)
			placeholders := make([]interface{}, 0)
			for i := 1; i < numArgs+1; i++ {
				placeholders = append(placeholders, i)
			}
			s = fmt.Sprintf(s, placeholders...)
		}
	} else {
		s = "SELECT " + strings.Join(p.cols, ", ") + ", count(*) over() FROM " + p.name + order + pagination
	}
	return s, where.args, nil
}

func (p *paginator) Response() PaginationResponse {
	p.response.PageNumber = p.pageNumber
	p.response.PageCount = p.pageCount
	p.response.TotalSize = p.totalSize

	if (p.pageNumber * p.pageSize) < p.totalSize {
		p.response.NextPageNumber = p.pageNumber + 1
		p.response.HasNextPage = true
	}
	if (p.pageNumber * p.pageSize) == p.totalSize {
		p.response.NextPageNumber = 0
		p.response.HasNextPage = false
	}
	if p.totalSize == 0 {
		p.response.NextPageNumber = 0
		p.response.HasNextPage = false
	}
	if p.response.PageNumber > 1 {
		p.response.HasPreviousPage = true
	}

	return p.response
}

// validateTable validates if the given table struct is valid.
func (p *paginator) validateTable() error {
	if p.rv.Type().Kind() != reflect.Struct {
		return fmt.Errorf("paginate: table should be of struct type")
	}

	if !p.rv.IsZero() {
		return fmt.Errorf("paginate: table struct should be empty with only the default zero values")
	}

	numOfIDs := 0

	// See usage below.
	countIDs := func(tags []string) int {
		s := "id"
		c := 0
		for _, tag := range tags {
			if tag == s {
				c++
			}
		}
		return c
	}

	for i := 0; i < p.rv.NumField(); i++ {
		field := p.rv.Type().Field(i)
		tags := strings.Split(field.Tag.Get("paginate"), ";")
		numOfIDs += countIDs(tags)
		fieldName := field.Name
		T := reflect.Indirect(p.rv).FieldByName(fieldName).Interface()
		switch T.(type) {
		case string:
			continue
		case int, int8, int16, int32, int64:
			continue
		case bool:
			continue
		case float32, float64:
			continue
		case time.Time:
			continue
		case NullInt, NullBool, NullString, NullTime, NullFloat64:
			continue
		default:
			return fmt.Errorf("paginate: invalid type for field %q", fieldName)
		}
	}

	if numOfIDs == 0 {
		return fmt.Errorf("paginate: id has not been defined in any " +
			"of the field of the given struct")
	} else if numOfIDs > 1 {
		return fmt.Errorf("paginate: more than one id has beend defined " +
			"in the fields of the given struct")
	}

	return nil
}

// getColsAndMapParameters does two things:
//
// (1) It infers the column names of the database table from the given ``table``
//     struct fields. If the fields have the tag "col", the column name will be taken
//     from there.
// (2) It will map the column names with request ``parameter`` names if the struct
//     fields have the tag "param" on it.
//
// Malformed "col" and "param" tags will be ignored silently.
func (p *paginator) getColsAndMapParameters() {

	getColNameFromTags := func(tags []string) (hasColTag bool, colName string) {
		for _, tag := range tags {
			kv := strings.Split(tag, "=")
			if len(kv) != 2 {
				continue
			}
			k, v := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])
			if k != col {
				continue
			}
			return true, v
		}
		return hasColTag, colName
	}

	getParamFromTags := func(tags []string) (hasParamTag bool, paramName string) {
		for _, tag := range tags {
			kv := strings.Split(tag, "=")
			if len(kv) != 2 {
				continue
			}
			k, v := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])
			if k != param {
				continue
			}
			return true, v
		}
		return hasParamTag, paramName
	}

	for i := 0; i < p.rv.NumField(); i++ {
		field := p.rv.Type().Field(i)
		tags := strings.Split(field.Tag.Get("paginate"), tagsep)
		if hasColTag, name := getColNameFromTags(tags); hasColTag {
			p.cols = append(p.cols, name)
			if hasParamTag, paramName := getParamFromTags(tags); hasParamTag {
				p.mappers.Add(name, paramName)
			}
			continue
		}
		fieldName := field.Name
		sneakName := parseCamelCaseToSnakeLowerCase(fieldName)
		if hasParamTag, paramName := getParamFromTags(tags); hasParamTag {
			p.mappers.Add(sneakName, paramName)
		}
		p.cols = append(p.cols, sneakName)
	}
}

func (p *paginator) getFieldNames() {
	for i := 0; i < p.rv.NumField(); i++ {
		field := p.rv.Type().Field(i)
		fieldName := field.Name
		p.fields = append(p.fields, fieldName)
	}
}

func (p *paginator) getFilters() {

	hasfilter := func(tags []string) bool {
		for _, tag := range tags {
			if tag == filter {
				return true
			}
		}
		return false
	}

	for i := 0; i < p.rv.NumField(); i++ {
		field := p.rv.Type().Field(i)
		tags := strings.Split(field.Tag.Get("paginate"), tagsep)
		if !hasfilter(tags) {
			continue
		}
		sneakName := parseCamelCaseToSnakeLowerCase(field.Name)
		p.filters = append(p.filters, sneakName)
	}
}

func (p *paginator) getID() {
	id := ""

	hasID := func(tags []string) bool {
		for _, tag := range tags {
			if tag == "id" {
				return true
			}
		}
		return false
	}

	for i := 0; i < p.rv.NumField(); i++ {
		field := p.rv.Type().Field(i)
		tags := strings.Split(field.Tag.Get("paginate"), tagsep)
		if hasID(tags) {
			id = parseCamelCaseToSnakeLowerCase(field.Name)
			break
		}
	}

	p.id = id
}

func (p *paginator) getTableName() {
	name := parseCamelCaseToSnakeLowerCase(p.rv.Type().Name())
	p.name = name
}

func (p *paginator) GetRowPtrArgs() []interface{} {
	if p.started {
		return nil
	}
	if len(p.tmp) > 0 {
		p.addRow()
	}
	for _, fieldName := range p.fields {
		I := reflect.Indirect(p.rv).FieldByName(fieldName).Interface()
		switch I.(type) {
		case NullInt:
			var ni NullInt
			p.tmp = append(p.tmp, &ni)
		case NullBool:
			var nb NullBool
			p.tmp = append(p.tmp, &nb)
		case NullString:
			var ni NullString
			p.tmp = append(p.tmp, &ni)
		case NullTime:
			var nt NullTime
			p.tmp = append(p.tmp, &nt)
		case NullFloat64:
			var n NullFloat64
			p.tmp = append(p.tmp, &n)
		case string:
			var s sql.NullString
			p.tmp = append(p.tmp, &s)
		case int:
			var i int
			p.tmp = append(p.tmp, &i)
		case int8:
			var i8 int8
			p.tmp = append(p.tmp, &i8)
		case int16:
			var i16 int16
			i16 = 0
			p.tmp = append(p.tmp, &i16)
		case int32:
			var i32 sql.NullInt32
			p.tmp = append(p.tmp, &i32)
		case int64:
			var i64 sql.NullInt64
			p.tmp = append(p.tmp, &i64)
		case bool:
			var b sql.NullBool
			p.tmp = append(p.tmp, &b)
		case float32:
			var f32 float32
			p.tmp = append(p.tmp, &f32)
		case float64:
			var f64 sql.NullFloat64
			p.tmp = append(p.tmp, &f64)
		case time.Time:
			var t sql.NullTime
			p.tmp = append(p.tmp, &t)
		}
	}

	// As an special case in tmp we will always
	// append at the end p.totalSize whose value
	// is going to be set when the query gets executed.
	p.tmp = append(p.tmp, &p.totalSize)

	return p.tmp
}

// addRow adds a new row in p.rows.
//
// The elements of p.rows will be instances of the given table struct.
//
// addRow will use the p.tmp temporarily values which are pointers
// scanned by the sql driver to fill in the table struct fields (a "row").
// Once the table struct fields are set addRow will add the table struct
// to p.rows. As an special case addRow will handle nullable fields with
// the following nullable types from the sql package:
//
// 		- sql.NullString
// 		- sql.NullInt32
// 		- sql.NullInt64
// 		- sql.NullFloat64
// 		- sql.NullBool
// 		- sql.NullTime
//
// Custom nullable values will not be handled, for example, for
// values like uint8, uint16, etc.
//
// It is up to GetRowPtrArgs to call addRow each time a new row is
// read by sql.Rows.Scan. NextData is also responsible to call addRow
// just in case there are values left in p.tmp. It might be possible
// that there are values left in p.tmp because the last call to GetRowPtrArgs,
// for example, will not call addRow to add the p.tmp values into p.rows.
// Finally, Scan will also call addRow only once in case there values left in
// p.tmp.
func (p *paginator) addRow() {
	row := p.table
	rowrv := reflect.ValueOf(&row).Elem()
	tmpRow := reflect.New(rowrv.Elem().Type()).Elem()
	tmpRow.Set(rowrv.Elem())

	// The below loop condition expression is len(tmp)-1
	// because of the extra field we are adding in p.tmp: totalSize.
	// len(tmp)-1  will give us exactly the field elements we want
	// from the given table.
	for i := 0; i < len(p.tmp)-1; i++ {
		I := reflect.Indirect(reflect.ValueOf(p.tmp[i])).Interface()
		tmpRowField := tmpRow.FieldByName(p.fields[i])

		switch I.(type) {
		case sql.NullString:
			ns := sql.NullString{}
			nsrv := reflect.ValueOf(&ns).Elem()
			nsrv.Set(reflect.ValueOf(p.tmp[i]).Elem())
			tmpRowField.SetString(ns.String)
		case sql.NullInt32:
			ni32 := sql.NullInt32{}
			ni32rv := reflect.ValueOf(&ni32).Elem()
			ni32rv.Set(reflect.ValueOf(p.tmp[i]).Elem())
			tmpRowField.Set(reflect.ValueOf(ni32.Int32))
		case sql.NullInt64:
			ni64 := sql.NullInt64{}
			ni64rv := reflect.ValueOf(&ni64).Elem()
			ni64rv.Set(reflect.ValueOf(p.tmp[i]).Elem())
			tmpRowField.Set(reflect.ValueOf(ni64.Int64))
		case sql.NullFloat64:
			nf64 := sql.NullFloat64{}
			nf64rv := reflect.ValueOf(&nf64).Elem()
			nf64rv.Set(reflect.ValueOf(p.tmp[i]).Elem())
			tmpRowField.Set(reflect.ValueOf(nf64.Float64))
		case sql.NullBool:
			nb := sql.NullBool{}
			nbrv := reflect.ValueOf(&nb).Elem()
			nbrv.Set(reflect.ValueOf(p.tmp[i]).Elem())
			tmpRowField.Set(reflect.ValueOf(nb.Bool))
		case sql.NullTime:
			nt := sql.NullTime{}
			ntrv := reflect.ValueOf(&nt).Elem()
			ntrv.Set(reflect.ValueOf(p.tmp[i]).Elem())
			tmpRowField.Set(reflect.ValueOf(nt.Time))
		default:
			val := reflect.ValueOf(p.tmp[i]).Elem()
			tmpRow.FieldByName(p.fields[i]).Set(val)
		}

		rowrv.Set(tmpRow)
	}

	// We need to clear p.tmp so we can reuse it later for another call
	// to addRow.
	p.tmp = make([]interface{}, 0)

	p.rows = append(p.rows, row)
}

func (p *paginator) NextData() bool {
	if len(p.tmp) > 0 {
		p.addRow()
	}
	if p.stop || p.closed {
		return false
	}
	return len(p.rows) > 0
}

func (p *paginator) Scan(dest interface{}) (err error) {
	defer func() {
		if err != nil {
			p.stop = true
		}
	}()

	if dest == nil {
		return fmt.Errorf("paginate: cannot pass nil as dest")
	}

	if p.closed {
		return ErrPaginatorIsClosed
	}

	if err = p.validateDest(dest); err != nil {
		return err
	}

	if len(p.rows) == 0 {
		return errors.New("paginate: Scan called without calling NextData")
	}

	p.once.Do(func() {
		// Order matters. If there is some data left in p.tmp,
		// p.addRow will add a new row with the p.tmp data affecting,
		// therefore, the value of p.pageCount.
		if len(p.tmp) > 0 {
			p.addRow()
		}
		p.started = true
		p.pageCount = len(p.rows)
	})

	destrv := reflect.ValueOf(dest)

	row := p.rows[0]
	for _, field := range p.fields {
		val := reflect.ValueOf(row).FieldByName(field)
		destrv.Elem().FieldByName(field).Set(val)
	}

	// Let's remove the row from p.rows.
	p.rows = p.rows[1:]

	// When all rows are consumed, we "close" the Paginator Scanner.
	if len(p.rows) == 0 {
		p.closed = true
	}

	return nil
}

func (p *paginator) validateDest(dest interface{}) error {
	destrv := reflect.ValueOf(dest)

	if destrv.Kind() != reflect.Ptr {
		return fmt.Errorf("paginate: the given "+
			"destination should be a pointer of type *%s; got %s",
			p.rv.Type().String(), destrv.Type().String())
	}

	destri := reflect.Indirect(destrv)

	if p.rv.Type() != destri.Type() {
		return fmt.Errorf("paginate: the given "+
			"destination should be a pointer of type *%s; got %s",
			p.rv.Type().String(), destri.Type().String())
	}

	if !destri.IsZero() {
		return fmt.Errorf("paginate: the given "+
			"destination should be the zero value of *%s",
			p.rv.Type().String())
	}

	return nil
}

func (p *paginator) AddWhereClause(clause RawWhereClause) error {
	re := regexp.MustCompile("[?]")
	occurrences := re.FindAll([]byte(clause.predicate), -1)
	if len(occurrences) == 0 && len(clause.args) > 0 {
		return fmt.Errorf("paginate: cannot receive arguments when placeholders are not defined")
	}
	if len(occurrences) > 0 && len(occurrences) != len(clause.args) {
		return fmt.Errorf("paginate: the number of placeholders and arguments in the where clause should be the same")
	}

	if err := dialectPlaceholder.CheckIfDialectIsSupported(clause.dialect); err != nil {
		return fmt.Errorf("the dialect specified in the RawWhereClause is not supported")
	}

	p.predicates = append(p.predicates, clause)
	return nil
}
