package paginate

import (
	"net/url"
	"testing"
)

func TestPaginatorMysql_HappyPath(t *testing.T) {
	type Employee struct {
		ID           int        `paginate:"id;col=id"`
		Name         string     `paginate:"col=name"`
		LastName     string     `paginate:"col=last_name"`
		WorkerNumber NullInt    `paginate:"col=worker_number"`
		DateJoined   NullTime   `paginate:"col=date_joined"`
		Salary       float64    `paginate:"col=salary"`
		NullText     NullString `paginate:"col=null_text"`
		NullVarchar  NullString `paginate:"col=null_varchar"`
		NullBool     NullBool   `paginate:"col=null_bool"`
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

	countNullBooleans := 0
	for _, r := range results {
		if !r.NullBool.Valid {
			countNullBooleans++
		}
	}

	if countNullBooleans != 8 {
		t.Errorf("we should have 8 nulls for the NullBool field")
	}
}
