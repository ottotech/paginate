package paginate

import (
	"fmt"
	"log"
	"net/url"
	//"reflect"
	"testing"
)

func ExampleNewPaginator_1() {
	tableName := "test"
	colNames := []string{"system"}
	u, err := url.Parse("http://ottotech.com?system=hipaca")
	if err != nil {
		log.Fatal(err)
	}
	paginator := NewPaginator(tableName, colNames, *u)
	sql, args, _ := paginator.Paginate()
	fmt.Println(sql)
	fmt.Printf("arg $1: %v\n", args[0])
	// Output:
	// SELECT system, count(*) over() FROM test WHERE system = $1 ORDER BY id LIMIT 30 OFFSET 0
	// arg $1: hipaca
}

func ExampleNewPaginator_2() {
	tableName := "test"
	colNames := []string{"system", "name"}
	u, err := url.Parse("http://ottotech.com?system=hipaca&name=martha")
	if err != nil {
		log.Fatal(err)
	}
	paginator := NewPaginator(tableName, colNames, *u)
	sql, args, _ := paginator.Paginate()
	fmt.Printf("arg $1: %v\n", args[0])
	fmt.Printf("arg $2: %v\n", args[1])
	fmt.Println(sql)
	// Unordered output:
	// SELECT system, name, count(*) over() FROM test WHERE system = $1 AND name = $2 ORDER BY id LIMIT 30 OFFSET 0
	// arg $1: hipaca
	// arg $2: martha
}

func ExampleNewPaginator_3() {
	tableName := "test"
	colNames := []string{"system", "name", "lastname"}
	u, err := url.Parse("http://ottotech.com?system=hipaca&page_size=8&lastname<>schuldt")
	if err != nil {
		log.Fatal(err)
	}
	paginator := NewPaginator(tableName, colNames, *u)
	sql, args, _ := paginator.Paginate()
	fmt.Println(sql)
	fmt.Printf("arg $1: %v\n", args[0])
	fmt.Printf("arg $2: %v\n", args[1])
	// Output:
	// SELECT system, name, lastname, count(*) over() FROM test WHERE system = $1 AND lastname <> $2 ORDER BY id LIMIT 8 OFFSET 0
	// arg $1: hipaca
	// arg $2: schuldt
}

func ExampleNewPaginator_4() {
	tableName := "test"
	colNames := []string{"system", "name", "lastname"}
	u, err := url.Parse("http://ottotech.com?system=hipaca&sort=+name,-lastname&page_size=7")
	if err != nil {
		log.Fatal(err)
	}
	paginator := NewPaginator(tableName, colNames, *u)
	sql, args, _ := paginator.Paginate()
	fmt.Println(sql)
	fmt.Printf("arg $1: %v\n", args[0])
	// Output:
	// SELECT system, name, lastname, count(*) over() FROM test WHERE system = $1 ORDER BY name ASC,lastname DESC,id LIMIT 7 OFFSET 0
	// arg $1: hipaca
}

func ExampleNewPaginator_5() {
	tableName := "test"
	colNames := []string{"column1", "column2", "column3", "column4"}
	u, err := url.Parse("http://ottotech.com?column1>2&column2<4&column3>=40&column4<=7")
	if err != nil {
		log.Fatal(err)
	}
	paginator := NewPaginator(tableName, colNames, *u)
	sql, args, _ := paginator.Paginate()
	fmt.Println(sql)
	fmt.Printf("arg $1: %v\n", args[0])
	fmt.Printf("arg $2: %v\n", args[1])
	fmt.Printf("arg $3: %v\n", args[2])
	fmt.Printf("arg $4: %v\n", args[3])
	// Output:
	// SELECT column1, column2, column3, column4, count(*) over() FROM test WHERE column1 > $1 AND column2 < $2 AND column3 >= $3 AND column4 <= $4 ORDER BY id LIMIT 30 OFFSET 0
	// arg $1: 2
	// arg $2: 4
	// arg $3: 40
	// arg $4: 7
}

func ExampleNewPaginatorWithLimit_1() {
	tableName := "test"
	colNames := []string{"system"}
	u, err := url.Parse("http://ottotech.com?system=hipaca")
	if err != nil {
		log.Fatal(err)
	}
	paginator := NewPaginatorWithLimit(10, tableName, colNames, *u)
	sql, args, _ := paginator.Paginate()
	fmt.Println(sql)
	fmt.Printf("arg $1: %v\n", args[0])
	// Output:
	// SELECT system, count(*) over() FROM test WHERE system = $1 ORDER BY id LIMIT 10 OFFSET 0
	// arg $1: hipaca
}

