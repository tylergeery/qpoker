package sql

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSelectQuery(t *testing.T) {
	t1, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	type TestCase struct {
		query          Query
		expectedSQL    string
		expectedValues []interface{}
	}
	testCases := []TestCase{
		TestCase{
			query:          NewSelect("test", []string{"id"}),
			expectedSQL:    "SELECT id FROM test",
			expectedValues: []interface{}{},
		},
		TestCase{
			query: NewSelect(
				"test", []string{"id", "username", "email"},
			).Limit(10).Offset(100),
			expectedSQL:    "SELECT id, username, email FROM test LIMIT 10 OFFSET 100",
			expectedValues: []interface{}{},
		},
		TestCase{
			query: NewSelect(
				"values_inc", []string{"foo", "bar"},
			).Sort("-foo", "id").Limit(5),
			expectedSQL:    "SELECT foo, bar FROM values_inc ORDER BY foo DESC, id ASC LIMIT 5",
			expectedValues: []interface{}{},
		},
		TestCase{
			query: NewSelect(
				"tmp_table_casd", []string{"foo", "bar", "1", "2", "hello"},
			).Filter(
				Equal("as", 1),
			).Sort("updated_at").Limit(5),
			expectedSQL:    "SELECT foo, bar, 1, 2, hello FROM tmp_table_casd WHERE as = $1 ORDER BY updated_at ASC LIMIT 5",
			expectedValues: []interface{}{1},
		},
		TestCase{
			query: NewSelect(
				"test", []string{"hello"},
			).Filter(
				Or(
					And(
						Equal("as", 1),
						Equal("long_column", "other_value"),
						GTE("created_at", t1),
					),
				),
			),
			expectedSQL:    "SELECT hello FROM test WHERE as = $1 AND long_column = $2 AND created_at >= $3",
			expectedValues: []interface{}{1, "other_value", t1},
		},
	}

	for _, c := range testCases {
		sql, prepared := c.query.ToSQL()
		assert.Equal(t, c.expectedSQL, sql)
		assert.Equal(t, c.expectedValues, prepared)
	}
}
