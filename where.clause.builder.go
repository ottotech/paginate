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
// Example of how to use the returned RawWhereClause object:
//
//  ...
//  paginator, _ := NewPaginator(MyTable{}, "mysql", *url)
//  rawWhereSql, err := NewRawWhereClause("mysql")
//  if err != nil {
//     // Handle error gracefully.
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

// String returns the RawWhereClause predicate without arguments as string.
func (raw RawWhereClause) String() string {
	if raw.dialect == "postgres" {
		pred := strings.Replace(raw.predicate, "?", "$%v", -1)
		return fmt.Sprint(pred)
	}
	return raw.predicate
}

// AddPredicate adds the given predicate to a RawWhereClause instance.
// If your sql "where" clause requires multiple arguments use the
// question mark symbol "?" as placeholders, later use AddArg to
// add the arguments you need.
func (raw *RawWhereClause) AddPredicate(predicate string) {
	raw.predicate = predicate
}

// AddArg adds the given argument to the RawWhereClause instance.
// If your custom raw sql where clause requires more than one argument
// you should call AddArg multiple times.
func (raw *RawWhereClause) AddArg(v interface{}) {
	raw.args = append(raw.args, v)
}


