package tests

import (
	"testing"
	"time"

	"github.com/pivaldi/presence/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetClause(t *testing.T) {
	t.Run("all fields set", func(t *testing.T) {
		clause, args := presence.SetClause(presence.Dollar,
			presence.Set("name", presence.FromValue("Alice")),
			presence.Set("age", presence.FromValue(30)),
		)

		assert.Equal(t, "name = $1, age = $2", clause)
		require.Len(t, args, 2)
		assert.Equal(t, presence.FromValue("Alice"), args[0])
		assert.Equal(t, presence.FromValue(30), args[1])
	})

	t.Run("unset fields are skipped and numbering stays contiguous", func(t *testing.T) {
		var age presence.Of[int]

		clause, args := presence.SetClause(presence.Dollar,
			presence.Set("name", presence.FromValue("Alice")),
			presence.Set("age", age),
			presence.Set("email", presence.FromValue("a@b.c")),
		)

		assert.Equal(t, "name = $1, email = $2", clause)
		require.Len(t, args, 2)
		assert.Equal(t, presence.FromValue("Alice"), args[0])
		assert.Equal(t, presence.FromValue("a@b.c"), args[1])
	})

	t.Run("explicit null is kept", func(t *testing.T) {
		clause, args := presence.SetClause(presence.Dollar,
			presence.Set("email", presence.Null[string]()),
		)

		assert.Equal(t, "email = $1", clause)
		require.Len(t, args, 1)
		assert.Equal(t, presence.Null[string](), args[0])
	})

	t.Run("nothing set returns empty clause and nil args", func(t *testing.T) {
		var name presence.Of[string]
		var age presence.Of[int]

		clause, args := presence.SetClause(presence.Dollar,
			presence.Set("name", name),
			presence.Set("age", age),
		)

		assert.Empty(t, clause)
		assert.Nil(t, args)
	})

	t.Run("no assignments returns empty clause and nil args", func(t *testing.T) {
		clause, args := presence.SetClause(presence.Dollar)

		assert.Empty(t, clause)
		assert.Nil(t, args)
	})

	t.Run("question placeholder", func(t *testing.T) {
		clause, args := presence.SetClause(presence.Question,
			presence.Set("name", presence.FromValue("Alice")),
			presence.Set("age", presence.FromValue(30)),
		)

		assert.Equal(t, "name = ?, age = ?", clause)
		assert.Len(t, args, 2)
	})

	t.Run("custom placeholder", func(t *testing.T) {
		named := func(n int) string { return "@p" + string(rune('0'+n)) }

		clause, _ := presence.SetClause(named,
			presence.Set("name", presence.FromValue("Alice")),
			presence.Set("age", presence.FromValue(30)),
		)

		assert.Equal(t, "name = @p1, age = @p2", clause)
	})

	t.Run("args work as driver values", func(t *testing.T) {
		_, args := presence.SetClause(presence.Dollar,
			presence.Set("name", presence.FromValue("Alice")),
			presence.Set("email", presence.Null[string]()),
		)

		require.Len(t, args, 2)

		nameValuer, ok := args[0].(presence.Of[string])
		require.True(t, ok)
		nameValue, err := nameValuer.Value()
		require.NoError(t, err)
		assert.Equal(t, "Alice", nameValue)

		emailValuer, ok := args[1].(presence.Of[string])
		require.True(t, ok)
		emailValue, err := emailValuer.Value()
		require.NoError(t, err)
		assert.Nil(t, emailValue)
	})
}

func TestSetClausePostgres(t *testing.T) {
	db := getDB(t)
	cleanupTables(t, db, "test")

	initial := now.Truncate(time.Second)

	var id int64
	err := db.QueryRow(
		`INSERT INTO test (name, date_to, data) VALUES ($1, $2, $3) RETURNING id`,
		presence.FromValue("original"),
		presence.FromValue(initial),
		presence.FromValue(map[string]any{"k": "v"}),
	).Scan(&id)
	require.NoError(t, err)

	// Patch: change name, clear data, leave date_to untouched.
	var dateTo presence.Of[time.Time]
	clause, args := presence.SetClause(presence.Dollar,
		presence.Set("name", presence.FromValue("updated")),
		presence.Set("date_to", dateTo),
		presence.Set("data", presence.Null[map[string]any]()),
	)
	require.NotEmpty(t, clause)

	args = append(args, id)
	_, err = db.Exec(`UPDATE test SET `+clause+` WHERE id = `+presence.Dollar(len(args)), args...)
	require.NoError(t, err)

	var got testedStruct[map[string]any]
	err = db.QueryRow(`SELECT id, name, date_to, data FROM test WHERE id = $1`, id).
		Scan(&got.ID, &got.Name, &got.DateTo, &got.Data)
	require.NoError(t, err)

	t.Run("set field is updated", func(t *testing.T) {
		assert.Equal(t, "updated", got.Name.MustGet())
	})

	t.Run("unset field keeps its current value", func(t *testing.T) {
		assert.True(t, got.DateTo.IsValue())
		assert.Equal(t, initial, got.DateTo.MustGet().Truncate(time.Second))
	})

	t.Run("explicit null clears the column", func(t *testing.T) {
		assert.True(t, got.Data.IsNull())
	})
}
