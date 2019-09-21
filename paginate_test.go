package paginate

import (
	"fmt"
	"log"
	"net/url"
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
	expectedSQL := "SELECT column1, column2, column3, column4, column5, count(*) over() FROM test WHERE column5 = $1 AND column1 > $2 AND column2 < $3 AND column3 >= $4 AND column4 <= $5 ORDER BY column1 ASC,column2 DESC,id LIMIT 15 OFFSET 0"
	paginator := NewPaginator(tableName, colNames, *u)
	paginator.SetPageCount(15)
	paginator.SetTotalResult(30)
	sql, args, err := paginator.Paginate()
	if sql != expectedSQL {
		t.Errorf("expected sql:\n %v; \ngot %v\n", expectedSQL, sql)
	}
	if args[0] != "otto" {
		t.Errorf("arg $1 should be otto; got %v", args[0])
	}
	if args[1] != 2 {
		t.Errorf("arg $2 should be 2; got %v", args[1])
	}
	if args[2] != 4 {
		t.Errorf("arg $3 should be 4; got %v", args[2])
	}
	if args[3] != 40 {
		t.Errorf("arg $4 should be 40; got %v", args[3])
	}
	if args[4] != 7 {
		t.Errorf("arg $5 should be 7; got %v", args[4])
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
	expectedSQL := "SELECT column1, column2, column3, column4, column5, count(*) over() FROM test WHERE column5 = $1 AND column1 > $2 AND column2 < $3 AND column3 >= $4 AND column4 <= $5 ORDER BY column1 ASC,column2 DESC,id LIMIT 10 OFFSET 10"
	paginator := NewPaginatorWithLimit(10, tableName, colNames, *u)
	paginator.SetPageCount(10)
	paginator.SetTotalResult(30)
	sql, args, err := paginator.Paginate()
	if sql != expectedSQL {
		t.Errorf("expected sql:\n %v; \ngot %v\n", expectedSQL, sql)
	}
	if args[0] != "otto" {
		t.Errorf("arg $1 should be otto; got %v", args[0])
	}
	if args[1] != 2 {
		t.Errorf("arg $2 should be 2; got %v", args[1])
	}
	if args[2] != 4 {
		t.Errorf("arg $3 should be 4; got %v", args[2])
	}
	if args[3] != 40 {
		t.Errorf("arg $4 should be 40; got %v", args[3])
	}
	if args[4] != 7 {
		t.Errorf("arg $5 should be 7; got %v", args[4])
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