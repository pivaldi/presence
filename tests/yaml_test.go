package tests

import (
	"testing"
	"time"

	"github.com/pivaldi/presence/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type yamlRow struct {
	Name  presence.Of[string]    `yaml:"name"`
	Age   presence.Of[int64]     `yaml:"age"`
	Score presence.Of[float64]   `yaml:"score"`
	OK    presence.Of[bool]      `yaml:"ok"`
	At    presence.Of[time.Time] `yaml:"at"`
	Tags  presence.Of[[]string]  `yaml:"tags"`
	Note  presence.Of[string]    `yaml:"note,omitempty"`
}

func TestMarshalYAML(t *testing.T) {
	t.Run("value encodes as the plain value", func(t *testing.T) {
		b, err := yaml.Marshal(presence.FromValue(int64(42)))
		require.NoError(t, err)
		assert.Equal(t, "42\n", string(b))
	})

	t.Run("null encodes as null", func(t *testing.T) {
		b, err := yaml.Marshal(presence.Null[string]())
		require.NoError(t, err)
		assert.Equal(t, "null\n", string(b))
	})

	t.Run("unset encodes as null", func(t *testing.T) {
		b, err := yaml.Marshal(presence.Of[string]{})
		require.NoError(t, err)
		assert.Equal(t, "null\n", string(b))
	})

	t.Run("unset field is omitted with omitempty", func(t *testing.T) {
		b, err := yaml.Marshal(yamlRow{Name: presence.FromValue("Dupont")})
		require.NoError(t, err)
		assert.NotContains(t, string(b), "note")
		assert.Contains(t, string(b), "name: Dupont")
	})

	t.Run("unset field is kept as null with omitempty under UnsetNull", func(t *testing.T) {
		row := yamlRow{}
		row.Note.SetMarshalUnset(presence.UnsetNull)
		b, err := yaml.Marshal(row)
		require.NoError(t, err)
		assert.Contains(t, string(b), "note: null")
	})
}

func TestUnmarshalYAML(t *testing.T) {
	src := `
name: Dupont
age: 42
score: 1.5
ok: true
at: 2026-10-07T12:30:00Z
tags: [a, b]
`
	var row yamlRow
	err := yaml.Unmarshal([]byte(src), &row)

	t.Run("decodes without error", func(t *testing.T) {
		require.NoError(t, err)
	})

	t.Run("string", func(t *testing.T) {
		assert.Equal(t, "Dupont", row.Name.MustGet())
	})

	t.Run("int", func(t *testing.T) {
		assert.Equal(t, int64(42), row.Age.MustGet())
	})

	t.Run("float", func(t *testing.T) {
		assert.InDelta(t, 1.5, row.Score.MustGet(), 0)
	})

	t.Run("bool", func(t *testing.T) {
		assert.True(t, row.OK.MustGet())
	})

	t.Run("timestamp", func(t *testing.T) {
		assert.True(t, time.Date(2026, 10, 7, 12, 30, 0, 0, time.UTC).Equal(row.At.MustGet()))
	})

	t.Run("sequence", func(t *testing.T) {
		assert.Equal(t, []string{"a", "b"}, row.Tags.MustGet())
	})

	t.Run("missing key stays unset", func(t *testing.T) {
		assert.True(t, row.Note.IsUnset())
	})

	// yaml.v3 never calls an unmarshaler for a null node, so a YAML null
	// cannot be told apart from a missing key.
	t.Run("null stays unset", func(t *testing.T) {
		var r yamlRow
		require.NoError(t, yaml.Unmarshal([]byte("age: null"), &r))
		assert.True(t, r.Age.IsUnset())
	})

	t.Run("type mismatch returns an error", func(t *testing.T) {
		var r yamlRow
		err := yaml.Unmarshal([]byte("age: [1, 2]"), &r)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "presence yaml unmarshal")
	})
}

func TestYAMLRoundTrip(t *testing.T) {
	in := yamlRow{
		Name: presence.FromValue("Dupont"),
		Age:  presence.FromValue(int64(42)),
		Tags: presence.FromValue([]string{"a"}),
	}
	b, err := yaml.Marshal(in)
	require.NoError(t, err)

	var out yamlRow
	require.NoError(t, yaml.Unmarshal(b, &out))

	t.Run("values survive", func(t *testing.T) {
		assert.Equal(t, "Dupont", out.Name.MustGet())
		assert.Equal(t, int64(42), out.Age.MustGet())
		assert.Equal(t, []string{"a"}, out.Tags.MustGet())
	})

	t.Run("unset survives", func(t *testing.T) {
		assert.True(t, out.Note.IsUnset())
	})
}
