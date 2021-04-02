package paginate

import "strings"

func NewInnerJoinClause(dialect string) (InnerJoin, error) {
	if err := dialectPlaceholder.CheckIfDialectIsSupported(dialect); err != nil {
		return InnerJoin{}, err
	}
	return InnerJoin{dialect: dialect}, nil
}

type InnerJoin struct {
	column       string
	targetTable  string
	targetColumn string
	dialect      string
}

func (clause *InnerJoin) On(column, targetTable, targetColumn string) {
	clause.column = column
	clause.targetTable = targetTable
	clause.targetColumn = targetColumn
}

func (clause *InnerJoin) clean() {
	clause.column = strings.TrimSpace(clause.column)
	clause.targetTable = strings.TrimSpace(clause.targetTable)
	clause.targetColumn = strings.TrimSpace(clause.targetColumn)
}
