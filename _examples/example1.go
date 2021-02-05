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

type HistoryEvent struct {
	Id               int       `json:"id" paginate:"id;col=id"`
	Performer        string    `json:"performer" paginate:"col=performer;param=performera"`
	Player           string    `json:"player" paginate:"col=player"`
	System           string    `json:"system" paginate:"col=system"`
	Event            string    `json:"event" paginate:"col=event"`
	DateCreated      time.Time `json:"date_created" paginate:"col=date_created"`
	ObjectIdentifier string    `json:"object_identifier" paginate:"col=object_identifier"`
	Notes            string    `json:"notes" paginate:"col=notes"`
	Dummy            int       `json:"dummy" paginate:"col=dummy"`
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

	u, err := url.Parse("http://localhost")
	if err != nil {
		log.Fatal(err)
	}

	paginator, err := paginate.NewPaginator(
		HistoryEvent{},
		*u,
		paginate.TableName("events"),
	)
	if err != nil {
		log.Fatal(err)
	}

	cmd, args, err := paginator.Paginate()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cmd, args)

	rows, err := db.Query(cmd, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(paginator.GetRowPtrArgs()...)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	result := make([]HistoryEvent, 0)
	for paginator.NextData() {
		event := HistoryEvent{}
		err = paginator.Scan(&event)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, event)
	}

	fmt.Printf("This is the result: %+v\n", result)
	_, err = json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(string(sb))
	fmt.Printf("%+v\n", paginator.Response())
}
