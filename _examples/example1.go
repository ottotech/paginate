package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/ottotech/paginate"
	"log"
	"net/url"
	"time"
)

type HistoryEvent struct {
	Id               int       `json:"id" paginate:"id"`
	Performer        string    `json:"performer"`
	Player           string    `json:"player"`
	System           string    `json:"system"`
	Event            string    `json:"event"`
	DateCreated      time.Time `json:"date_created"`
	ObjectIdentifier string    `json:"object_identifier"`
	Notes            string    `json:"notes"`
}

func main() {
	dbUrl := "user=ottosg password=secret host=localhost " +
		"port=8888 dbname=events_db sslmode=disable"

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You connected to your database.")

	u, err := url.Parse("http://localhost?system=olms")
	if err != nil {
		log.Fatal(err)
	}

	paginator, err := paginate.NewPaginator(HistoryEvent{}, "events", *u)
	if err != nil {
		log.Fatal(err)
	}

	cmd, args, err := paginator.Paginate()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query(cmd, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	n := 0
	for rows.Next() {

		err = rows.Scan(paginator.GetRowPtrArgs()...)
		if err != nil {
			log.Fatal(err)
		}
		n++
		fmt.Println(n)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	result := make([]HistoryEvent, 0)
	for paginator.Next() {
		event := HistoryEvent{}
		err = paginator.Scan(&event)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, event)
	}

	fmt.Printf("This is the result: %+v\n", result)
}
