package paginate

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
	"time"
)

// variables that hold the global access to the psql and mysql testing databases.
var (
	mysqlTestDB *sql.DB
	psqlTestDB  *sql.DB
)

// variables that represent the uris to connect to the psql and mysql testing databases.
var (
	mysqlDatabaseUri = "root:secret@tcp(localhost:3306)/%s?multiStatements=true&parseTime=true"
	psqlDatabaseUri  = "user=postgres password=secret host=localhost port=5432 dbname=%s sslmode=disable"
)

func createMysqlDatabaseAndTestingTable() error {
	defaultDB, err := sql.Open("mysql", fmt.Sprintf(mysqlDatabaseUri, "mysql"))
	if err != nil {
		return err
	}

	_, err = defaultDB.Exec(fmt.Sprintf("CREATE DATABASE %s;", "paginate_test"))
	if err != nil {
		return err
	}
	defaultDB.Close()

	db, err := sql.Open("mysql", fmt.Sprintf(mysqlDatabaseUri, "paginate_test"))
	if err != nil {
		return err
	}
	mysqlTestDB = db

	createTableQuery := `
CREATE TABLE employees
  (
     id            INT auto_increment PRIMARY KEY,
     name          VARCHAR(200) NOT NULL,
     last_name     VARCHAR(200) NOT NULL,
     worker_number INT NOT NULL,
     date_joined   TIMESTAMP NULL,
     salary        FLOAT NULL,
     null_text     TEXT NULL,
     null_varchar  VARCHAR(100) NULL,
     null_bool     TINYINT(1) NULL,
	 null_date     TIMESTAMP NULL,
	 null_int      INT NULL,
	 null_float    FLOAT NULL,
     CONSTRAINT employee_worker_number_uindex UNIQUE (worker_number)
  );
`

	ctx := context.Background()
	tx, err := mysqlTestDB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	for _, q := range []string{createTableQuery} {
		_, err = tx.Exec(q)

		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return err
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func createPsqlDatabaseAndTestingTable() error {
	defaultDB, err := sql.Open("postgres", fmt.Sprintf(psqlDatabaseUri, "postgres"))
	if err != nil {
		return err
	}
	_, err = defaultDB.Exec("CREATE DATABASE paginate_test;")
	if err != nil {
		return err
	}
	defaultDB.Close()

	db, err := sql.Open("postgres", fmt.Sprintf(psqlDatabaseUri, "paginate_test"))
	if err != nil {
		return err
	}
	psqlTestDB = db

	createTableQuery := `
create table employees
(
   id            serial       not null
       constraint employees_pk
           primary key,
   name          varchar(200) not null,
   last_name     varchar(200) not null,
   worker_number integer      not null,
   date_joined   timestamp with time zone,
   salary        double precision,
   null_text     text,
   null_varchar  varchar(100),
   null_bool     boolean,
   null_date     timestamp with time zone,
   null_int      integer,
   null_float    double precision
);

create unique index employees_id_uindex
   on employees (id);

create unique index employees_worker_number_uindex
   on employees (worker_number);
`

	ctx := context.Background()
	tx, err := psqlTestDB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	for _, q := range []string{createTableQuery} {
		_, err = tx.Exec(q)

		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return err
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func addDataToDatabaseTable(db *sql.DB, dialect string) error {
	type employee struct {
		Name, LastName string
		WorkNumber     int
		DateJoined     time.Time
		Salary         float64
		NullBool       interface{}
	}

	employees := []employee{
		{
			Name:       "Ringo",
			LastName:   "Star",
			WorkNumber: 1,
			DateJoined: time.Now(),
			Salary:     5400,
		},
		{
			Name:       "Bill",
			LastName:   "Gates",
			WorkNumber: 2,
			DateJoined: time.Now(),
			Salary:     1200000,
			NullBool:   true,
		},
		{
			Name:       "Mark",
			LastName:   "Smith",
			WorkNumber: 3,
			DateJoined: time.Now(),
			Salary:     8000,
		},
		{
			Name:       "John",
			LastName:   "Smith",
			WorkNumber: 4,
			DateJoined: time.Now(),
			Salary:     4650.90,
		},
		{
			Name:       "Fred",
			LastName:   "Smith",
			WorkNumber: 5,
			DateJoined: time.Now(),
			Salary:     7550,
			NullBool:   true,
		},
		{
			Name:       "Rob",
			LastName:   "Williams",
			WorkNumber: 6,
			DateJoined: time.Now(),
			Salary:     9880,
		},
		{
			Name:       "Juliana",
			LastName:   "Collier",
			WorkNumber: 7,
			DateJoined: time.Now(),
			Salary:     7788,
		},
		{
			Name:       "Erika",
			LastName:   "Smith",
			WorkNumber: 8,
			DateJoined: time.Now(),
			Salary:     4455,
		},
		{
			Name:       "Maria",
			LastName:   "Gomez",
			WorkNumber: 9,
			DateJoined: time.Now(),
			Salary:     7550,
		},
		{
			Name:       "Rafael",
			LastName:   "Smith",
			WorkNumber: 10,
			DateJoined: time.Now(),
			Salary:     7550,
		},
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	var sqlStatement string

	if dialect == "mysql" {
		sqlStatement = `
		INSERT INTO employees (name, last_name, worker_number, date_joined, salary, null_bool)
		VALUES (?, ?, ?, ?, ?, ?)`
	} else if dialect == "postgres" {
		sqlStatement = `
		INSERT INTO employees (name, last_name, worker_number, date_joined, salary, null_bool)
		VALUES ($1, $2, $3, $4, $5, $6)`
	}

	for _, e := range employees {
		_, err := tx.Exec(sqlStatement, e.Name, e.LastName, e.WorkNumber, e.DateJoined, e.Salary, e.NullBool)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return err
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func removeMysqlDatabase() error {
	db, err := sql.Open("mysql", fmt.Sprintf(mysqlDatabaseUri, "mysql"))
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("DROP DATABASE %s;", "paginate_test"))
	if err != nil {
		return err
	}
	return nil
}

func removePsqlDatabase() error {
	db, err := sql.Open("postgres", fmt.Sprintf(psqlDatabaseUri, "postgres"))
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("DROP DATABASE %s;", "paginate_test"))
	if err != nil {
		return err
	}
	return nil
}

func TestMain(m *testing.M) {
	var err error
	err = createMysqlDatabaseAndTestingTable()
	if err != nil {
		removeMysqlDatabase()
		log.Fatalln(err)
	}

	err = addDataToDatabaseTable(mysqlTestDB, "mysql")
	if err != nil {
		removeMysqlDatabase()
		log.Fatalln(err)
	}

	err = createPsqlDatabaseAndTestingTable()
	if err != nil {
		removePsqlDatabase()
		log.Fatalln(err)
	}

	err = addDataToDatabaseTable(psqlTestDB, "postgres")
	if err != nil {
		removePsqlDatabase()
		log.Fatalln(err)
	}

	code := m.Run()

	mysqlTestDB.Close()
	psqlTestDB.Close()

	err = removePsqlDatabase()
	if err != nil {
		log.Println(err)
	}

	err = removeMysqlDatabase()
	if err != nil {
		log.Println(err)
	}

	os.Exit(code)
}
