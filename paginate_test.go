package paginate

import (
	"fmt"
	"log"
	"net/url"
	"testing"
)

func ExampleNewPaginator_1() {
	type Application struct {
		ID     int    `paginate:"id;col=id"`
		System string `paginate:"filter"`
	}
	opt := TableName("test")
	app := Application{}
	u, err := url.Parse("http://ottotech.com?system=platform")
	if err != nil {
		log.Fatal(err)
	}
	paginator, err := NewPaginator(app, *u, opt)
	if err != nil {
		log.Fatalln(err)
	}
	sql, args, err := paginator.Paginate()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(sql)
	fmt.Printf("arg $1: %v\n", args[0])
	// Output:
	// SELECT id, system, count(*) over() FROM test WHERE system = $1 ORDER BY id LIMIT 30 OFFSET 0
	// arg $1: platform
}

func ExampleNewPaginator_2() {
	type Person struct {
		ID       int    `paginate:"id"`
		Name     string `paginate:"filter"`
		LastName string `paginate:"filter"`
	}
	u, err := url.Parse("http://ottotech.com?name=Ringo&last_name=Star")
	if err != nil {
		log.Fatal(err)
	}
	paginator, err := NewPaginator(Person{}, *u)
	if err != nil {
		log.Fatalln(err)
	}
	sql, args, err := paginator.Paginate()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("arg $1: %v\n", args[0])
	fmt.Printf("arg $2: %v\n", args[1])
	fmt.Println(sql)
	// Unordered output:
	// SELECT id, name, last_name, count(*) over() FROM person WHERE name = $1 AND last_name = $2 ORDER BY id LIMIT 30 OFFSET 0
	// arg $1: Ringo
	// arg $2: Star
}

func ExampleNewPaginator_3() {
	type Person struct {
		ID       int    `paginate:"id"`
		Name     string `paginate:"filter"`
		LastName string `paginate:"filter"`
	}
	u, err := url.Parse("http://ottotech.com?name=Ringo&page_size=8&last_name<>Star")
	if err != nil {
		log.Fatal(err)
	}
	paginator, err := NewPaginator(Person{}, *u)
	if err != nil {
		log.Fatalln(err)
	}
	sql, args, err := paginator.Paginate()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(sql)
	fmt.Printf("arg $1: %v\n", args[0])
	fmt.Printf("arg $2: %v\n", args[1])
	// Output:
	// SELECT id, name, last_name, count(*) over() FROM person WHERE name = $1 AND last_name <> $2 ORDER BY id LIMIT 8 OFFSET 0
	// arg $1: Ringo
	// arg $2: Star
}

func ExampleNewPaginator_4() {
	type Person struct {
		ID       int    `paginate:"id"`
		Name     string `paginate:"filter"`
		LastName string `paginate:"filter"`
	}
	u, err := url.Parse("http://ottotech.com?name=Ringo&sort=+name,-last_name&page_size=7")
	if err != nil {
		log.Fatal(err)
	}
	paginator, err := NewPaginator(Person{}, *u)
	if err != nil {
		log.Fatalln(err)
	}
	sql, args, err := paginator.Paginate()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(sql)
	fmt.Printf("arg $1: %v\n", args[0])
	// Output:
	// SELECT id, name, last_name, count(*) over() FROM person WHERE name = $1 ORDER BY name ASC,last_name DESC,id LIMIT 7 OFFSET 0
	// arg $1: Ringo
}

func ExampleNewPaginator_5() {
	type Person struct {
		ID       int     `paginate:"id"`
		Name     string  `paginate:"filter"`
		LastName string  `paginate:"filter"`
		Age      int     `paginate:"filter"`
		Salary   float64 `paginate:"filter"`
	}
	u, err := url.Parse("http://ottotech.com?age>20&salary<80000")
	if err != nil {
		log.Fatal(err)
	}
	paginator, err := NewPaginator(Person{}, *u)
	if err != nil {
		log.Fatalln(err)
	}
	sql, args, err := paginator.Paginate()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(sql)
	fmt.Printf("arg $1: %v\n", args[0])
	fmt.Printf("arg $2: %v\n", args[1])
	// Output:
	// SELECT id, name, last_name, age, salary, count(*) over() FROM person WHERE age > $1 AND salary < $2 ORDER BY id LIMIT 30 OFFSET 0
	// arg $1: 20
	// arg $2: 80000
}

