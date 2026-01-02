package tests

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/pivaldi/presence"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPresenceEdgeCases(t *testing.T) {
	t.Run("SetValueP with nil pointer", func(t *testing.T) {
		test := testedStruct[embeddedStruct]{
			Name: presence.Of[string]{},
		}
		test.Name.SetValueP(nil)

		assert.True(t, test.Name.IsNull(), "SetValueP(nil) should result in NULL value")
	})

	t.Run("SetValueP with value pointer", func(t *testing.T) {
		value := "test value"
		test := testedStruct[embeddedStruct]{
			Name: presence.Of[string]{},
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
		var n presence.PresenceI[string] = &presence.Of[string]{}
		assert.True(t, n.IsUnset())
	})

	t.Run("interface has IsSet method", func(t *testing.T) {
		var n presence.PresenceI[string] = &presence.Of[string]{}
		assert.False(t, n.IsSet())
	})

	t.Run("interface has Unset method", func(t *testing.T) {
		val := presence.FromValue("test")
		var n presence.PresenceI[string] = &val
		n.Unset()
		assert.True(t, n.IsUnset())
	})
}

func TestNullConstructor(t *testing.T) {
	t.Run("Null returns explicitly null value", func(t *testing.T) {
		n := presence.Null[string]()
		assert.True(t, n.IsNull(), "Null() should return IsNull=true")
		assert.False(t, n.IsUnset(), "Null() should return IsUnset=false")
		assert.True(t, n.IsSet(), "Null() should return IsSet=true")
	})

	t.Run("zero value is unset not null", func(t *testing.T) {
		var n presence.Of[string]
		assert.False(t, n.IsNull(), "zero value should not be null")
		assert.True(t, n.IsUnset(), "zero value should be unset")
		assert.False(t, n.IsSet(), "zero value should not be set")
	})
}

func TestUnsetMethod(t *testing.T) {
	t.Run("Unset resets to unset state", func(t *testing.T) {
		n := presence.FromValue("test")
		n.Unset()
		assert.True(t, n.IsUnset())
		assert.False(t, n.IsNull())
		assert.Nil(t, n.GetValue())
	})

	t.Run("Unset on null value", func(t *testing.T) {
		n := presence.Null[int]()
		n.Unset()
		assert.True(t, n.IsUnset())
		assert.False(t, n.IsNull())
	})
}

func TestScanNullBehavior(t *testing.T) {
	t.Run("scan null with ScanNullAsNull", func(t *testing.T) {
		n := presence.Of[string]{}
		n.SetScanNull(presence.ScanNullAsNull)
		err := n.Scan(nil)
		require.NoError(t, err)
		assert.True(t, n.IsNull())
		assert.False(t, n.IsUnset())
	})

	t.Run("scan null with ScanNullAsUnset", func(t *testing.T) {
		n := presence.Of[string]{}
		n.SetScanNull(presence.ScanNullAsUnset)
		err := n.Scan(nil)
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.True(t, n.IsUnset())
	})

	t.Run("scan value sets isSet", func(t *testing.T) {
		n := presence.Of[string]{}
		err := n.Scan("test")
		require.NoError(t, err)
		assert.True(t, n.IsSet())
		assert.False(t, n.IsNull())
		assert.False(t, n.IsUnset())
	})
}

// Tests for Get method
func TestGet(t *testing.T) {
	t.Run("Get on value returns value and true", func(t *testing.T) {
		n := presence.FromValue("hello")
		v, ok := n.Get()
		assert.True(t, ok)
		assert.Equal(t, "hello", v)
	})

	t.Run("Get on null returns zero and false", func(t *testing.T) {
		n := presence.Null[string]()
		v, ok := n.Get()
		assert.False(t, ok)
		assert.Equal(t, "", v)
	})

	t.Run("Get on unset returns zero and false", func(t *testing.T) {
		var n presence.Of[int]
		v, ok := n.Get()
		assert.False(t, ok)
		assert.Equal(t, 0, v)
	})

	t.Run("Get on nil receiver returns zero and false", func(t *testing.T) {
		var n *presence.Of[string]
		v, ok := n.Get()
		assert.False(t, ok)
		assert.Equal(t, "", v)
	})
}

// Tests for GetOr method
func TestGetOr(t *testing.T) {
	t.Run("GetOr on value returns value", func(t *testing.T) {
		n := presence.FromValue(42)
		v := n.GetOr(100)
		assert.Equal(t, 42, v)
	})

	t.Run("GetOr on null returns default", func(t *testing.T) {
		n := presence.Null[int]()
		v := n.GetOr(100)
		assert.Equal(t, 100, v)
	})

	t.Run("GetOr on unset returns default", func(t *testing.T) {
		var n presence.Of[string]
		v := n.GetOr("default")
		assert.Equal(t, "default", v)
	})

	t.Run("GetOr on nil receiver returns default", func(t *testing.T) {
		var n *presence.Of[float64]
		v := n.GetOr(3.14)
		assert.Equal(t, 3.14, v)
	})
}