func ExampleNewPaginatorWithLimit_2() {
	tableName := "test"
	colNames := []string{"system", "name"}
	u, err := url.Parse("http://ottotech.com?system=hipaca&name=martha")
	if err != nil {
		log.Fatal(err)
	}
	paginator := NewPaginatorWithLimit(50, tableName, colNames, *u)
	sql, args, _ := paginator.Paginate()
	fmt.Printf("arg $1: %v\n", args[0])
	fmt.Printf("arg $2: %v\n", args[1])
	fmt.Println(sql)
	// Unordered output:
	// SELECT system, name, count(*) over() FROM test WHERE system = $1 AND name = $2 ORDER BY id LIMIT 30 OFFSET 0
	// arg $1: hipaca
	// arg $2: martha
}

func ExampleNewPaginatorWithLimit_3() {
	tableName := "test"
	colNames := []string{"system", "name", "lastname"}
	u, err := url.Parse("http://ottotech.com?system=hipaca&page_size=8&lastname<>schuldt")
	if err != nil {
		log.Fatal(err)
	}
	paginator := NewPaginatorWithLimit(8, tableName, colNames, *u)
	sql, args, _ := paginator.Paginate()
	fmt.Println(sql)
	fmt.Printf("arg $1: %v\n", args[0])
	fmt.Printf("arg $2: %v\n", args[1])
	// Output:
	// SELECT system, name, lastname, count(*) over() FROM test WHERE system = $1 AND lastname <> $2 ORDER BY id LIMIT 8 OFFSET 0
	// arg $1: hipaca
	// arg $2: schuldt
}

func ExampleNewPaginatorWithLimit_4() {
	tableName := "test"
	colNames := []string{"system", "name", "lastname"}
	u, err := url.Parse("http://ottotech.com?system=hipaca&sort=+name,-lastname&page_size=7")
	if err != nil {
		log.Fatal(err)
	}
	paginator := NewPaginatorWithLimit(7, tableName, colNames, *u)
	sql, args, _ := paginator.Paginate()
	fmt.Println(sql)
	fmt.Printf("arg $1: %v\n", args[0])
	// Output:
	// SELECT system, name, lastname, count(*) over() FROM test WHERE system = $1 ORDER BY name ASC,lastname DESC,id LIMIT 7 OFFSET 0
	// arg $1: hipaca
}

func ExampleNewPaginatorWithLimit_5() {
	tableName := "test"
	colNames := []string{"column1", "column2", "column3", "column4"}
	u, err := url.Parse("http://ottotech.com?column1>2&column2<4&column3>=40&column4<=7")
	if err != nil {
		log.Fatal(err)
	}
	paginator := NewPaginatorWithLimit(30, tableName, colNames, *u)
	sql, args, _ := paginator.Paginate()
	fmt.Println(sql)
	fmt.Printf("arg $1: %v\n", args[0])
	fmt.Printf("arg $2: %v\n", args[1])
	fmt.Printf("arg $3: %v\n", args[2])
	fmt.Printf("arg $4: %v\n", args[3])
	// Output:
	// SELECT column1, column2, column3, column4, count(*) over() FROM test WHERE column1 > $1 AND column2 < $2 AND column3 >= $3 AND column4 <= $4 ORDER BY id LIMIT 30 OFFSET 0
	// arg $1: 2
	// arg $2: 4
	// arg $3: 40
	// arg $4: 7
}

