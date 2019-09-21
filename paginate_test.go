package paginate

import (
	"fmt"
	"log"
	"net/url"
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