package paginate

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
	"time"
)

func TestPaginatorMysql_HappyPath(t *testing.T) {
	type Employee struct {
		ID           int         `paginate:"id;col=id"`
		Name         string      `paginate:"col=name"`
		LastName     string      `paginate:"col=last_name"`
		WorkerNumber NullInt     `paginate:"col=worker_number"`
		DateJoined   time.Time   `paginate:"col=date_joined"`
		Salary       float64     `paginate:"col=salary"`
		NullText     NullString  `paginate:"col=null_text"`
		NullVarchar  NullString  `paginate:"col=null_varchar"`
		NullBool     NullBool    `paginate:"col=null_bool"`
		NullDate     NullTime    `paginate:"col=null_date"`
		NullInt      NullInt     `paginate:"col=null_int"`
		NullFloat    NullFloat64 `paginate:"col=null_float"`
	}

	u, err := url.Parse("http://localhost")
	if err != nil {
		t.Fatal(err)
	}

	pag, err := NewPaginator(Employee{}, "mysql", *u, TableName("employees"))
	if err != nil {
		t.Fatal(err)
	}

	cmd, args, err := pag.Paginate()
	if err != nil {
		t.Fatal(err)
	}

	rows, err := mysqlTestDB.Query(cmd, args...)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(pag.GetRowPtrArgs()...)
		if err != nil {
			t.Fatal(err)
		}
	}

	if err = rows.Err(); err != nil {
		t.Fatal(err)
	}

	results := make([]Employee, 0)

	for pag.NextData() {
		employee := Employee{}
		err = pag.Scan(&employee)
		if err != nil {
			t.Fatal(err)
		}
		results = append(results, employee)
	}

	if len(results) != 10 {
		t.Errorf("we should have 10 records in result; got %d", len(results))
	}

	countNullBooleans, countNullStrings, countNullDates, countNullInts, countNullFloats := 0, 0, 0, 0, 0
	for _, r := range results {
		if !r.NullBool.Valid {
			countNullBooleans++
		}
		if !r.NullVarchar.Valid {
			countNullStrings++
		}
		if !r.NullDate.Valid {
			countNullDates++
		}
		if !r.NullInt.Valid {
			countNullInts++
		}
		if !r.NullFloat.Valid {
			countNullFloats++
		}
	}

	if countNullBooleans != 8 {
		t.Errorf("we should have 8 nulls for the NullBool field")
	}

	if countNullStrings != 10 {
		t.Errorf("we should have 10 nulls for the NullVarchar field")
	}

	if countNullDates != 10 {
		t.Errorf("we should have 10 nulls for the NullDate field")
	}

	if countNullInts != 10 {
		t.Errorf("we should have 10 nulls for the NullInt field")
	}

	if countNullFloats != 10 {
		t.Errorf("we should have 10 nulls for the NullFloat field")
	}
}

func TestPaginatorMysql_JsonMarshalling(t *testing.T) {
	type Employee struct {
		ID           int         `json:"id" paginate:"id;col=id"`
		Name         string      `json:"name" paginate:"col=name"`
		LastName     string      `json:"last_name" paginate:"col=last_name"`
		WorkerNumber NullInt     `json:"worker_number" paginate:"col=worker_number"`
		DateJoined   time.Time   `json:"date_joined" paginate:"col=date_joined"`
		Salary       float64     `json:"salary" paginate:"col=salary"`
		NullText     NullString  `json:"null_text" paginate:"col=null_text"`
		NullVarchar  NullString  `json:"null_varchar" paginate:"col=null_varchar"`
		NullBool     NullBool    `json:"null_bool" paginate:"col=null_bool"`
		NullDate     NullTime    `json:"null_date" paginate:"col=null_date"`
		NullInt      NullInt     `json:"null_int" paginate:"col=null_int"`
		NullFloat    NullFloat64 `json:"null_float" paginate:"col=null_float"`
	}

	u, err := url.Parse("http://localhost")
	if err != nil {
		t.Fatal(err)
	}

	pag, err := NewPaginator(Employee{}, "mysql", *u, TableName("employees"))
	if err != nil {
		t.Fatal(err)
	}

	cmd, args, err := pag.Paginate()
	if err != nil {
		t.Fatal(err)
	}

	rows, err := mysqlTestDB.Query(cmd, args...)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(pag.GetRowPtrArgs()...)
		if err != nil {
			t.Fatal(err)
		}
	}

	if err = rows.Err(); err != nil {
		t.Fatal(err)
	}

	results := make([]Employee, 0)

	for pag.NextData() {
		employee := Employee{}
		err = pag.Scan(&employee)
		if err != nil {
			t.Fatal(err)
		}
		results = append(results, employee)
	}

	_, err = json.Marshal(results)
	if err != nil {
		t.Errorf("we should be able to marshal results into json; got err %s", err)
	}
}

