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
	createTable2Query := `
CREATE TABLE developer
  (
     employee_id          INT NOT NULL,
     programming_language VARCHAR(200) NOT NULL,
     FOREIGN KEY (employee_id) REFERENCES employees (id) ON DELETE CASCADE
  );

CREATE UNIQUE INDEX developer_employee_id_uindex ON developer (employee_id); 
`
	createTable3Query := `
CREATE TABLE manager
  (
     employee_id INT NOT NULL,
     FOREIGN KEY (employee_id) REFERENCES employees (id) ON DELETE CASCADE
  );

CREATE UNIQUE INDEX manager_employee_id_uindex ON manager (employee_id); 
`

	ctx := context.Background()
	tx, err := mysqlTestDB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	for _, q := range []string{createTableQuery, createTable2Query, createTable3Query} {
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
CREATE TABLE employees
  (
     id            SERIAL NOT NULL CONSTRAINT employees_pk PRIMARY KEY,
     NAME          VARCHAR(200) NOT NULL,
     last_name     VARCHAR(200) NOT NULL,
     worker_number INTEGER NOT NULL,
     date_joined   TIMESTAMP WITH time zone,
     salary        DOUBLE PRECISION,
     null_text     TEXT,
     null_varchar  VARCHAR(100),
     null_bool     BOOLEAN,
     null_date     TIMESTAMP WITH time zone,
     null_int      INTEGER,
     null_float    DOUBLE PRECISION
  );

CREATE UNIQUE INDEX employees_id_uindex
  ON employees (id);

CREATE UNIQUE INDEX employees_worker_number_uindex
  ON employees (worker_number); 
`

	create2TableQuery := `
CREATE TABLE developer
  (
     employee_id          BIGINT NOT NULL CONSTRAINT developer_employees_id_fk
     REFERENCES
     employees ON DELETE
     CASCADE,
     programming_language VARCHAR(200) NOT NULL
  );

CREATE UNIQUE INDEX developer_employee_id_uindex ON developer (employee_id);
`

	create3TableQuery := `
CREATE TABLE manager
  (
     employee_id BIGINT NOT NULL CONSTRAINT manager_employees_id_fk REFERENCES
     employees ON DELETE
     CASCADE
  );

CREATE UNIQUE INDEX manager_employee_id_uindex ON manager (employee_id); 
`

	ctx := context.Background()
	tx, err := psqlTestDB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	for _, q := range []string{createTableQuery, create2TableQuery, create3TableQuery} {
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

type employee struct {
	Name, LastName string
	WorkNumber     int
	DateJoined     time.Time
	Salary         float64
	NullBool       interface{}
}

// We are adding 10 employees into the "employees" table.
// The first 5 employees on the list will be developers
// 3 of them will know Go and the rest 2 know Python as a programming
// language. The rest 5 employees will be managers.
var employees = []employee{
	// The next two employees
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
	// Below all employees are manager.
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

func addDataToDatabaseTables(db *sql.DB, dialect string) error {
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
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	}

	firstFiveIDs := make([]int64, 0, 5)
	lastFiveIDs := make([]int64, 0, 5)

	for i, e := range employees {
		var lastInsertedID int64

		switch dialect {
		case "mysql":
			res, err := tx.Exec(sqlStatement, e.Name, e.LastName, e.WorkNumber, e.DateJoined, e.Salary, e.NullBool)
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					return err
				}
				return err
			}

			id, err := res.LastInsertId()
			if err != nil {
				return err
			}
			lastInsertedID = id
		case "postgres":
			err := tx.QueryRow(sqlStatement, e.Name, e.LastName, e.WorkNumber, e.DateJoined, e.Salary, e.NullBool).Scan(&lastInsertedID)
			if err != nil {
				return err
			}
		}

		if i < 5 {
			firstFiveIDs = append(firstFiveIDs, lastInsertedID)
		} else {
			lastFiveIDs = append(lastFiveIDs, lastInsertedID)
		}
	}

	// Let's create 5 developers now. The first 3 will be Go developers.
	// The rest 2 will be Python developers.
	for i := 0; i < len(firstFiveIDs); i++ {
		var programmingLanguage string

		if i < 3 {
			programmingLanguage = "Go"
		} else {
			programmingLanguage = "Python"
		}

		switch dialect {
		case "mysql":
			_, err := tx.Exec("INSERT INTO developer (employee_id, programming_language) VALUES (?, ?);", firstFiveIDs[i], programmingLanguage)
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					return err
				}
				return err
			}
		case "postgres":
			_, err := tx.Exec("INSERT INTO developer (employee_id, programming_language) VALUES ($1, $2);", firstFiveIDs[i], programmingLanguage)
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					return err
				}
				return err
			}
		}
	}

	// Finally let's create five managers.
	for _, id := range lastFiveIDs {
		switch dialect {
		case "mysql":
			_, err := tx.Exec("INSERT INTO manager (employee_id) VALUES (?);", id)
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					return err
				}
				return err
			}
		case "postgres":
			_, err := tx.Exec("INSERT INTO manager (employee_id) VALUES ($1);", id)
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					return err
				}
				return err
			}
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

	err = addDataToDatabaseTables(mysqlTestDB, "mysql")
	if err != nil {
		removeMysqlDatabase()
		log.Fatalln(err)
	}

	err = createPsqlDatabaseAndTestingTable()
	if err != nil {
		removePsqlDatabase()
		log.Fatalln(err)
	}

	err = addDataToDatabaseTables(psqlTestDB, "postgres")
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
