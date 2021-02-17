package paginate

import (
	"fmt"
	"strings"
)

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