// This test will check if the paginator can infer properly
// the name of the target database table and its columns
// from the given struct.
func TestNewPaginatorMysql_DefaultTableAndColumnsInferring(t *testing.T) {
	type Employees struct {
		ID           int `paginate:"id"` // this is a mandatory tag.
		Name         string
		LastName     string
		WorkerNumber NullInt
		DateJoined   time.Time
		Salary       float64
		NullText     NullString
		NullVarchar  NullString
		NullBool     NullBool
		NullDate     NullTime
		NullInt      NullInt
		NullFloat    NullFloat64
	}

	u, err := url.Parse("http://localhost")
	if err != nil {
		t.Fatal(err)
	}

	pag, err := NewPaginator(Employees{}, "mysql", *u)
	if err != nil {
		t.Fatal(err)
	}

	cmd, args, err := pag.Paginate()
	if err != nil {
		t.Fatal(err)
	}

	rows, err := mysqlTestDB.Query(cmd, args...)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(pag.GetRowPtrArgs()...)
		if err != nil {
			t.Fatal(err)
		}
	}

	if err = rows.Err(); err != nil {
		t.Fatal(err)
	}

	results := make([]Employees, 0)

	for pag.NextData() {
		employee := Employees{}
		err = pag.Scan(&employee)
		if err != nil {
			t.Fatal(err)
		}
		results = append(results, employee)
	}

	if len(results) != 10 {
		t.Errorf("we should have 10 records in result; got %d", len(results))
	}
}

func TestNewPaginatorMysql_RequestParameter_Equal(t *testing.T) {
	type Employee struct {
		ID           int         `paginate:"filter;id;col=id"`
		Name         string      `paginate:"filter;col=name"`
		LastName     string      `paginate:"filter;col=last_name"`
		WorkerNumber NullInt     `paginate:"filter;col=worker_number"`
		DateJoined   time.Time   `paginate:"filter;col=date_joined"`
		Salary       float64     `paginate:"filter;col=salary"`
		NullText     NullString  `paginate:"filter;col=null_text"`
		NullVarchar  NullString  `paginate:"filter;col=null_varchar"`
		NullBool     NullBool    `paginate:"filter;col=null_bool"`
		NullDate     NullTime    `paginate:"filter;col=null_date"`
		NullInt      NullInt     `paginate:"filter;col=null_int"`
		NullFloat    NullFloat64 `paginate:"filter;col=null_float"`
	}

	u, err := url.Parse("http://localhost?name=Ringo")
	if err != nil {
		t.Fatal(err)
	}

	pag, err := NewPaginator(Employee{}, "mysql", *u, TableName("employees"))
	if err != nil {
		t.Fatal(err)
	}

	cmd, args, err := pag.Paginate()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(cmd)
	fmt.Println(args)

	rows, err := mysqlTestDB.Query(cmd, args...)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(pag.GetRowPtrArgs()...)
		if err != nil {
			t.Fatal(err)
		}
	}

	if err = rows.Err(); err != nil {
		t.Fatal(err)
	}

	results := make([]Employee, 0)

	for pag.NextData() {
		employee := Employee{}
		err = pag.Scan(&employee)
		if err != nil {
			t.Fatal(err)
		}
		results = append(results, employee)
	}

	if len(results) != 1 {
		t.Errorf("we should have 1 record in result; got %d", len(results))
	}

	if results[0].Name != "Ringo" {
		t.Errorf("expected Ringo in record inside results; got %s", results[0].Name)
	}

	if results[0].LastName != "Star" {
		t.Errorf("expected Star in record inside results; got %s", results[0].LastName)
	}

	if results[0].WorkerNumber.Int != 1 {
		t.Errorf("expected WorkNumber = 1 in record inside results; got %d", results[0].WorkerNumber.Int)
	}

	todayDateStr := time.Now().Format("2006/01/02")
	dateJoinedStr := results[0].DateJoined.Format("2006/01/02")

	if todayDateStr != dateJoinedStr {
		t.Errorf("expected DateJoined to be %s; got %s instead", todayDateStr, dateJoinedStr)
	}
}
