# paginate 
[![Build Status](https://travis-ci.org/ottotech/paginate.svg?branch=master)](https://travis-ci.org/ottotech/paginate)
[![MIT license](http://img.shields.io/badge/license-MIT-brightgreen.svg)](http://opensource.org/licenses/MIT)
[![GoDoc](https://godoc.org/github.com/ottotech/paginate?status.svg)](https://godoc.org/github.com/ottotech/paginate)

## Overview

Package **paginate** provides basic pagination capabilities to paginate sql database tables with Go.

## How to use?

```go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/ottotech/paginate"
	"log"
	"net/url"
	"time"
)

type Employee struct {
	ID           int                  `paginate:"id;col=id;param=employee_id"`
	Name         string               `paginate:"filter;col=name"`
	LastName     string               `paginate:"filter;col=last_name"`
	WorkerNumber paginate.NullInt     `paginate:"col=worker_number"`
	DateJoined   time.Time            `paginate:"col=date_joined"`
	Salary       float64              `paginate:"filter;col=salary"`
	NullText     paginate.NullString  `paginate:"col=null_text"`
	NullVarchar  paginate.NullString  `paginate:"col=null_varchar"`
	NullBool     paginate.NullBool    `paginate:"col=null_bool"`
	NullDate     paginate.NullTime    `paginate:"col=null_date"`
	NullInt      paginate.NullInt     `paginate:"col=null_int"`
	NullFloat    paginate.NullFloat64 `paginate:"col=null_float"`
}

func main() {
    dbUri := "user=root password=secret host=localhost " +
        "port=5432 dbname=events_db sslmode=disable"
    
    db, err := sql.Open("postgres", dbUri)
    if err != nil {
        panic(err)
    }
    
    if err = db.Ping(); err != nil {
        panic(err)
    }
    fmt.Println("You connected to your database.")
    
    u, err := url.Parse("http://localhost?name=Ringo")
    if err != nil {
        // Handle error gracefully...
    }
    
    paginator, err := paginate.NewPaginator(Employee{}, "postgres", *u, paginate.TableName("employees"))
    if err != nil {
        // Handle error gracefully...
    }
    
    cmd, args, err := paginator.Paginate()
    if err != nil {
        // Handle error gracefully...
    }
    
    rows, err := db.Query(cmd, args...)
    if err != nil {
        // Handle error gracefully...
    }
    defer rows.Close()
    
    for rows.Next() {
        err = rows.Scan(paginator.GetRowPtrArgs()...)
        if err != nil {
            // Handle error gracefully...
        }
    }
    
    if err = rows.Err(); err != nil {
        // Handle error gracefully...
    }
    
    results := make([]Employee, 0)
    
    for paginator.NextData() {
        employee := Employee{}
        err = paginator.Scan(&employee)
        if err != nil {
            // Handle error gracefully...
        }
        results = append(results, employee)
    }
    
    // You should be able to see the paginated data inside results. 
    fmt.Println(results)
    
    // You should be able to serialize into json the paginated data.
    sb, err := json.Marshal(results)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(sb)
}
```

## Contributing

## Notes

Paginator does not take into consideration performance since it uses the OFFSET sql argument
which reads and counts all rows from the beginning until it reaches the requested page. For
not too big datasets Paginator will just work fine. If you care about performance because you
are dealing with heavy data you might want to write a custom solution for that.

## License
Released under MIT license, see [LICENSE](https://github.com/ottotech/paginate/blob/master/LICENSE.md) for details.