package paginate

import (
	"fmt"
	"strings"
)

type RawWhereClause struct {
	predicate string
	args      []interface{}
}

func (raw RawWhereClause) String() string {
	pred := strings.Replace(raw.predicate, "?", "$%v", -1)
	return fmt.Sprint(pred)
}

func (raw RawWhereClause) isEmpty() bool {
	if raw.predicate == "" || len(raw.args) == 0 {
		return true
	}
	return false
}

func (raw *RawWhereClause) AddArg(v interface{}) {
	raw.args = append(raw.args, v)
}

func (raw *RawWhereClause) AddPredicate(predicate string) {
	raw.predicate = predicate
}
