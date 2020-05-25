package sql

import (
	"fmt"
	"strings"
)

type singleSort struct {
	column    string
	direction string
}

type sort struct {
	sorts []singleSort
}

// ToString gets sort to query string
func (s sort) ToString() string {
	sortParts := []string{}

	for _, sort := range s.sorts {
		sortParts = append(sortParts, fmt.Sprintf("%s %s", sort.column, sort.direction))
	}

	return strings.Join(sortParts, ", ")
}

// Sort creates a sort for query
func Sort(sorts []string) sort {
	sortParts := []singleSort{}

	for _, col := range sorts {
		dir := "ASC"
		if col[0] == '-' {
			dir, col = "DESC", col[1:]
		}

		sortParts = append(sortParts, singleSort{col, dir})
	}

	return sort{sortParts}
}
