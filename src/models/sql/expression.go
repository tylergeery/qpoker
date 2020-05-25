package sql

import (
	"fmt"
	"strings"
)

// Expression holds interface for quering by values
type Expression interface {
	ToQuery(i int) (string, []interface{}, int)
}

// SimpleExpression holds parts of a simple query expression
type SimpleExpression struct {
	column     string
	comparison string
	value      interface{}
}

// ToQuery returns sql query, prepared values, and new value offset
func (e SimpleExpression) ToQuery(i int) (string, []interface{}, int) {
	query := fmt.Sprintf("%s %s $%d", e.column, e.comparison, i)

	return query, []interface{}{e.value}, i + 1
}

// CombinedExpression holds parts of combined query expression
type CombinedExpression struct {
	children    []Expression
	conjunction string
}

// ToQuery returns sql query, prepared values, and new value offset
func (e CombinedExpression) ToQuery(i int) (string, []interface{}, int) {
	queryParts := []string{}
	values := []interface{}{}

	for _, expression := range e.children {
		query, prepared, next := expression.ToQuery(i)
		values, i = append(values, prepared...), next
		queryParts = append(queryParts, query)
	}

	return strings.Join(queryParts, fmt.Sprintf(" %s ", e.conjunction)), values, i
}

// And returns an expression of a combined "and" Expressions
func And(expressions ...Expression) CombinedExpression {
	return CombinedExpression{expressions, "AND"}
}

// Or returns an expression of a combined "or" Expressions
func Or(expressions ...Expression) CombinedExpression {
	return CombinedExpression{expressions, "OR"}
}

// Equal creates a new expression with "=" operator
func Equal(col string, val interface{}) SimpleExpression {
	return SimpleExpression{col, "=", val}
}

// GT creates a new expression with "=" operator
func GT(col string, val interface{}) SimpleExpression {
	return SimpleExpression{col, ">", val}
}

// GTE creates a new expression with ">=" operator
func GTE(col string, val interface{}) SimpleExpression {
	return SimpleExpression{col, ">=", val}
}

// LT creates a new expression with "<" operator
func LT(col string, val interface{}) SimpleExpression {
	return SimpleExpression{col, "<", val}
}

// LTE creates a new expression with "<=" operator
func LTE(col string, val interface{}) Expression {
	return SimpleExpression{col, "<=", val}
}