func ExampleNewPaginatorWithPageSize15() {
	type Person struct {
		ID       int    `paginate:"id"`
		Name     string `paginate:"filter"`
		LastName string `paginate:"filter"`
	}
	u, err := url.Parse("http://ottotech.com?name=Ringo")
	if err != nil {
		log.Fatal(err)
	}
	paginator, err := NewPaginator(Person{}, *u, PageSize(15))
	if err != nil {
		log.Fatalln(err)
	}
	sql, args, err := paginator.Paginate()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sql)
	fmt.Printf("arg $1: %v\n", args[0])
	// Output:
	// SELECT id, name, last_name, count(*) over() FROM person WHERE name = $1 ORDER BY id LIMIT 15 OFFSET 0
	// arg $1: Ringo
}

func ExampleTestINSqlClause() {
	type Person struct {
		ID       int    `paginate:"id"`
		Name     string `paginate:"filter"`
		LastName string `paginate:"filter"`
	}
	u, err := url.Parse("http://ottotech.com?name=Ringo&name=Rob")
	if err != nil {
		log.Fatal(err)
	}
	paginator, err := NewPaginator(Person{}, *u)
	if err != nil {
		log.Fatalln(err)
	}
	sql, args, err := paginator.Paginate()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sql)
	fmt.Printf("arg $1: %v\n", args[0])
	fmt.Printf("arg $2: %v\n", args[1])
	// Output:
	// SELECT id, name, last_name, count(*) over() FROM person WHERE name IN($1,$2) ORDER BY id LIMIT 30 OFFSET 0
	// arg $1: Ringo
	// arg $2: Rob
}

func ExampleTestNotINSqlClause() {
	type Person struct {
		ID       int    `paginate:"id"`
		Name     string `paginate:"filter"`
		LastName string `paginate:"filter"`
	}
	u, err := url.Parse("http://ottotech.com?name<>Ringo&name<>Rob")
	if err != nil {
		log.Fatal(err)
	}
	paginator, err := NewPaginator(Person{}, *u)
	if err != nil {
		log.Fatalln(err)
	}
	sql, args, err := paginator.Paginate()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sql)
	fmt.Printf("arg $1: %v\n", args[0])
	fmt.Printf("arg $2: %v\n", args[1])
	// Output:
	// SELECT id, name, last_name, count(*) over() FROM person WHERE name NOT IN($1,$2) ORDER BY id LIMIT 30 OFFSET 0
	// arg $1: Ringo
	// arg $2: Rob
}

func ExampleTestNOFilterTag() {
	type Person struct {
		ID       int `paginate:"id"`
		Name     string
		LastName string
	}
	// We will try to filter the database table by name and last_name.
	// However, this shouldn't work since the "filter" tag is not set in
	// any of the fields in the Person struct.
	u, err := url.Parse("http://ottotech.com?name=Ringo&last_name=Star")
	if err != nil {
		log.Fatal(err)
	}
	paginator, err := NewPaginator(Person{}, *u)
	if err != nil {
		log.Fatalln(err)
	}
	sql, args, err := paginator.Paginate()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sql)
	fmt.Printf("args length: %v\n", len(args))
	// Output:
	// SELECT id, name, last_name, count(*) over() FROM person ORDER BY id LIMIT 30 OFFSET 0
	// args length: 0
}

func TestCreateWhereClauseMultipleFilters(t *testing.T) {
	colNames := []string{"age", "skills", "cars"}
	param1 := parameter{"age", ">", "33"}
	param2 := parameter{"skills", "<>", "golang"}
	param3 := parameter{"cars", ">=", "2"}
	param4 := parameter{"cars", ">", "4"}
	param5 := parameter{"cars", "<", "4"}
	param6 := parameter{"cars", "<=", "5"}
	params := parameters{param1, param2, param3, param4, param5, param6}
	c := make(chan whereClause)
	go createWhereClause(colNames, params, []RawWhereClause{}, c)
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
	colNames := []string{"id", "name", "lastname", "age", "address"}
	params := parameters{{"sort", "=", "+name,-lastname,-age,+address"}}
	c := make(chan string)
	go createOrderByClause(params, colNames, "id", c)
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
	go createOrderByClause(params, colNames, "id", c)
	clause := <-c
	expectedCLAUSE := " ORDER BY id"
	if clause != expectedCLAUSE {
		t.Errorf("expected clause should be %v, got %v", expectedCLAUSE, clause)
	}
}
