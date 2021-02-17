package paginate

import (
	"fmt"
	"strings"
)

// NewRawWhereClause will give you a validated instance of a RawWhereClause object.
//
// Use this constructor whenever you want to create a custom sql
// "where" clause that Paginator can use to paginate your table
// accordingly.
//
// Example of how to use the returned RawWhereClause:
//
//  ...
//  paginator, _ := NewPaginator(MyTable{}, "mysql", *url)
//  rawWhereSql, err := NewRawWhereClause("mysql")
//  if err != nil {
//  	t.Fatal(err)
//  }
//
//  rawWhereSql.AddPredicate("name LIKE ? OR last_name LIKE ?")
//  rawWhereSql.AddArg("%ringo%")
//  rawWhereSql.AddArg("%smith%")
//
//  err = paginator.AddWhereClause(rawWhereSql)
//  if err != nil {
//     // Handle error gracefully.
//  }
//
func NewRawWhereClause(dialect string) (RawWhereClause, error) {
	if err := dialectPlaceholder.CheckIfDialectIsSupported(dialect); err != nil {
		return RawWhereClause{}, err
	}

	return RawWhereClause{
		dialect: dialect,
	}, nil
}

type RawWhereClause struct {
	predicate string
	args      []interface{}
	dialect   string
}

func (raw RawWhereClause) String() string {
	if raw.dialect == "postgres" {
		pred := strings.Replace(raw.predicate, "?", "$%v", -1)
		return fmt.Sprint(pred)
	}
	return raw.predicate
}

func (raw *RawWhereClause) AddArg(v interface{}) {
	raw.args = append(raw.args, v)
}

func (raw *RawWhereClause) AddPredicate(predicate string) {
	raw.predicate = predicate
}
