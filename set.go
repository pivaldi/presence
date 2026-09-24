package presence

import (
	"strconv"
	"strings"
)

// Placeholder renders the SQL placeholder for the n-th (1-based) bound argument.
// Use Dollar for PostgreSQL, Question for MySQL/SQLite, or provide your own.
type Placeholder func(n int) string

// Assignment pairs a column name with a presence value for SetClause.
// Build it with Set.
type Assignment struct {
	Column string
	isSet  bool
	value  any
}

// Dollar renders PostgreSQL-style placeholders: $1, $2, ...
func Dollar(n int) string {
	return "$" + strconv.Itoa(n)
}

// Question renders MySQL/SQLite-style placeholders: ?, ?, ...
func Question(_ int) string {
	return "?"
}

// Set creates an Assignment of column to the given presence value.
func Set[T any](column string, v Of[T]) Assignment {
	return Assignment{Column: column, isSet: v.IsSet(), value: v}
}

// SetClause builds the body of a SQL UPDATE ... SET clause from the given
// assignments, keeping only those whose value is set.
//
// An unset value (zero Of[T]) is skipped, so the column keeps its current
// value in the database. An explicit null is kept and writes NULL. A concrete
// value is kept and writes that value.
//
// The returned clause has the form "col1 = $1, col2 = $2" and args holds the
// matching presence values, which implement driver.Valuer. Placeholders are
// numbered from 1, so extra parameters such as a WHERE condition must be
// appended after args and numbered from len(args)+1.
//
// When no assignment is set, SetClause returns an empty clause and nil args:
// the caller should then skip the UPDATE entirely.
func SetClause(placeholder Placeholder, assignments ...Assignment) (clause string, args []any) {
	var sb strings.Builder

	for _, a := range assignments {
		if !a.isSet {
			continue
		}

		if len(args) > 0 {
			sb.WriteString(", ")
		}

		args = append(args, a.value)
		sb.WriteString(a.Column)
		sb.WriteString(" = ")
		sb.WriteString(placeholder(len(args)))
	}

	return sb.String(), args
}