// Tests for MustGet method
func TestMustGet(t *testing.T) {
	t.Run("MustGet on value returns value", func(t *testing.T) {
		n := presence.FromValue("test")
		v := n.MustGet()
		assert.Equal(t, "test", v)
	})

	t.Run("MustGet on null panics", func(t *testing.T) {
		n := presence.Null[int]()
		assert.Panics(t, func() {
			n.MustGet()
		})
	})

	t.Run("MustGet on unset panics", func(t *testing.T) {
		var n presence.Of[string]
		assert.Panics(t, func() {
			n.MustGet()
		})
	})

	t.Run("MustGet on nil receiver panics", func(t *testing.T) {
		var n *presence.Of[int]
		assert.Panics(t, func() {
			n.MustGet()
		})
	})
}

// Tests for Ptr method
func TestPtr(t *testing.T) {
	t.Run("Ptr on value returns pointer to value", func(t *testing.T) {
		n := presence.FromValue(42)
		ptr := n.Ptr()
		require.NotNil(t, ptr)
		assert.Equal(t, 42, *ptr)
	})

	t.Run("Ptr on null returns nil", func(t *testing.T) {
		n := presence.Null[string]()
		ptr := n.Ptr()
		assert.Nil(t, ptr)
	})

	t.Run("Ptr on unset returns nil", func(t *testing.T) {
		var n presence.Of[int]
		ptr := n.Ptr()
		assert.Nil(t, ptr)
	})

	t.Run("Ptr on nil receiver returns nil", func(t *testing.T) {
		var n *presence.Of[string]
		ptr := n.Ptr()
		assert.Nil(t, ptr)
	})
}

// Tests for IsValue method
func TestIsValue(t *testing.T) {
	t.Run("IsValue on value returns true", func(t *testing.T) {
		n := presence.FromValue("hello")
		assert.True(t, n.IsValue())
	})

	t.Run("IsValue on null returns false", func(t *testing.T) {
		n := presence.Null[int]()
		assert.False(t, n.IsValue())
	})

	t.Run("IsValue on unset returns false", func(t *testing.T) {
		var n presence.Of[string]
		assert.False(t, n.IsValue())
	})

	t.Run("IsValue on nil receiver returns false", func(t *testing.T) {
		var n *presence.Of[int]
		assert.False(t, n.IsValue())
	})
}

// Tests for Map function
func TestMap(t *testing.T) {
	t.Run("Map on value transforms value", func(t *testing.T) {
		n := presence.FromValue(42)
		result := presence.Map(n, func(v int) string {
			return fmt.Sprintf("value: %d", v)
		})
		assert.True(t, result.IsValue())
		assert.Equal(t, "value: 42", result.MustGet())
	})

	t.Run("Map on null returns null", func(t *testing.T) {
		n := presence.Null[int]()
		result := presence.Map(n, func(v int) string {
			return fmt.Sprintf("value: %d", v)
		})
		assert.True(t, result.IsNull())
		assert.False(t, result.IsUnset())
	})

	t.Run("Map on unset returns unset", func(t *testing.T) {
		var n presence.Of[int]
		result := presence.Map(n, func(v int) string {
			return fmt.Sprintf("value: %d", v)
		})
		assert.True(t, result.IsUnset())
		assert.False(t, result.IsNull())
	})

	t.Run("Map with type conversion", func(t *testing.T) {
		n := presence.FromValue("123")
		result := presence.Map(n, func(v string) int {
			i, _ := strconv.Atoi(v)
			return i
		})
		assert.True(t, result.IsValue())
		assert.Equal(t, 123, result.MustGet())
	})
}

// Tests for MapOr function
func TestMapOr(t *testing.T) {
	t.Run("MapOr on value transforms value", func(t *testing.T) {
		n := presence.FromValue(10)
		result := presence.MapOr(n, "default", func(v int) string {
			return fmt.Sprintf("got %d", v)
		})
		assert.Equal(t, "got 10", result)
	})

	t.Run("MapOr on null returns default", func(t *testing.T) {
		n := presence.Null[int]()
		result := presence.MapOr(n, "default", func(v int) string {
			return fmt.Sprintf("got %d", v)
		})
		assert.Equal(t, "default", result)
	})

	t.Run("MapOr on unset returns default", func(t *testing.T) {
		var n presence.Of[string]
		result := presence.MapOr(n, 0, func(v string) int {
			return len(v)
		})
		assert.Equal(t, 0, result)
	})
}

