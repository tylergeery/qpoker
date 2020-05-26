package sql

import (
	"fmt"
	"strings"
)

// Query is interface for SQL query
type Query interface {
	ToSQL() (string, []interface{})
}

// BaseQuery holds query aspects related to all queries
type BaseQuery struct {
	table      string
	expression Expression
	sort       sort
	limit      int
	offset     int

	hasExpression bool
	hasSort       bool
}

// Select contains entirety of single SQL query
type Select struct {
	BaseQuery
	columns []string
}

// Filter calls BaseQuery Filter
func (s Select) Filter(expression Expression) Select {
	s.BaseQuery.expression = expression
	s.BaseQuery.hasExpression = true

	return s
}

// Sort calls BaseQuery Sort
func (s Select) Sort(sorts ...string) Select {
	s.BaseQuery.sort = buildSort(sorts...)
	s.BaseQuery.hasSort = true

	return s
}

// Limit sets query limit
func (s Select) Limit(limit int) Select {
	s.BaseQuery.limit = limit

	return s
}

// Offset calls BaseQuery Offset
func (s Select) Offset(offset int) Select {
	s.BaseQuery.offset = offset

	return s
}

// NewSelect creates new Query struct
func NewSelect(table string, columns []string) Select {
	return Select{
		columns:   columns,
		BaseQuery: BaseQuery{table: table},
	}
}

// ToSQL returns a sql statement and prepared values for query
func (s Select) ToSQL() (string, []interface{}) {
	values := []interface{}{}

	query := fmt.Sprintf(
		`SELECT %s FROM %s`,
		strings.Join(s.columns, ", "),
		s.BaseQuery.table,
	)

	if s.BaseQuery.hasExpression {
		expression, prepared, _ := s.BaseQuery.expression.ToQuery(1)
		query += " WHERE " + expression
		values = prepared
	}

	if s.BaseQuery.hasSort {
		query += fmt.Sprintf(" ORDER BY %s", s.BaseQuery.sort.ToString())
	}

	if s.BaseQuery.limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", s.BaseQuery.limit)
	}

	if s.BaseQuery.offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", s.BaseQuery.offset)
	}

	fmt.Println("query:", query, "values:", values)
	return query, values
}