func TestNewPaginator(t *testing.T) {
	tableName := "test"
	colNames := []string{"column1", "column2", "column3", "column4", "column5"}
	u, err := url.Parse("http://ottotech.com?column1>2&column2<4&column3>=40&column4<=7&column5=otto&sort=+column1,-column2&page_size=15")
	if err != nil {
		log.Fatal(err)
	}
	expectedSQL := "SELECT column1, column2, column3, column4, column5, count(*) over() FROM test WHERE column1 > $1 AND column2 < $2 AND column3 >= $3 AND column4 <= $4 AND column5 = $5 ORDER BY column1 ASC,column2 DESC,id LIMIT 15 OFFSET 0"
	paginator := NewPaginator(tableName, colNames, *u)
	paginator.SetPageCount(15)
	paginator.SetTotalResult(30)
	sql, args, err := paginator.Paginate()
	if sql != expectedSQL {
		t.Errorf("expected sql:\n %v; \ngot %v\n", expectedSQL, sql)
	}
	rightARGS := []interface{}{"2", "4", "40", "7", "otto"}
	for i := 0; i < len(args); i++ {
		if args[i] != rightARGS[i] {
			t.Errorf("arg $%v should be %v; got %v", i, rightARGS[i], args[i])
		}
	}
	res := paginator.Response()
	if res.PageNumber != 1 {
		t.Errorf("page number in pagination response should be 1; got %v", res.PageNumber)
	}
	if res.HasNextPage != true {
		t.Errorf("has_next_page in pagination response should be true; got %v", res.HasNextPage)
	}
	if res.HasPreviousPage != false {
		t.Errorf("has_previous_page in pagination response should be false; got %v", res.HasPreviousPage)
	}
	if res.NextPageNumber != 2 {
		t.Errorf("next_page_number in pagination response should be 2; got %v", res.NextPageNumber)
	}
	if res.TotalSize != 30 {
		t.Errorf("total_size in pagination response should be 30; got %v", res.TotalSize)
	}
	if res.PageCount != 15 {
		t.Errorf("page_count in pagination response should be 15; got %v", res.PageCount)
	}
}

func TestNewPaginatorWithLimit(t *testing.T) {
	tableName := "test"
	colNames := []string{"column1", "column2", "column3", "column4", "column5"}
	u, err := url.Parse("http://ottotech.com?column1>2&column2<4&column3>=40&column4<=7&column5=otto&sort=+column1,-column2&page_size=15&page=2")
	if err != nil {
		log.Fatal(err)
	}
	expectedSQL := "SELECT column1, column2, column3, column4, column5, count(*) over() FROM test WHERE column1 > $1 AND column2 < $2 AND column3 >= $3 AND column4 <= $4 AND column5 = $5 ORDER BY column1 ASC,column2 DESC,id LIMIT 10 OFFSET 10"
	paginator := NewPaginatorWithLimit(10, tableName, colNames, *u)
	paginator.SetPageCount(10)
	paginator.SetTotalResult(30)
	sql, args, err := paginator.Paginate()
	if sql != expectedSQL {
		t.Errorf("expected sql:\n %v; \ngot %v\n", expectedSQL, sql)
	}
	rightARGS := []interface{}{"2", "4", "40", "7", "otto"}
	for i := 0; i < len(args); i++ {
		if args[i] != rightARGS[i] {
			t.Errorf("arg $%v should be %v; got %v", i, rightARGS[i], args[i])
		}
	}
	res := paginator.Response()
	if res.PageNumber != 2 {
		t.Errorf("page number in pagination response should be 2; got %v", res.PageNumber)
	}
	if res.HasNextPage != true {
		t.Errorf("has_next_page in pagination response should be true; got %v", res.HasNextPage)
	}
	if res.HasPreviousPage != true {
		t.Errorf("has_previous_page in pagination response should be true; got %v", res.HasPreviousPage)
	}
	if res.NextPageNumber != 3 {
		t.Errorf("next_page_number in pagination response should be 2; got %v", res.NextPageNumber)
	}
	if res.TotalSize != 30 {
		t.Errorf("total_size in pagination response should be 30; got %v", res.TotalSize)
	}
	if res.PageCount != 10 {
		t.Errorf("page_count in pagination response should be 10; got %v", res.PageCount)
	}
}

func TestGetRequestData_with_data(t *testing.T) {
	values := url.Values{}
	values.Add("page", "1")
	values.Add("page_size", "20")
	res := getRequestData(values)
	if res.pageSize != 20 {
		t.Errorf("page_size in pagination response should be %v; got %v", 20, res.pageSize)
	}
	if res.pageNumber != 1 {
		t.Errorf("page_number in pagination response should be %v; got %v", 1, res.pageNumber)
	}
}

func TestGetRequestData_with_no_data(t *testing.T) {
	values := url.Values{}
	res := getRequestData(values)
	if res.pageSize != PageSize {
		t.Errorf("page_size in pagination response should be %v; got %v", PageSize, res.pageSize)
	}
	if res.pageNumber != 1 {
		t.Errorf("page_number in pagination response should be %v; got %v", 1, res.pageNumber)
	}
}