// Tests for FlatMap function
func TestFlatMap(t *testing.T) {
	// Helper function that returns presence
	parsePositive := func(s string) presence.Of[int] {
		i, err := strconv.Atoi(s)
		if err != nil || i < 0 {
			return presence.Null[int]()
		}
		return presence.FromValue(i)
	}

	t.Run("FlatMap on value with successful transform", func(t *testing.T) {
		n := presence.FromValue("42")
		result := presence.FlatMap(n, parsePositive)
		assert.True(t, result.IsValue())
		assert.Equal(t, 42, result.MustGet())
	})

	t.Run("FlatMap on value with failing transform", func(t *testing.T) {
		n := presence.FromValue("-5")
		result := presence.FlatMap(n, parsePositive)
		assert.True(t, result.IsNull())
	})

	t.Run("FlatMap on null returns null", func(t *testing.T) {
		n := presence.Null[string]()
		result := presence.FlatMap(n, parsePositive)
		assert.True(t, result.IsNull())
		assert.False(t, result.IsUnset())
	})

	t.Run("FlatMap on unset returns unset", func(t *testing.T) {
		var n presence.Of[string]
		result := presence.FlatMap(n, parsePositive)
		assert.True(t, result.IsUnset())
		assert.False(t, result.IsNull())
	})

	t.Run("FlatMap chaining", func(t *testing.T) {
		n := presence.FromValue("100")
		result := presence.FlatMap(n, func(s string) presence.Of[int] {
			i, _ := strconv.Atoi(s)
			return presence.FromValue(i)
		})
		result2 := presence.FlatMap(result, func(i int) presence.Of[string] {
			if i > 50 {
				return presence.FromValue("large")
			}
			return presence.FromValue("small")
		})
		assert.Equal(t, "large", result2.MustGet())
	})
}

// Tests for Filter function
func TestFilter(t *testing.T) {
	isPositive := func(n int) bool { return n > 0 }
	isLong := func(s string) bool { return len(s) > 5 }

	t.Run("Filter on value that passes predicate", func(t *testing.T) {
		n := presence.FromValue(42)
		result := presence.Filter(n, isPositive)
		assert.True(t, result.IsValue())
		assert.Equal(t, 42, result.MustGet())
	})

	t.Run("Filter on value that fails predicate", func(t *testing.T) {
		n := presence.FromValue(-5)
		result := presence.Filter(n, isPositive)
		assert.True(t, result.IsNull())
		assert.False(t, result.IsUnset())
	})

	t.Run("Filter on null returns null", func(t *testing.T) {
		n := presence.Null[int]()
		result := presence.Filter(n, isPositive)
		assert.True(t, result.IsNull())
	})

	t.Run("Filter on unset returns unset", func(t *testing.T) {
		var n presence.Of[int]
		result := presence.Filter(n, isPositive)
		assert.True(t, result.IsUnset())
	})

	t.Run("Filter with string predicate", func(t *testing.T) {
		n := presence.FromValue("hello world")
		result := presence.Filter(n, isLong)
		assert.True(t, result.IsValue())
		assert.Equal(t, "hello world", result.MustGet())
	})

	t.Run("Filter with failing string predicate", func(t *testing.T) {
		n := presence.FromValue("hi")
		result := presence.Filter(n, isLong)
		assert.True(t, result.IsNull())
	})
}

// Tests for Or function
func TestOr(t *testing.T) {
	t.Run("Or returns first value", func(t *testing.T) {
		a := presence.FromValue("first")
		b := presence.FromValue("second")
		c := presence.FromValue("third")
		result := presence.Or(a, b, c)
		assert.Equal(t, "first", result.MustGet())
	})

	t.Run("Or skips null and returns first value", func(t *testing.T) {
		a := presence.Null[string]()
		b := presence.FromValue("second")
		c := presence.FromValue("third")
		result := presence.Or(a, b, c)
		assert.Equal(t, "second", result.MustGet())
	})

	t.Run("Or skips unset and returns first value", func(t *testing.T) {
		var a presence.Of[int]
		b := presence.FromValue(42)
		result := presence.Or(a, b)
		assert.Equal(t, 42, result.MustGet())
	})

	t.Run("Or skips null and unset", func(t *testing.T) {
		var a presence.Of[string]
		b := presence.Null[string]()
		c := presence.FromValue("third")
		result := presence.Or(a, b, c)
		assert.Equal(t, "third", result.MustGet())
	})

	t.Run("Or returns null when all are null or unset", func(t *testing.T) {
		var a presence.Of[int]
		b := presence.Null[int]()
		var c presence.Of[int]
		result := presence.Or(a, b, c)
		assert.True(t, result.IsNull())
	})

	t.Run("Or with empty args returns null", func(t *testing.T) {
		result := presence.Or[string]()
		assert.True(t, result.IsNull())
	})

	t.Run("Or with single value", func(t *testing.T) {
		a := presence.FromValue(100)
		result := presence.Or(a)
		assert.Equal(t, 100, result.MustGet())
	})
}

