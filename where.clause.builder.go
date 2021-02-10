package paginate

import (
	"fmt"
	"strings"
)

type RawWhereClause struct {
	Predicate string
	Args      []interface{}
}

func (raw RawWhereClause) String() string {
	pred := strings.Replace(raw.Predicate, "?", "$%v", -1)
	return fmt.Sprint(pred)
}

func (raw RawWhereClause) isEmpty() bool {
	if raw.Predicate == "" || len(raw.Args) == 0 {
		return true
	}
	return false
}

func (raw *RawWhereClause) AddArg(v interface{}) {
	raw.Args = append(raw.Args, v)
}

func (raw *RawWhereClause) AddPredicate(predicate string) {
	raw.Predicate = predicate
}