func TestCreateWhereClause(t *testing.T) {
	colNames := []string{"name", "system", "age"}
	param1 := parameter{"name", "=", "otto"}
	param2 := parameter{"system", "=", "hipaca"}
	param3 := parameter{"age", "=", "33"}
	params := parameters{param1, param2, param3}
	c := make(chan whereClause)
	go createWhereClause(colNames, params, c)
	where := <-c
	if !where.exists {
		t.Errorf("where clauses should exists; got %v", where.exists)
	}
	expectedCLAUSE := " WHERE name = $%v AND system = $%v AND age = $%v"
	if where.clause != expectedCLAUSE {
		t.Errorf("where clauses should be %v; got %v", expectedCLAUSE, where.clause)
	}
	rightARGS := []interface{}{"otto", "hipaca", "33"}
	for i := 0; i < len(where.args); i++ {
		if where.args[i] != rightARGS[i] {
			t.Errorf("where clause arg number %v should be %v; got %v", i, rightARGS[i], where.args[i])
		}
	}
}

func TestCreateWhereClause_with_filters(t *testing.T) {
	colNames := []string{"age", "skills", "cars"}
	param1 := parameter{"age", ">", "33"}
	param2 := parameter{"skills", "<>", "golang"}
	param3 := parameter{"cars", ">=", "2"}
	param4 := parameter{"cars", ">", "4"}
	param5 := parameter{"cars", "<", "4"}
	param6 := parameter{"cars", "<=", "5"}
	params := parameters{param1, param2, param3, param4, param5, param6}
	c := make(chan whereClause)
	go createWhereClause(colNames, params, c)
	where := <-c
	if !where.exists {
		t.Errorf("where clauses should exists; got %v", where.exists)
	}
	expectedCLAUSE := " WHERE age > $%v AND skills <> $%v AND cars >= $%v AND cars > $%v AND cars < $%v AND cars <= $%v"
	if where.clause != expectedCLAUSE {
		t.Errorf("filter clause should be %v; got %v", expectedCLAUSE, where.clause)
	}
	rightARGS := []interface{}{"33", "golang", "2", "4", "4", "5"}
	for i := 0; i < len(where.args); i++ {
		if where.args[i] != rightARGS[i] {
			t.Errorf("where clause arg number %v should be %v; got %v", i, rightARGS[i], where.args[i])
		}
	}
}

func TestCreatePaginationClause_with_page_gt_1(t *testing.T) {
	pageNumber := 2
	pageSize := 30
	c := make(chan string)
	go createPaginationClause(pageNumber, pageSize, c)
	clause := <-c
	expectedCLAUSE := " LIMIT 30 OFFSET 30"
	if clause != expectedCLAUSE {
		t.Errorf("expected clause should be %v; got %v", expectedCLAUSE, clause)
	}
}

func TestCreatePaginationClause_with_page_eq_1(t *testing.T) {
	pageNumber := 1
	pageSize := 30
	c := make(chan string)
	go createPaginationClause(pageNumber, pageSize, c)
	clause := <-c
	expectedCLAUSE := " LIMIT 30 OFFSET 0"
	if clause != expectedCLAUSE {
		t.Errorf("expected clause should be %v; got %v", expectedCLAUSE, clause)
	}
}

func TestCreatePaginationClause_with_page_lt_1(t *testing.T) {
	pageNumber := -4
	pageSize := 30
	c := make(chan string)
	go createPaginationClause(pageNumber, pageSize, c)
	clause := <-c
	expectedCLAUSE := " LIMIT 30 OFFSET 0"
	if clause != expectedCLAUSE {
		t.Errorf("expected clause should be %v; got %v", expectedCLAUSE, clause)
	}
}

func TestCreateOrderByClause_with_sorting_options(t *testing.T) {
	colNames := []string{"name", "lastname", "age", "address"}
	params := parameters{{"sort", "=", "+name,-lastname,-age,+address"}}
	c := make(chan string)
	go createOrderByClause(params, colNames, c)
	clause := <-c
	expectedCLAUSE := " ORDER BY name ASC,lastname DESC,age DESC,address ASC,id"
	if clause != expectedCLAUSE {
		t.Errorf("expected clause should be %v, got %v", expectedCLAUSE, clause)
	}
}

func TestCreateOrderByClause_with_no_sorting_options(t *testing.T) {
	colNames := []string{"name", "lastname", "age", "address"}
	params := parameters{}
	c := make(chan string)
	go createOrderByClause(params, colNames, c)
	clause := <-c
	expectedCLAUSE := " ORDER BY id"
	if clause != expectedCLAUSE {
		t.Errorf("expected clause should be %v, got %v", expectedCLAUSE, clause)
	}
}
