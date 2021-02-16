package paginate

import (
	"encoding/json"
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

	pag, err := NewPaginator(Employee{}, *u, TableName("employees"))
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

	pag, err := NewPaginator(Employee{}, *u, TableName("employees"))
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

	_, err = json.Marshal(results)
	if err != nil {
		t.Errorf("we should be able to marshal results into json; got err %s", err)
	}
}