// Tests for FromPtr function
func TestFromPtr(t *testing.T) {
	t.Run("FromPtr with non-nil pointer", func(t *testing.T) {
		value := "hello"
		n := presence.FromPtr(&value)
		assert.True(t, n.IsValue())
		assert.Equal(t, "hello", n.MustGet())
	})

	t.Run("FromPtr with nil pointer", func(t *testing.T) {
		var ptr *int
		n := presence.FromPtr(ptr)
		assert.True(t, n.IsNull())
		assert.False(t, n.IsUnset())
	})

	t.Run("FromPtr with zero value pointer", func(t *testing.T) {
		zero := 0
		n := presence.FromPtr(&zero)
		assert.True(t, n.IsValue())
		assert.Equal(t, 0, n.MustGet())
	})

	t.Run("FromPtr with empty string pointer", func(t *testing.T) {
		empty := ""
		n := presence.FromPtr(&empty)
		assert.True(t, n.IsValue())
		assert.Equal(t, "", n.MustGet())
	})
}

// Tests for FromBool function
func TestFromBool(t *testing.T) {
	t.Run("FromBool with true", func(t *testing.T) {
		n := presence.FromBool("value", true)
		assert.True(t, n.IsValue())
		assert.Equal(t, "value", n.MustGet())
	})

	t.Run("FromBool with false", func(t *testing.T) {
		n := presence.FromBool("value", false)
		assert.True(t, n.IsNull())
		assert.False(t, n.IsUnset())
	})

	t.Run("FromBool with zero value and true", func(t *testing.T) {
		n := presence.FromBool(0, true)
		assert.True(t, n.IsValue())
		assert.Equal(t, 0, n.MustGet())
	})

	t.Run("FromBool from map lookup", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2}
		v, ok := m["a"]
		n := presence.FromBool(v, ok)
		assert.True(t, n.IsValue())
		assert.Equal(t, 1, n.MustGet())

		v, ok = m["c"]
		n = presence.FromBool(v, ok)
		assert.True(t, n.IsNull())
	})
}

// Tests for interface compliance
func TestInterfaceCompliance(t *testing.T) {
	t.Run("Of implements PresenceI with new methods", func(t *testing.T) {
		val := presence.FromValue("test")
		var n presence.PresenceI[string] = &val

		// Test IsValue
		assert.True(t, n.IsValue())

		// Test Get
		v, ok := n.Get()
		assert.True(t, ok)
		assert.Equal(t, "test", v)

		// Test GetOr
		assert.Equal(t, "test", n.GetOr("default"))

		// Test MustGet
		assert.Equal(t, "test", n.MustGet())

		// Test Ptr
		ptr := n.Ptr()
		require.NotNil(t, ptr)
		assert.Equal(t, "test", *ptr)
	})
}

// Tests for combined functional operations
func TestCombinedOperations(t *testing.T) {
	t.Run("Map then Filter", func(t *testing.T) {
		n := presence.FromValue("hello")
		mapped := presence.Map(n, func(s string) int {
			return len(s)
		})
		filtered := presence.Filter(mapped, func(i int) bool {
			return i > 3
		})
		assert.True(t, filtered.IsValue())
		assert.Equal(t, 5, filtered.MustGet())
	})

	t.Run("Filter then Map", func(t *testing.T) {
		n := presence.FromValue(10)
		filtered := presence.Filter(n, func(i int) bool {
			return i > 5
		})
		mapped := presence.Map(filtered, func(i int) string {
			return fmt.Sprintf("number: %d", i)
		})
		assert.Equal(t, "number: 10", mapped.MustGet())
	})

	t.Run("Chain with null propagation", func(t *testing.T) {
		n := presence.FromValue(-5)
		// Filter out negative numbers
		filtered := presence.Filter(n, func(i int) bool {
			return i > 0
		})
		// This should not execute since filtered is null
		mapped := presence.Map(filtered, func(i int) string {
			return fmt.Sprintf("positive: %d", i)
		})
		assert.True(t, mapped.IsNull())
	})

	t.Run("Or with Map fallback", func(t *testing.T) {
		primary := presence.Null[string]()
		secondary := presence.FromValue("backup")

		result := presence.Or(primary, secondary)
		mapped := presence.Map(result, strings.ToUpper)
		assert.Equal(t, "BACKUP", mapped.MustGet())
	})
}
