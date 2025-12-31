package tests

import (
	"testing"

	"github.com/pivaldi/nullable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNullableEdgeCases(t *testing.T) {
	t.Run("SetValueP with nil pointer", func(t *testing.T) {
		test := testedStruct[embeddedStruct]{
			Name: nullable.Of[string]{},
		}
		test.Name.SetValueP(nil)

		assert.True(t, test.Name.IsNull(), "SetValueP(nil) should result in NULL value")
	})

	t.Run("SetValueP with value pointer", func(t *testing.T) {
		value := "test value"
		test := testedStruct[embeddedStruct]{
			Name: nullable.Of[string]{},
		}
		test.Name.SetValueP(&value)

		require.False(t, test.Name.IsNull(), "SetValueP(&value) should not result in NULL")
		assert.Equal(t, value, *test.Name.GetValue())
	})

	t.Run("IsNull on zero value", func(t *testing.T) {
		var test testedStruct[embeddedStruct]
		assert.False(t, test.Name.IsNull(), "Zero value should be unset, not null")
		assert.True(t, test.Name.IsUnset(), "Zero value should be unset")
	})
}

func TestThreeStateInterface(t *testing.T) {
	t.Run("interface has IsUnset method", func(t *testing.T) {
		var n nullable.NullableI[string] = &nullable.Of[string]{}
		assert.True(t, n.IsUnset())
	})

	t.Run("interface has IsSet method", func(t *testing.T) {
		var n nullable.NullableI[string] = &nullable.Of[string]{}
		assert.False(t, n.IsSet())
	})

	t.Run("interface has Unset method", func(t *testing.T) {
		val := nullable.FromValue("test")
		var n nullable.NullableI[string] = &val
		n.Unset()
		assert.True(t, n.IsUnset())
	})
}

func TestNullConstructor(t *testing.T) {
	t.Run("Null returns explicitly null value", func(t *testing.T) {
		n := nullable.Null[string]()
		assert.True(t, n.IsNull(), "Null() should return IsNull=true")
		assert.False(t, n.IsUnset(), "Null() should return IsUnset=false")
		assert.True(t, n.IsSet(), "Null() should return IsSet=true")
	})

	t.Run("zero value is unset not null", func(t *testing.T) {
		var n nullable.Of[string]
		assert.False(t, n.IsNull(), "zero value should not be null")
		assert.True(t, n.IsUnset(), "zero value should be unset")
		assert.False(t, n.IsSet(), "zero value should not be set")
	})
}

func TestUnsetMethod(t *testing.T) {
	t.Run("Unset resets to unset state", func(t *testing.T) {
		n := nullable.FromValue("test")
		n.Unset()
		assert.True(t, n.IsUnset())
		assert.False(t, n.IsNull())
		assert.Nil(t, n.GetValue())
	})

	t.Run("Unset on null value", func(t *testing.T) {
		n := nullable.Null[int]()
		n.Unset()
		assert.True(t, n.IsUnset())
		assert.False(t, n.IsNull())
	})
}

func TestScanNullBehavior(t *testing.T) {
	t.Run("scan null with ScanNullAsNull", func(t *testing.T) {
		n := nullable.Of[string]{}
		n.SetScanNull(nullable.ScanNullAsNull)
		err := n.Scan(nil)
		require.NoError(t, err)
		assert.True(t, n.IsNull())
		assert.False(t, n.IsUnset())
	})

	t.Run("scan null with ScanNullAsUnset", func(t *testing.T) {
		n := nullable.Of[string]{}
		n.SetScanNull(nullable.ScanNullAsUnset)
		err := n.Scan(nil)
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.True(t, n.IsUnset())
	})

	t.Run("scan value sets isSet", func(t *testing.T) {
		n := nullable.Of[string]{}
		err := n.Scan("test")
		require.NoError(t, err)
		assert.True(t, n.IsSet())
		assert.False(t, n.IsNull())
		assert.False(t, n.IsUnset())
	})
}
