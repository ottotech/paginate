package paginate

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"testing"
	"time"
)

func TestPaginatorPsql_HappyPath(t *testing.T) {
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

	pag, err := NewPaginator(Employee{}, "postgres", *u, TableName("employees"))
	if err != nil {
		t.Fatal(err)
	}

	cmd, args, err := pag.Paginate()
	if err != nil {
		t.Fatal(err)
	}

	rows, err := psqlTestDB.Query(cmd, args...)
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

func TestPaginatorPsql_JsonMarshalling(t *testing.T) {
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

	pag, err := NewPaginator(Employee{}, "postgres", *u, TableName("employees"))
	if err != nil {
		t.Fatal(err)
	}

	cmd, args, err := pag.Paginate()
	if err != nil {
		t.Fatal(err)
	}

	rows, err := psqlTestDB.Query(cmd, args...)
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
func TestNewPaginatorPsql_DefaultTableAndColumnsInferring(t *testing.T) {
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

	pag, err := NewPaginator(Employees{}, "postgres", *u)
	if err != nil {
		t.Fatal(err)
	}

	cmd, args, err := pag.Paginate()
	if err != nil {
		t.Fatal(err)
	}

	rows, err := psqlTestDB.Query(cmd, args...)
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

func TestNewPaginatorPsql_RequestParameter_Equal(t *testing.T) {
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

	pag, err := NewPaginator(Employee{}, "postgres", *u, TableName("employees"))
	if err != nil {
		t.Fatal(err)
	}

	cmd, args, err := pag.Paginate()
	if err != nil {
		t.Fatal(err)
	}

	rows, err := psqlTestDB.Query(cmd, args...)
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

	todayDateStr := time.Now().Local().Format("2006/01/02")
	dateJoinedStr := results[0].DateJoined.Local().Format("2006/01/02")

	if todayDateStr != dateJoinedStr {
		t.Errorf("expected DateJoined to be %s; got %s instead", todayDateStr, dateJoinedStr)
	}
}

func TestNewPaginatorPsql_RequestParameter_LessThan(t *testing.T) {
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

	u, err := url.Parse("http://localhost?salary<5000")
	if err != nil {
		t.Fatal(err)
	}

	pag, err := NewPaginator(Employee{}, "postgres", *u, TableName("employees"))
	if err != nil {
		t.Fatal(err)
	}

	cmd, args, err := pag.Paginate()
	if err != nil {
		t.Fatal(err)
	}

	rows, err := psqlTestDB.Query(cmd, args...)
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

	if len(results) != 2 {
		t.Errorf("we should have 2 records in result; got %d", len(results))
	}

	expectedResults := []struct {
		Name, LastName string
		WorkNumber     int
		Salary         float64
	}{
		{
			Name:       "John",
			LastName:   "Smith",
			WorkNumber: 4,
			Salary:     4650.90,
		}, {
			Name:       "Erika",
			LastName:   "Smith",
			WorkNumber: 8,
			Salary:     4455,
		},
	}

	for _, e := range expectedResults {
		isThere := false
		for _, r := range results {
			if r.Name == e.Name && r.LastName == e.LastName && r.WorkerNumber.Int == e.WorkNumber && r.Salary == e.Salary {
				isThere = true
				break
			}
		}
		if !isThere {
			t.Errorf("expected (%+v) in results", e)
		}
	}
}

func TestNewPaginatorPsql_RequestParameter_GreaterThan(t *testing.T) {
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

	u, err := url.Parse("http://localhost?salary>10000")
	if err != nil {
		t.Fatal(err)
	}

	pag, err := NewPaginator(Employee{}, "postgres", *u, TableName("employees"))
	if err != nil {
		t.Fatal(err)
	}

	cmd, args, err := pag.Paginate()
	if err != nil {
		t.Fatal(err)
	}

	rows, err := psqlTestDB.Query(cmd, args...)
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

	expectedResults := []struct {
		Name, LastName string
		WorkNumber     int
		Salary         float64
	}{
		{
			Name:       "Bill",
			LastName:   "Gates",
			WorkNumber: 2,
			Salary:     1200000,
		},
	}

	for _, e := range expectedResults {
		isThere := false
		for _, r := range results {
			if r.Name == e.Name && r.LastName == e.LastName && r.WorkerNumber.Int == e.WorkNumber && r.Salary == e.Salary {
				isThere = true
				break
			}
		}
		if !isThere {
			t.Errorf("expected (%+v) in results", e)
		}
	}
}

func TestNewPaginatorPsql_RequestParameter_GreaterEqualThan(t *testing.T) {
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

	u, err := url.Parse("http://localhost?salary>=1200000")
	if err != nil {
		t.Fatal(err)
	}

	pag, err := NewPaginator(Employee{}, "postgres", *u, TableName("employees"))
	if err != nil {
		t.Fatal(err)
	}

	cmd, args, err := pag.Paginate()
	if err != nil {
		t.Fatal(err)
	}

	rows, err := psqlTestDB.Query(cmd, args...)
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

	expectedResults := []struct {
		Name, LastName string
		WorkNumber     int
		Salary         float64
	}{
		{
			Name:       "Bill",
			LastName:   "Gates",
			WorkNumber: 2,
			Salary:     1200000,
		},
	}

	for _, e := range expectedResults {
		isThere := false
		for _, r := range results {
			if r.Name == e.Name && r.LastName == e.LastName && r.WorkerNumber.Int == e.WorkNumber && r.Salary == e.Salary {
				isThere = true
				break
			}
		}
		if !isThere {
			t.Errorf("expected (%+v) in results", e)
		}
	}
}

func TestNewPaginatorPsql_RequestParameter_LessEqualThan(t *testing.T) {
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

	u, err := url.Parse("http://localhost?salary<=4455")
	if err != nil {
		t.Fatal(err)
	}

	pag, err := NewPaginator(Employee{}, "postgres", *u, TableName("employees"))
	if err != nil {
		t.Fatal(err)
	}

	cmd, args, err := pag.Paginate()
	if err != nil {
		t.Fatal(err)
	}

	rows, err := psqlTestDB.Query(cmd, args...)
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

	expectedResults := []struct {
		Name, LastName string
		WorkNumber     int
		Salary         float64
	}{
		{
			Name:       "Erika",
			LastName:   "Smith",
			WorkNumber: 8,
			Salary:     4455,
		},
	}

	for _, e := range expectedResults {
		isThere := false
		for _, r := range results {
			if r.Name == e.Name && r.LastName == e.LastName && r.WorkerNumber.Int == e.WorkNumber && r.Salary == e.Salary {
				isThere = true
				break
			}
		}
		if !isThere {
			t.Errorf("expected (%+v) in results", e)
		}
	}
}

func TestNewPaginatorPsql_RequestParameter_NotEqual(t *testing.T) {
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

	u, err := url.Parse("http://localhost?salary<>4455")
	if err != nil {
		t.Fatal(err)
	}

	pag, err := NewPaginator(Employee{}, "postgres", *u, TableName("employees"))
	if err != nil {
		t.Fatal(err)
	}

	cmd, args, err := pag.Paginate()
	if err != nil {
		t.Fatal(err)
	}

	rows, err := psqlTestDB.Query(cmd, args...)
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

	if len(results) != 9 {
		t.Errorf("we should have 9 records in result; got %d", len(results))
	}

	expectedResults := []struct {
		Name, LastName string
		WorkNumber     int
		Salary         float64
	}{
		{
			Name:       "Ringo",
			LastName:   "Star",
			WorkNumber: 1,
			Salary:     5400,
		},
		{
			Name:       "Bill",
			LastName:   "Gates",
			WorkNumber: 2,
			Salary:     1200000,
		}, {
			Name:       "Mark",
			LastName:   "Smith",
			WorkNumber: 3,
			Salary:     8000,
		}, {
			Name:       "John",
			LastName:   "Smith",
			WorkNumber: 4,
			Salary:     4650.90,
		}, {
			Name:       "Fred",
			LastName:   "Smith",
			WorkNumber: 5,
			Salary:     7550,
		}, {
			Name:       "Rob",
			LastName:   "Williams",
			WorkNumber: 6,
			Salary:     9880,
		}, {
			Name:       "Juliana",
			LastName:   "Collier",
			WorkNumber: 7,
			Salary:     7788,
		}, {
			Name:       "Maria",
			LastName:   "Gomez",
			WorkNumber: 9,
			Salary:     7550,
		}, {
			Name:       "Rafael",
			LastName:   "Smith",
			WorkNumber: 10,
			Salary:     7550,
		},
	}

	for _, e := range expectedResults {
		isThere := false
		for _, r := range results {
			if r.Name == e.Name && r.LastName == e.LastName && r.WorkerNumber.Int == e.WorkNumber && r.Salary == e.Salary {
				isThere = true
				break
			}
		}
		if !isThere {
			t.Errorf("expected (%+v) in results", e)
		}
	}
}

func TestNewPaginatorPsql_RequestParameter_MultipleEqual_IN_Clause(t *testing.T) {
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

	u, err := url.Parse("http://localhost?salary=4455&salary=1200000&salary=8000")
	if err != nil {
		t.Fatal(err)
	}

	pag, err := NewPaginator(Employee{}, "postgres", *u, TableName("employees"))
	if err != nil {
		t.Fatal(err)
	}

	cmd, args, err := pag.Paginate()
	if err != nil {
		t.Fatal(err)
	}

	rows, err := psqlTestDB.Query(cmd, args...)
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

	if len(results) != 3 {
		t.Errorf("we should have 3 records in result; got %d", len(results))
	}

	expectedResults := []struct {
		Name, LastName string
		WorkNumber     int
		Salary         float64
	}{
		{
			Name:       "Mark",
			LastName:   "Smith",
			WorkNumber: 3,
			Salary:     8000,
		}, {
			Name:       "Bill",
			LastName:   "Gates",
			WorkNumber: 2,
			Salary:     1200000,
		}, {
			Name:       "Erika",
			LastName:   "Smith",
			WorkNumber: 8,
			Salary:     4455,
		},
	}

	for _, e := range expectedResults {
		isThere := false
		for _, r := range results {
			if r.Name == e.Name && r.LastName == e.LastName && r.WorkerNumber.Int == e.WorkNumber && r.Salary == e.Salary {
				isThere = true
				break
			}
		}
		if !isThere {
			t.Errorf("expected (%+v) in results", e)
		}
	}
}

func TestNewPaginatorPsql_RequestParameter_MultipleNotEqual_NOT_IN_Clause(t *testing.T) {
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

	u, err := url.Parse(
		"http://localhost?name<>Ringo&name<>Bill&name<>Mark&name<>John&name<>Fred&name<>Rob&name<>Juliana")
	if err != nil {
		t.Fatal(err)
	}

	pag, err := NewPaginator(Employee{}, "postgres", *u, TableName("employees"))
	if err != nil {
		t.Fatal(err)
	}

	cmd, args, err := pag.Paginate()
	if err != nil {
		t.Fatal(err)
	}

	rows, err := psqlTestDB.Query(cmd, args...)
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

	if len(results) != 3 {
		t.Errorf("we should have 3 records in result; got %d", len(results))
	}

	expectedResults := []struct {
		Name, LastName string
		WorkNumber     int
		Salary         float64
	}{
		{
			Name:       "Erika",
			LastName:   "Smith",
			WorkNumber: 8,
			Salary:     4455,
		},
		{
			Name:       "Maria",
			LastName:   "Gomez",
			WorkNumber: 9,
			Salary:     7550,
		},
		{
			Name:       "Rafael",
			LastName:   "Smith",
			WorkNumber: 10,
			Salary:     7550,
		},
	}

	for _, e := range expectedResults {
		isThere := false
		for _, r := range results {
			if r.Name == e.Name && r.LastName == e.LastName && r.WorkerNumber.Int == e.WorkNumber && r.Salary == e.Salary {
				isThere = true
				break
			}
		}
		if !isThere {
			t.Errorf("expected (%+v) in results", e)
		}
	}
}

func TestNewPaginatorPsql_RequestParameter_Raw_Sql_LIKE_Clause(t *testing.T) {
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

	rawSql, err := NewRawWhereClause("postgres")
	if err != nil {
		t.Fatal(err)
	}
	rawSql.AddPredicate("name ILIKE ? OR last_name ILIKE ?")
	rawSql.AddArg("ringo")
	rawSql.AddArg("smith")

	pag, err := NewPaginator(Employee{}, "postgres", *u, TableName("employees"))
	if err != nil {
		t.Fatal(err)
	}

	err = pag.AddWhereClause(rawSql)
	if err != nil {
		t.Fatal(err)
	}

	cmd, args, err := pag.Paginate()
	if err != nil {
		t.Fatal(err)
	}

	rows, err := psqlTestDB.Query(cmd, args...)
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

	if len(results) != 6 {
		t.Errorf("we should have 6 records in result; got %d", len(results))
	}

	expectedResults := []struct {
		Name, LastName string
		WorkNumber     int
		Salary         float64
	}{
		{
			Name:       "Ringo",
			LastName:   "Star",
			WorkNumber: 1,
			Salary:     5400,
		}, {
			Name:       "Mark",
			LastName:   "Smith",
			WorkNumber: 3,
			Salary:     8000,
		}, {
			Name:       "John",
			LastName:   "Smith",
			WorkNumber: 4,
			Salary:     4650.90,
		}, {
			Name:       "Fred",
			LastName:   "Smith",
			WorkNumber: 5,
			Salary:     7550,
		}, {
			Name:       "Erika",
			LastName:   "Smith",
			WorkNumber: 8,
			Salary:     4455,
		}, {
			Name:       "Rafael",
			LastName:   "Smith",
			WorkNumber: 10,
			Salary:     7550,
		},
	}

	for _, e := range expectedResults {
		isThere := false
		for _, r := range results {
			if r.Name == e.Name && r.LastName == e.LastName && r.WorkerNumber.Int == e.WorkNumber && r.Salary == e.Salary {
				isThere = true
				break
			}
		}
		if !isThere {
			t.Errorf("expected (%+v) in results", e)
		}
	}
}

func TestNewPaginatorPsql_RequestParameter_Raw_Sql_IS_NULL_Clause(t *testing.T) {
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

	rawSql, err := NewRawWhereClause("postgres")
	if err != nil {
		t.Fatal(err)
	}
	rawSql.AddPredicate("null_float IS NULL")

	pag, err := NewPaginator(Employee{}, "postgres", *u, TableName("employees"))
	if err != nil {
		t.Fatal(err)
	}

	err = pag.AddWhereClause(rawSql)
	if err != nil {
		t.Fatal(err)
	}

	cmd, args, err := pag.Paginate()
	if err != nil {
		t.Fatal(err)
	}

	rows, err := psqlTestDB.Query(cmd, args...)
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
}

func TestNewPaginatorPsql_RequestParameter_Sort_ASC(t *testing.T) {
	type Employee struct {
		ID           int         `paginate:"id;col=id"`
		Name         string      `paginate:"filter;col=name"`
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

	u, err := url.Parse("http://localhost?sort=+name")
	if err != nil {
		t.Fatal(err)
	}

	pag, err := NewPaginator(Employee{}, "postgres", *u, TableName("employees"), PageSize(1))
	if err != nil {
		t.Fatal(err)
	}

	cmd, args, err := pag.Paginate()
	if err != nil {
		t.Fatal(err)
	}

	rows, err := psqlTestDB.Query(cmd, args...)
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

	expectedResults := []struct {
		Name, LastName string
		WorkNumber     int
		Salary         float64
	}{
		{
			Name:       "Bill",
			LastName:   "Gates",
			WorkNumber: 2,
			Salary:     1200000,
		},
	}

	for _, e := range expectedResults {
		isThere := false
		for _, r := range results {
			if r.Name == e.Name && r.LastName == e.LastName && r.WorkerNumber.Int == e.WorkNumber && r.Salary == e.Salary {
				isThere = true
				break
			}
		}
		if !isThere {
			t.Errorf("expected (%+v) in results", e)
		}
	}
}

func TestNewPaginatorPsql_RequestParameter_Sort_DESC(t *testing.T) {
	type Employee struct {
		ID           int         `paginate:"id;col=id"`
		Name         string      `paginate:"filter;col=name"`
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

	u, err := url.Parse("http://localhost?sort=-name")
	if err != nil {
		t.Fatal(err)
	}

	pag, err := NewPaginator(Employee{}, "postgres", *u, TableName("employees"), PageSize(1))
	if err != nil {
		t.Fatal(err)
	}

	cmd, args, err := pag.Paginate()
	if err != nil {
		t.Fatal(err)
	}

	rows, err := psqlTestDB.Query(cmd, args...)
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

	expectedResults := []struct {
		Name, LastName string
		WorkNumber     int
		Salary         float64
	}{
		{
			Name:       "Rob",
			LastName:   "Williams",
			WorkNumber: 6,
			Salary:     9880,
		},
	}

	for _, e := range expectedResults {
		isThere := false
		for _, r := range results {
			if r.Name == e.Name && r.LastName == e.LastName && r.WorkerNumber.Int == e.WorkNumber && r.Salary == e.Salary {
				isThere = true
				break
			}
		}
		if !isThere {
			t.Errorf("expected (%+v) in results (%+v)", e, results)
		}
	}
}

func TestNewPaginatorPsql_With_Custom_OrderByAsc_Clauses(t *testing.T) {
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

	pag, err := NewPaginator(Employee{}, "postgres", *u, TableName("employees"), PageSize(1), OrderByAsc("name"))
	if err != nil {
		t.Fatal(err)
	}

	cmd, args, err := pag.Paginate()
	if err != nil {
		t.Fatal(err)
	}

	rows, err := psqlTestDB.Query(cmd, args...)
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

	expectedResults := []struct {
		Name, LastName string
		WorkNumber     int
		Salary         float64
	}{
		{
			Name:       "Bill",
			LastName:   "Gates",
			WorkNumber: 2,
			Salary:     1200000,
		},
	}

	for _, e := range expectedResults {
		isThere := false
		for _, r := range results {
			if r.Name == e.Name && r.LastName == e.LastName && r.WorkerNumber.Int == e.WorkNumber && r.Salary == e.Salary {
				isThere = true
				break
			}
		}
		if !isThere {
			t.Errorf("expected (%+v) in results (%+v)", e, results)
		}
	}
}

func TestNewPaginatorPsql_With_Custom_OrderByDesc_Clauses(t *testing.T) {
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

	pag, err := NewPaginator(Employee{}, "postgres", *u, TableName("employees"), PageSize(1), OrderByDesc("name"))
	if err != nil {
		t.Fatal(err)
	}

	cmd, args, err := pag.Paginate()
	if err != nil {
		t.Fatal(err)
	}

	rows, err := psqlTestDB.Query(cmd, args...)
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

	expectedResults := []struct {
		Name, LastName string
		WorkNumber     int
		Salary         float64
	}{
		{
			Name:       "Rob",
			LastName:   "Williams",
			WorkNumber: 6,
			Salary:     9880,
		},
	}

	for _, e := range expectedResults {
		isThere := false
		for _, r := range results {
			if r.Name == e.Name && r.LastName == e.LastName && r.WorkerNumber.Int == e.WorkNumber && r.Salary == e.Salary {
				isThere = true
				break
			}
		}
		if !isThere {
			t.Errorf("expected (%+v) in results (%+v)", e, results)
		}
	}
}

func ExampleTest_NewPaginatorPsql_With_Custom_OrderByDesc_Clause_With_ID() {
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

	u, err := url.Parse("http://localhost?sort=-id")
	if err != nil {
		log.Fatalln(err)
	}

	pag, err := NewPaginator(Employee{}, "postgres", *u, TableName("employees"), OrderByDesc("id"))
	if err != nil {
		log.Fatalln(err)
	}

	cmd, args, err := pag.Paginate()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(cmd)
	fmt.Printf("args length: %v\n", len(args))
	// Output:
	// SELECT id, name, last_name, worker_number, date_joined, salary, null_text, null_varchar, null_bool, null_date, null_int, null_float, count(*) over() FROM employees ORDER BY id LIMIT 30 OFFSET 0
	// args length: 0
}

func Test_InnerJoin_Psql_Query_String_HappyPath(t *testing.T) {
	type Employee struct {
		ID       int    `paginate:"id;col=id"`
		Name     string `paginate:"col=name"`
		LastName string `paginate:"col=last_name"`
	}

	u, err := url.Parse("http://localhost?sort=-id")
	if err != nil {
		log.Fatalln(err)
	}

	pag, err := NewPaginator(Employee{}, "postgres", *u, TableName("employees"))
	if err != nil {
		log.Fatalln(err)
	}

	rawSql, err := NewRawWhereClause("postgres")
	if err != nil {
		t.Fatal(err)
	}

	rawSql.AddPredicate("name ILIKE ?")
	rawSql.AddArg("%ringo%")

	err = pag.AddWhereClause(rawSql)
	if err != nil {
		t.Fatal(err)
	}

	innerClause, err := NewInnerJoinClause("postgres")
	if err != nil {
		t.Fatal(err)
	}

	innerClause.On("id", "managers", "employee_id")

	err = pag.AddJoinClause(innerClause)
	if err != nil {
		t.Fatal(err)
	}

	sql, args, err := pag.Paginate()
	if err != nil {
		t.Fatal(err)
	}

	expectedSql := "SELECT id, name, last_name, count(*) over() FROM employees JOIN managers ON employees.id = managers.employee_id WHERE name ILIKE $1 ORDER BY id LIMIT 30 OFFSET 0"
	expectedArg := "%ringo%"

	if sql != expectedSql {
		t.Errorf("expected sql %q; got %q instead", expectedSql, sql)
	}

	if args[0] != expectedArg {
		t.Errorf("expected arg %q; got %q instead", expectedArg, args[0])
	}
}

func Test_InnerJoin_Psql_Employees_That_Are_Managers(t *testing.T) {
	type Employee struct {
		ID         int       `paginate:"id;col=id"`
		Name       string    `paginate:"col=name"`
		LastName   string    `paginate:"col=last_name"`
		WorkNumber int64     `paginate:"col=worker_number"`
		DateJoined time.Time `paginate:"col=date_joined"`
		Salary     int64     `paginate:"col=salary"`
	}

	u, err := url.Parse("http://localhost")
	if err != nil {
		log.Fatalln(err)
	}

	pag, err := NewPaginator(Employee{}, "postgres", *u, TableName("employees"))
	if err != nil {
		log.Fatalln(err)
	}

	innerClause, err := NewInnerJoinClause("postgres")
	if err != nil {
		t.Fatal(err)
	}

	innerClause.On("id", "manager", "employee_id")

	err = pag.AddJoinClause(innerClause)
	if err != nil {
		t.Fatal(err)
	}

	sql, args, err := pag.Paginate()
	if err != nil {
		t.Fatal(err)
	}

	rows, err := psqlTestDB.Query(sql, args...)
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

	managers := employees[5:]

	if len(results) != len(managers) {
		t.Errorf("expected to have %[1]d results since there are %[1]d managers.", len(managers))
	}

	for i := 0; i < len(results); i++ {
		expectedManager := managers[i]
		resultManager := results[i]

		if expectedManager.Name != resultManager.Name {
			t.Errorf("expected manager name to be %s; got %s", expectedManager.Name, resultManager.Name)
		}
		if expectedManager.LastName != resultManager.LastName {
			t.Errorf("expected manager last name to be %s; got %s", expectedManager.LastName, resultManager.LastName)
		}
		if int64(expectedManager.Salary) != resultManager.Salary {
			t.Errorf("expected manager salary to be %d; got %d", int64(expectedManager.Salary), resultManager.Salary)
		}
		if expectedManager.DateJoined.Format("02-01-2006") != resultManager.DateJoined.Format("02-01-2006") {
			t.Errorf("expected manager date joined to be %s; got %s", expectedManager.DateJoined.Format("02-01-2006"), resultManager.DateJoined.Format("02-01-2006"))
		}
		if expectedManager.WorkNumber != int(resultManager.WorkNumber) {
			t.Errorf("expected manager worker number to be %d; got %d", expectedManager.WorkNumber, resultManager.WorkNumber)
		}
	}
}

func Test_InnerJoin_Psql_Employees_That_Are_Go_Developers(t *testing.T) {
	// ProgrammingLanguage it is a column from the joined table ("developer").
	// In this test we are able to prove that it is possible to filter
	// joined columns. However, if there are clashes with the column names
	// between two or more joined tables, to get the correct name of the columns
	// it's a bit hard. I need to think more about how to smooth this problem.
	type Employee struct {
		ID                  int       `paginate:"id;col=id"`
		Name                string    `paginate:"col=name"`
		LastName            string    `paginate:"col=last_name"`
		WorkNumber          int64     `paginate:"col=worker_number"`
		DateJoined          time.Time `paginate:"col=date_joined"`
		Salary              float64   `paginate:"col=salary"`
		ProgrammingLanguage string    `paginate:"filter;col=programming_language;param=lg"`
	}

	u, err := url.Parse("http://localhost?lg=Go")
	if err != nil {
		log.Fatalln(err)
	}

	pag, err := NewPaginator(Employee{}, "postgres", *u, TableName("employees"))
	if err != nil {
		log.Fatalln(err)
	}

	innerClause, err := NewInnerJoinClause("postgres")
	if err != nil {
		t.Fatal(err)
	}

	innerClause.On("id", "developer", "employee_id")

	err = pag.AddJoinClause(innerClause)
	if err != nil {
		t.Fatal(err)
	}

	sql, args, err := pag.Paginate()
	if err != nil {
		t.Fatal(err)
	}

	rows, err := psqlTestDB.Query(sql, args...)
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

	goDevelopers := employees[:3]

	if len(results) != len(goDevelopers) {
		t.Errorf("expected to have %[1]d results since there are %[1]d Go developers.", len(goDevelopers))
	}

	for i := 0; i < len(results); i++ {
		expectedGoDeveloper := goDevelopers[i]
		resultGoDeveloper := results[i]

		if expectedGoDeveloper.Name != resultGoDeveloper.Name {
			t.Errorf("expected manager name to be %s; got %s", expectedGoDeveloper.Name, resultGoDeveloper.Name)
		}
		if expectedGoDeveloper.LastName != resultGoDeveloper.LastName {
			t.Errorf("expected manager last name to be %s; got %s", expectedGoDeveloper.LastName, resultGoDeveloper.LastName)
		}
		if expectedGoDeveloper.Salary != resultGoDeveloper.Salary {
			t.Errorf("expected manager salary to be %f; got %f", expectedGoDeveloper.Salary, resultGoDeveloper.Salary)
		}
		if expectedGoDeveloper.DateJoined.Format("02-01-2006") != resultGoDeveloper.DateJoined.Format("02-01-2006") {
			t.Errorf("expected manager date joined to be %s; got %s", expectedGoDeveloper.DateJoined.Format("02-01-2006"), resultGoDeveloper.DateJoined.Format("02-01-2006"))
		}
		if expectedGoDeveloper.WorkNumber != int(resultGoDeveloper.WorkNumber) {
			t.Errorf("expected manager worker number to be %d; got %d", expectedGoDeveloper.WorkNumber, resultGoDeveloper.WorkNumber)
		}
	}

	// Let's check that all developers have "Go" as their programming language.
	allAreGophers := true
	for _, r := range results {
		if r.ProgrammingLanguage != "Go" {
			allAreGophers = false
			break
		}
	}
	if !allAreGophers {
		t.Error("All developers should be gophers")
	}
}

func Test_InnerJoin_Psql_Employees_That_Are_Python_Developers_With_Pagination(t *testing.T) {
	type Employee struct {
		ID                  int       `paginate:"id;col=id"`
		Name                string    `paginate:"col=name"`
		LastName            string    `paginate:"col=last_name"`
		WorkNumber          int64     `paginate:"col=worker_number"`
		DateJoined          time.Time `paginate:"col=date_joined"`
		Salary              float64   `paginate:"col=salary"`
		ProgrammingLanguage string    `paginate:"filter;col=programming_language;param=lg"`
	}

	u, err := url.Parse("http://localhost?lg=Python&page=2")
	if err != nil {
		log.Fatalln(err)
	}

	opt := PageSize(1)

	pag, err := NewPaginator(Employee{}, "postgres", *u, TableName("employees"), opt)
	if err != nil {
		log.Fatalln(err)
	}

	innerClause, err := NewInnerJoinClause("postgres")
	if err != nil {
		t.Fatal(err)
	}

	innerClause.On("id", "developer", "employee_id")

	err = pag.AddJoinClause(innerClause)
	if err != nil {
		t.Fatal(err)
	}

	sql, args, err := pag.Paginate()
	if err != nil {
		t.Fatal(err)
	}

	rows, err := psqlTestDB.Query(sql, args...)
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
		t.Fatal("expected to have only one result, " +
			"since there are 2 python developers and we are paginating " +
			"the results and we are in the second page.")
	}

	expectedPythonDeveloper := employees[4]
	resultPythonDeveloper := results[0]

	if expectedPythonDeveloper.Name != resultPythonDeveloper.Name {
		t.Fatalf("expected to have %s python developer; got %s", expectedPythonDeveloper.Name, resultPythonDeveloper.Name)
	}
}
