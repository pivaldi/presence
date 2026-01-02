package tests

import (
	"encoding/json"
	"math"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pivaldi/presence"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalJSON_NullValues(t *testing.T) {
	t.Run("null string", func(t *testing.T) {
		n := presence.Null[string]()
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("null"), data)
	})

	t.Run("null int", func(t *testing.T) {
		n := presence.Null[int]()
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("null"), data)
	})

	t.Run("null bool", func(t *testing.T) {
		n := presence.Null[bool]()
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("null"), data)
	})

	t.Run("null float64", func(t *testing.T) {
		n := presence.Null[float64]()
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("null"), data)
	})

	t.Run("null UUID", func(t *testing.T) {
		n := presence.Null[uuid.UUID]()
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("null"), data)
	})
}

func TestMarshalJSON_PrimitiveTypes(t *testing.T) {
	t.Run("string value", func(t *testing.T) {
		n := presence.FromValue("hello world")
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte(`"hello world"`), data)
	})

	t.Run("empty string", func(t *testing.T) {
		n := presence.FromValue("")
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte(`""`), data)
	})

	t.Run("int value", func(t *testing.T) {
		n := presence.FromValue(42)
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("42"), data)
	})

	t.Run("int16 value", func(t *testing.T) {
		n := presence.FromValue(int16(123))
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("123"), data)
	})

	t.Run("int32 value", func(t *testing.T) {
		n := presence.FromValue(int32(456))
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("456"), data)
	})

	t.Run("int64 value", func(t *testing.T) {
		n := presence.FromValue(int64(789))
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("789"), data)
	})

	t.Run("zero int", func(t *testing.T) {
		n := presence.FromValue(0)
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("0"), data)
	})

	t.Run("negative int", func(t *testing.T) {
		n := presence.FromValue(-42)
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("-42"), data)
	})

	t.Run("bool true", func(t *testing.T) {
		n := presence.FromValue(true)
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("true"), data)
	})

	t.Run("bool false", func(t *testing.T) {
		n := presence.FromValue(false)
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("false"), data)
	})

	t.Run("float64 value", func(t *testing.T) {
		n := presence.FromValue(3.14159)
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("3.14159"), data)
	})

	t.Run("float64 zero", func(t *testing.T) {
		n := presence.FromValue(0.0)
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("0"), data)
	})

	t.Run("float64 negative", func(t *testing.T) {
		n := presence.FromValue(-2.5)
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("-2.5"), data)
	})
}

func TestMarshalJSON_SpecialFloatValues(t *testing.T) {
	t.Run("float64 NaN", func(t *testing.T) {
		n := presence.FromValue(math.NaN())
		_, err := n.MarshalJSON()
		// JSON doesn't support NaN, so this should error
		assert.Error(t, err)
	})

	t.Run("float64 positive infinity", func(t *testing.T) {
		n := presence.FromValue(math.Inf(1))
		_, err := n.MarshalJSON()
		// JSON doesn't support Inf, so this should error
		assert.Error(t, err)
	})

	t.Run("float64 negative infinity", func(t *testing.T) {
		n := presence.FromValue(math.Inf(-1))
		_, err := n.MarshalJSON()
		// JSON doesn't support -Inf, so this should error
		assert.Error(t, err)
	})
}

func TestMarshalJSON_UUID(t *testing.T) {
	t.Run("valid UUID", func(t *testing.T) {
		testUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
		n := presence.FromValue(testUUID)
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte(`"550e8400-e29b-41d4-a716-446655440000"`), data)
	})

	t.Run("zero UUID", func(t *testing.T) {
		n := presence.FromValue(uuid.UUID{})
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte(`"00000000-0000-0000-0000-000000000000"`), data)
	})
}

func TestMarshalJSON_JSONType(t *testing.T) {
	t.Run("simple map", func(t *testing.T) {
		obj := map[string]any{"key": "value", "number": 42}
		n := presence.FromValue(obj)
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		// Can't assert exact JSON due to map ordering, so unmarshal and compare
		var result map[string]any
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)
		assert.Equal(t, "value", result["key"])
		assert.Equal(t, float64(42), result["number"])
	})

	t.Run("nested structure", func(t *testing.T) {
		obj := map[string]any{
			"nested": map[string]any{
				"inner": "value",
			},
		}
		n := presence.FromValue(obj)
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		var result map[string]any
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)
		nested := result["nested"].(map[string]any)
		assert.Equal(t, "value", nested["inner"])
	})

	t.Run("array", func(t *testing.T) {
		obj := []any{1, 2, 3, "four"}
		n := presence.FromValue(obj)
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.JSONEq(t, `[1,2,3,"four"]`, string(data))
	})
}

func TestUnmarshalJSON_NullValues(t *testing.T) {
	t.Run("null keyword to string", func(t *testing.T) {
		var n presence.Of[string]
		err := n.UnmarshalJSON([]byte("null"))
		require.NoError(t, err)
		assert.True(t, n.IsNull())
		assert.Nil(t, n.GetValue())
	})

	t.Run("null keyword to int", func(t *testing.T) {
		var n presence.Of[int]
		err := n.UnmarshalJSON([]byte("null"))
		require.NoError(t, err)
		assert.True(t, n.IsNull())
	})

	t.Run("nil byte slice", func(t *testing.T) {
		var n presence.Of[string]
		err := n.UnmarshalJSON(nil)
		require.NoError(t, err)
		assert.True(t, n.IsNull())
	})

	t.Run("null to previously set value", func(t *testing.T) {
		n := presence.FromValue("previous value")
		err := n.UnmarshalJSON([]byte("null"))
		require.NoError(t, err)
		assert.True(t, n.IsNull())
		assert.Nil(t, n.GetValue())
	})
}

func TestUnmarshalJSON_PrimitiveTypes(t *testing.T) {
	t.Run("string value", func(t *testing.T) {
		var n presence.Of[string]
		err := n.UnmarshalJSON([]byte(`"hello world"`))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, "hello world", *n.GetValue())
	})

	t.Run("empty string", func(t *testing.T) {
		var n presence.Of[string]
		err := n.UnmarshalJSON([]byte(`""`))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, "", *n.GetValue())
	})

	t.Run("string with quotes", func(t *testing.T) {
		var n presence.Of[string]
		err := n.UnmarshalJSON([]byte(`"say \"hello\""`))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, `say "hello"`, *n.GetValue())
	})

	t.Run("int value", func(t *testing.T) {
		var n presence.Of[int]
		err := n.UnmarshalJSON([]byte("42"))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, 42, *n.GetValue())
	})

	t.Run("int16 value", func(t *testing.T) {
		var n presence.Of[int16]
		err := n.UnmarshalJSON([]byte("123"))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, int16(123), *n.GetValue())
	})

	t.Run("int32 value", func(t *testing.T) {
		var n presence.Of[int32]
		err := n.UnmarshalJSON([]byte("456"))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, int32(456), *n.GetValue())
	})

	t.Run("int64 value", func(t *testing.T) {
		var n presence.Of[int64]
		err := n.UnmarshalJSON([]byte("789"))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, int64(789), *n.GetValue())
	})

	t.Run("zero int", func(t *testing.T) {
		var n presence.Of[int]
		err := n.UnmarshalJSON([]byte("0"))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, 0, *n.GetValue())
	})

	t.Run("negative int", func(t *testing.T) {
		var n presence.Of[int]
		err := n.UnmarshalJSON([]byte("-42"))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, -42, *n.GetValue())
	})

	t.Run("bool true", func(t *testing.T) {
		var n presence.Of[bool]
		err := n.UnmarshalJSON([]byte("true"))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, true, *n.GetValue())
	})

	t.Run("bool false", func(t *testing.T) {
		var n presence.Of[bool]
		err := n.UnmarshalJSON([]byte("false"))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, false, *n.GetValue())
	})

	t.Run("float64 value", func(t *testing.T) {
		var n presence.Of[float64]
		err := n.UnmarshalJSON([]byte("3.14159"))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, 3.14159, *n.GetValue())
	})

	t.Run("float64 zero", func(t *testing.T) {
		var n presence.Of[float64]
		err := n.UnmarshalJSON([]byte("0.0"))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, 0.0, *n.GetValue())
	})

	t.Run("float64 negative", func(t *testing.T) {
		var n presence.Of[float64]
		err := n.UnmarshalJSON([]byte("-2.5"))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, -2.5, *n.GetValue())
	})

	t.Run("float64 scientific notation", func(t *testing.T) {
		var n presence.Of[float64]
		err := n.UnmarshalJSON([]byte("1.23e-4"))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.InDelta(t, 0.000123, *n.GetValue(), 0.0000001)
	})
}

func TestUnmarshalJSON_UUID(t *testing.T) {
	t.Run("valid UUID", func(t *testing.T) {
		var n presence.Of[uuid.UUID]
		err := n.UnmarshalJSON([]byte(`"550e8400-e29b-41d4-a716-446655440000"`))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		expected := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
		assert.Equal(t, expected, *n.GetValue())
	})

	t.Run("zero UUID", func(t *testing.T) {
		var n presence.Of[uuid.UUID]
		err := n.UnmarshalJSON([]byte(`"00000000-0000-0000-0000-000000000000"`))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, uuid.UUID{}, *n.GetValue())
	})

	t.Run("invalid UUID format", func(t *testing.T) {
		var n presence.Of[uuid.UUID]
		err := n.UnmarshalJSON([]byte(`"not-a-uuid"`))
		assert.Error(t, err)
	})
}

func TestUnmarshalJSON_InvalidJSON(t *testing.T) {
	t.Run("invalid JSON for string", func(t *testing.T) {
		var n presence.Of[string]
		err := n.UnmarshalJSON([]byte(`not valid json`))
		assert.Error(t, err)
	})

	t.Run("invalid JSON for int", func(t *testing.T) {
		var n presence.Of[int]
		err := n.UnmarshalJSON([]byte(`"not a number"`))
		assert.Error(t, err)
	})

	t.Run("invalid JSON for bool", func(t *testing.T) {
		var n presence.Of[bool]
		err := n.UnmarshalJSON([]byte(`"not a bool"`))
		assert.Error(t, err)
	})

	t.Run("invalid JSON for float", func(t *testing.T) {
		var n presence.Of[float64]
		err := n.UnmarshalJSON([]byte(`"not a float"`))
		assert.Error(t, err)
	})

	t.Run("number overflow for int16", func(t *testing.T) {
		var n presence.Of[int16]
		err := n.UnmarshalJSON([]byte("100000"))
		assert.Error(t, err)
	})

	t.Run("number overflow for int32", func(t *testing.T) {
		var n presence.Of[int32]
		err := n.UnmarshalJSON([]byte("10000000000"))
		assert.Error(t, err)
	})
}

func TestUnmarshalJSON_JSONType(t *testing.T) {
	t.Run("simple map", func(t *testing.T) {
		var n presence.Of[any]
		err := n.UnmarshalJSON([]byte(`{"key":"value","number":42}`))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		result := (*n.GetValue()).(map[string]any)
		assert.Equal(t, "value", result["key"])
		assert.Equal(t, float64(42), result["number"])
	})

	t.Run("nested structure", func(t *testing.T) {
		var n presence.Of[any]
		err := n.UnmarshalJSON([]byte(`{"nested":{"inner":"value"}}`))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		result := (*n.GetValue()).(map[string]any)
		nested := result["nested"].(map[string]any)
		assert.Equal(t, "value", nested["inner"])
	})

	t.Run("array", func(t *testing.T) {
		var n presence.Of[any]
		err := n.UnmarshalJSON([]byte(`[1,2,3,"four"]`))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		result := (*n.GetValue()).([]any)
		assert.Len(t, result, 4)
		assert.Equal(t, float64(1), result[0])
		assert.Equal(t, "four", result[3])
	})
}

func TestMarshalUnmarshal_RoundTrip(t *testing.T) {
	t.Run("string round trip", func(t *testing.T) {
		original := presence.FromValue("test value")
		data, err := original.MarshalJSON()
		require.NoError(t, err)

		var restored presence.Of[string]
		err = restored.UnmarshalJSON(data)
		require.NoError(t, err)
		assert.Equal(t, *original.GetValue(), *restored.GetValue())
	})

	t.Run("int round trip", func(t *testing.T) {
		original := presence.FromValue(42)
		data, err := original.MarshalJSON()
		require.NoError(t, err)

		var restored presence.Of[int]
		err = restored.UnmarshalJSON(data)
		require.NoError(t, err)
		assert.Equal(t, *original.GetValue(), *restored.GetValue())
	})

	t.Run("bool round trip", func(t *testing.T) {
		original := presence.FromValue(true)
		data, err := original.MarshalJSON()
		require.NoError(t, err)

		var restored presence.Of[bool]
		err = restored.UnmarshalJSON(data)
		require.NoError(t, err)
		assert.Equal(t, *original.GetValue(), *restored.GetValue())
	})

	t.Run("float64 round trip", func(t *testing.T) {
		original := presence.FromValue(3.14159)
		data, err := original.MarshalJSON()
		require.NoError(t, err)

		var restored presence.Of[float64]
		err = restored.UnmarshalJSON(data)
		require.NoError(t, err)
		assert.Equal(t, *original.GetValue(), *restored.GetValue())
	})

	t.Run("UUID round trip", func(t *testing.T) {
		original := presence.FromValue(uuid.New())
		data, err := original.MarshalJSON()
		require.NoError(t, err)

		var restored presence.Of[uuid.UUID]
		err = restored.UnmarshalJSON(data)
		require.NoError(t, err)
		assert.Equal(t, *original.GetValue(), *restored.GetValue())
	})

	t.Run("null round trip", func(t *testing.T) {
		original := presence.Null[string]()
		data, err := original.MarshalJSON()
		require.NoError(t, err)

		var restored presence.Of[string]
		err = restored.UnmarshalJSON(data)
		require.NoError(t, err)
		assert.True(t, restored.IsNull())
	})

	t.Run("JSON type round trip", func(t *testing.T) {
		obj := map[string]any{"key": "value", "number": float64(42)}
		original := presence.FromValue(obj)
		data, err := original.MarshalJSON()
		require.NoError(t, err)

		var restored presence.Of[any]
		err = restored.UnmarshalJSON(data)
		require.NoError(t, err)
		result := (*restored.GetValue()).(map[string]any)
		assert.Equal(t, "value", result["key"])
		assert.Equal(t, float64(42), result["number"])
	})
}

func TestMarshalUnmarshal_InStructs(t *testing.T) {
	type TestStruct struct {
		Name   presence.Of[string]  `json:"name"`
		Age    presence.Of[int]     `json:"age"`
		Active presence.Of[bool]    `json:"active"`
		Score  presence.Of[float64] `json:"score"`
	}

	t.Run("struct with all values", func(t *testing.T) {
		original := TestStruct{
			Name:   presence.FromValue("John"),
			Age:    presence.FromValue(30),
			Active: presence.FromValue(true),
			Score:  presence.FromValue(95.5),
		}

		data, err := json.Marshal(original)
		require.NoError(t, err)

		var restored TestStruct
		err = json.Unmarshal(data, &restored)
		require.NoError(t, err)

		assert.Equal(t, *original.Name.GetValue(), *restored.Name.GetValue())
		assert.Equal(t, *original.Age.GetValue(), *restored.Age.GetValue())
		assert.Equal(t, *original.Active.GetValue(), *restored.Active.GetValue())
		assert.Equal(t, *original.Score.GetValue(), *restored.Score.GetValue())
	})

	t.Run("struct with null values", func(t *testing.T) {
		original := TestStruct{
			Name:   presence.Null[string](),
			Age:    presence.Null[int](),
			Active: presence.Null[bool](),
			Score:  presence.Null[float64](),
		}

		data, err := json.Marshal(original)
		require.NoError(t, err)
		assert.JSONEq(t, `{"name":null,"age":null,"active":null,"score":null}`, string(data))

		var restored TestStruct
		err = json.Unmarshal(data, &restored)
		require.NoError(t, err)

		assert.True(t, restored.Name.IsNull())
		assert.True(t, restored.Age.IsNull())
		assert.True(t, restored.Active.IsNull())
		assert.True(t, restored.Score.IsNull())
	})

	t.Run("struct with mixed null and non-null", func(t *testing.T) {
		original := TestStruct{
			Name:   presence.FromValue("Jane"),
			Age:    presence.Null[int](),
			Active: presence.FromValue(false),
			Score:  presence.Null[float64](),
		}

		data, err := json.Marshal(original)
		require.NoError(t, err)

		var restored TestStruct
		err = json.Unmarshal(data, &restored)
		require.NoError(t, err)

		assert.Equal(t, "Jane", *restored.Name.GetValue())
		assert.True(t, restored.Age.IsNull())
		assert.Equal(t, false, *restored.Active.GetValue())
		assert.True(t, restored.Score.IsNull())
	})
}

func TestUnmarshalJSON_OverwritingExisting(t *testing.T) {
	t.Run("overwrite value with new value", func(t *testing.T) {
		n := presence.FromValue("original")
		err := n.UnmarshalJSON([]byte(`"new value"`))
		require.NoError(t, err)
		assert.Equal(t, "new value", *n.GetValue())
	})

	t.Run("overwrite value with null", func(t *testing.T) {
		n := presence.FromValue(42)
		err := n.UnmarshalJSON([]byte("null"))
		require.NoError(t, err)
		assert.True(t, n.IsNull())
	})

	t.Run("overwrite null with value", func(t *testing.T) {
		n := presence.Null[int]()
		err := n.UnmarshalJSON([]byte("123"))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, 123, *n.GetValue())
	})
}

func TestMarshalUnmarshal_ComplexStructures(t *testing.T) {
	// Define complex nested structures
	type Address struct {
		Street     presence.Of[string]  `json:"street"`
		City       presence.Of[string]  `json:"city"`
		PostalCode presence.Of[string]  `json:"postalCode"`
		Country    presence.Of[string]  `json:"country"`
		Verified   presence.Of[bool]    `json:"verified"`
		Lat        presence.Of[float64] `json:"lat"`
		Lng        presence.Of[float64] `json:"lng"`
	}

	type ContactInfo struct {
		Email       presence.Of[string]  `json:"email"`
		Phone       presence.Of[string]  `json:"phone"`
		Address     presence.Of[Address] `json:"address"`
		IsPrimary   presence.Of[bool]    `json:"isPrimary"`
		LastUpdated presence.Of[int64]   `json:"lastUpdated"`
	}

	type Metadata struct {
		Tags        presence.Of[[]string]          `json:"tags"`
		Properties  presence.Of[map[string]string] `json:"properties"`
		Version     presence.Of[int]               `json:"version"`
		IsActive    presence.Of[bool]              `json:"isActive"`
		CreatedBy   presence.Of[string]            `json:"createdBy"`
		CreatedByID presence.Of[uuid.UUID]         `json:"createdById"`
	}

	type Profile struct {
		Bio         presence.Of[string]         `json:"bio"`
		Website     presence.Of[string]         `json:"website"`
		AvatarURL   presence.Of[string]         `json:"avatarUrl"`
		Contacts    presence.Of[[]ContactInfo]  `json:"contacts"`
		Preferences presence.Of[map[string]any] `json:"preferences"`
		Metadata    presence.Of[Metadata]       `json:"metadata"`
		Score       presence.Of[float64]        `json:"score"`
		Level       presence.Of[int32]          `json:"level"`
	}

	type User struct {
		ID          presence.Of[uuid.UUID]       `json:"id"`
		Username    presence.Of[string]          `json:"username"`
		Email       presence.Of[string]          `json:"email"`
		FirstName   presence.Of[string]          `json:"firstName"`
		LastName    presence.Of[string]          `json:"lastName"`
		Age         presence.Of[int]             `json:"age"`
		IsActive    presence.Of[bool]            `json:"isActive"`
		Balance     presence.Of[float64]         `json:"balance"`
		Profile     presence.Of[Profile]         `json:"profile"`
		Roles       presence.Of[[]string]        `json:"roles"`
		Permissions presence.Of[map[string]bool] `json:"permissions"`
		CreatedAt   presence.Of[int64]           `json:"createdAt"`
	}

	t.Run("deeply nested structure with all values", func(t *testing.T) {
		// Create deeply nested structure
		address := Address{
			Street:     presence.FromValue("123 Main St"),
			City:       presence.FromValue("New York"),
			PostalCode: presence.FromValue("10001"),
			Country:    presence.FromValue("USA"),
			Verified:   presence.FromValue(true),
			Lat:        presence.FromValue(40.7128),
			Lng:        presence.FromValue(-74.0060),
		}

		contact := ContactInfo{
			Email:       presence.FromValue("user@example.com"),
			Phone:       presence.FromValue("+1-555-0100"),
			Address:     presence.FromValue[Address](address),
			IsPrimary:   presence.FromValue(true),
			LastUpdated: presence.FromValue(int64(1234567890)),
		}

		metadata := Metadata{
			Tags:        presence.FromValue([]string{"premium", "verified"}),
			Properties:  presence.FromValue(map[string]string{"theme": "dark", "language": "en"}),
			Version:     presence.FromValue(3),
			IsActive:    presence.FromValue(true),
			CreatedBy:   presence.FromValue("admin"),
			CreatedByID: presence.FromValue(uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")),
		}

		profile := Profile{
			Bio:         presence.FromValue("Software developer"),
			Website:     presence.FromValue("https://example.com"),
			AvatarURL:   presence.FromValue("https://example.com/avatar.jpg"),
			Contacts:    presence.FromValue([]ContactInfo{contact}),
			Preferences: presence.FromValue(map[string]any{"notifications": true, "theme": "dark"}),
			Metadata:    presence.FromValue[Metadata](metadata),
			Score:       presence.FromValue(98.5),
			Level:       presence.FromValue(int32(42)),
		}

		user := User{
			ID:          presence.FromValue(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")),
			Username:    presence.FromValue("johndoe"),
			Email:       presence.FromValue("john@example.com"),
			FirstName:   presence.FromValue("John"),
			LastName:    presence.FromValue("Doe"),
			Age:         presence.FromValue(30),
			IsActive:    presence.FromValue(true),
			Balance:     presence.FromValue(1234.56),
			Profile:     presence.FromValue(profile),
			Roles:       presence.FromValue([]string{"admin", "user"}),
			Permissions: presence.FromValue(map[string]bool{"read": true, "write": true, "delete": false}),
			CreatedAt:   presence.FromValue(int64(1609459200)),
		}

		// Marshal
		data, err := json.Marshal(user)
		require.NoError(t, err)
		require.NotEmpty(t, data)

		// Unmarshal
		var restored User
		err = json.Unmarshal(data, &restored)
		require.NoError(t, err)

		// Verify top-level fields
		assert.Equal(t, *user.ID.GetValue(), *restored.ID.GetValue())
		assert.Equal(t, *user.Username.GetValue(), *restored.Username.GetValue())
		assert.Equal(t, *user.Email.GetValue(), *restored.Email.GetValue())
		assert.Equal(t, *user.Age.GetValue(), *restored.Age.GetValue())
		assert.Equal(t, *user.Balance.GetValue(), *restored.Balance.GetValue())

		// Verify roles (array)
		roles := *restored.Roles.GetValue()
		assert.Len(t, roles, 2)
		assert.Equal(t, "admin", roles[0])
		assert.Equal(t, "user", roles[1])

		// Verify permissions (map)
		permissions := *restored.Permissions.GetValue()
		assert.Equal(t, true, permissions["read"])
		assert.Equal(t, true, permissions["write"])
		assert.Equal(t, false, permissions["delete"])

		// Verify nested profile
		profileData := restored.Profile.GetValue()
		assert.Equal(t, "Software developer", *profileData.Bio.GetValue())
		assert.Equal(t, "https://example.com", *profileData.Website.GetValue())
		assert.Equal(t, 98.5, *profileData.Score.GetValue())
		assert.Equal(t, int32(42), *profileData.Level.GetValue())

		// Verify deeply nested metadata
		metadataData := *profileData.Metadata.GetValue()
		assert.Equal(t, 3, *metadataData.Version.GetValue())
		assert.Equal(t, true, *metadataData.IsActive.GetValue())
		assert.Equal(t, "admin", *metadataData.CreatedBy.GetValue())
		assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", metadataData.CreatedByID.GetValue().String())

		// Verify deeply nested tags
		tags := *metadataData.Tags.GetValue()
		assert.Len(t, tags, 2)
		assert.Equal(t, "premium", tags[0])
		assert.Equal(t, "verified", tags[1])
	})

	t.Run("deeply nested structure with mixed null values", func(t *testing.T) {
		// Create structure with some null values at different levels
		contact := ContactInfo{
			Email:       presence.FromValue("contact@example.com"),
			Phone:       presence.Null[string](),  // Null phone
			Address:     presence.Null[Address](), // Null address
			IsPrimary:   presence.FromValue(false),
			LastUpdated: presence.FromValue(int64(9876543210)),
		}

		metadata := Metadata{
			Tags:        presence.FromValue([]string{"new"}),
			Properties:  presence.Null[map[string]string](), // Null properties
			Version:     presence.FromValue(1),
			IsActive:    presence.Null[bool](), // Null isActive
			CreatedBy:   presence.FromValue("system"),
			CreatedByID: presence.Null[uuid.UUID](), // Null UUID
		}

		profile := Profile{
			Bio:         presence.Null[string](), // Null bio
			Website:     presence.FromValue("https://site.com"),
			AvatarURL:   presence.Null[string](), // Null avatar
			Contacts:    presence.FromValue([]ContactInfo{contact}),
			Preferences: presence.Null[map[string]any](), // Null preferences
			Metadata:    presence.FromValue(metadata),
			Score:       presence.FromValue(75.0),
			Level:       presence.Null[int32](), // Null level
		}

		user := User{
			ID:          presence.FromValue(uuid.MustParse("abcd1234-e89b-12d3-a456-426614174000")),
			Username:    presence.FromValue("janedoe"),
			Email:       presence.Null[string](), // Null email
			FirstName:   presence.FromValue("Jane"),
			LastName:    presence.Null[string](), // Null last name
			Age:         presence.FromValue(25),
			IsActive:    presence.FromValue(true),
			Balance:     presence.Null[float64](), // Null balance
			Profile:     presence.FromValue(profile),
			Roles:       presence.Null[[]string](), // Null roles
			Permissions: presence.FromValue(map[string]bool{"read": true}),
			CreatedAt:   presence.FromValue(int64(1609459200)),
		}

		// Marshal
		data, err := json.Marshal(user)
		require.NoError(t, err)
		require.NotEmpty(t, data)

		// Unmarshal
		var restored User
		err = json.Unmarshal(data, &restored)
		require.NoError(t, err)

		// Verify null fields
		assert.True(t, restored.Email.IsNull())
		assert.True(t, restored.LastName.IsNull())
		assert.True(t, restored.Balance.IsNull())
		assert.True(t, restored.Roles.IsNull())

		// Verify non-null fields
		assert.Equal(t, *user.ID.GetValue(), *restored.ID.GetValue())
		assert.Equal(t, *user.Username.GetValue(), *restored.Username.GetValue())
		assert.Equal(t, *user.Age.GetValue(), *restored.Age.GetValue())

		// Verify nested profile with nulls
		profileData := *restored.Profile.GetValue()
		assert.True(t, profileData.Bio.IsNull())
		assert.Equal(t, "https://site.com", *profileData.Website.GetValue())
		assert.True(t, profileData.AvatarURL.IsNull())
		assert.True(t, profileData.Preferences.IsNull())
		assert.Equal(t, 75.0, *profileData.Score.GetValue())
		assert.True(t, profileData.Level.IsNull())

		// Verify deeply nested metadata with nulls
		metadataData := *profileData.Metadata.GetValue()
		assert.Equal(t, 1, *metadataData.Version.GetValue())
		assert.True(t, metadataData.IsActive.IsNull())
		assert.True(t, metadataData.Properties.IsNull())
		assert.True(t, metadataData.CreatedByID.IsNull())
	})

	t.Run("deeply nested structure with all null values", func(t *testing.T) {
		user := User{
			ID:          presence.Null[uuid.UUID](),
			Username:    presence.Null[string](),
			Email:       presence.Null[string](),
			FirstName:   presence.Null[string](),
			LastName:    presence.Null[string](),
			Age:         presence.Null[int](),
			IsActive:    presence.Null[bool](),
			Balance:     presence.Null[float64](),
			Profile:     presence.Null[Profile](),
			Roles:       presence.Null[[]string](),
			Permissions: presence.Null[map[string]bool](),
			CreatedAt:   presence.Null[int64](),
		}

		// Marshal
		data, err := json.Marshal(user)
		require.NoError(t, err)
		require.NotEmpty(t, data)

		// Verify all fields are null in JSON
		var jsonMap map[string]any
		err = json.Unmarshal(data, &jsonMap)
		require.NoError(t, err)
		for key, value := range jsonMap {
			assert.Nil(t, value, "Field %s should be null", key)
		}

		// Unmarshal
		var restored User
		err = json.Unmarshal(data, &restored)
		require.NoError(t, err)

		// Verify all fields are null
		assert.True(t, restored.ID.IsNull())
		assert.True(t, restored.Username.IsNull())
		assert.True(t, restored.Email.IsNull())
		assert.True(t, restored.FirstName.IsNull())
		assert.True(t, restored.LastName.IsNull())
		assert.True(t, restored.Age.IsNull())
		assert.True(t, restored.IsActive.IsNull())
		assert.True(t, restored.Balance.IsNull())
		assert.True(t, restored.Profile.IsNull())
		assert.True(t, restored.Roles.IsNull())
		assert.True(t, restored.Permissions.IsNull())
		assert.True(t, restored.CreatedAt.IsNull())
	})

	t.Run("array of complex structures", func(t *testing.T) {
		contact1 := ContactInfo{
			Email:       presence.FromValue("contact1@example.com"),
			Phone:       presence.FromValue("+1-555-0101"),
			IsPrimary:   presence.FromValue(true),
			LastUpdated: presence.FromValue(int64(1000000)),
		}

		contact2 := ContactInfo{
			Email:       presence.FromValue("contact2@example.com"),
			Phone:       presence.Null[string](),
			IsPrimary:   presence.FromValue(false),
			LastUpdated: presence.Null[int64](),
		}

		contact3 := ContactInfo{
			Email:       presence.Null[string](),
			Phone:       presence.Null[string](),
			IsPrimary:   presence.Null[bool](),
			LastUpdated: presence.Null[int64](),
		}

		contacts := []ContactInfo{contact1, contact2, contact3}

		// Marshal
		data, err := json.Marshal(contacts)
		require.NoError(t, err)
		require.NotEmpty(t, data)

		// Unmarshal
		var restored []ContactInfo
		err = json.Unmarshal(data, &restored)
		require.NoError(t, err)

		assert.Len(t, restored, 3)

		// Verify first contact (all fields present)
		assert.Equal(t, "contact1@example.com", *restored[0].Email.GetValue())
		assert.Equal(t, "+1-555-0101", *restored[0].Phone.GetValue())
		assert.Equal(t, true, *restored[0].IsPrimary.GetValue())

		// Verify second contact (some nulls)
		assert.Equal(t, "contact2@example.com", *restored[1].Email.GetValue())
		assert.True(t, restored[1].Phone.IsNull())
		assert.Equal(t, false, *restored[1].IsPrimary.GetValue())
		assert.True(t, restored[1].LastUpdated.IsNull())

		// Verify third contact (all nulls)
		assert.True(t, restored[2].Email.IsNull())
		assert.True(t, restored[2].Phone.IsNull())
		assert.True(t, restored[2].IsPrimary.IsNull())
		assert.True(t, restored[2].LastUpdated.IsNull())
	})

	t.Run("map with complex presence values", func(t *testing.T) {
		data := map[string]presence.Of[any]{
			"user1": presence.FromValue[any](map[string]any{
				"name":   "Alice",
				"age":    30,
				"active": true,
			}),
			"user2": presence.FromValue[any](map[string]any{
				"name": "Bob",
				"age":  25,
			}),
			"user3": presence.Null[any](),
		}

		// Marshal
		jsonData, err := json.Marshal(data)
		require.NoError(t, err)
		require.NotEmpty(t, jsonData)

		// Unmarshal
		var restored map[string]presence.Of[any]
		err = json.Unmarshal(jsonData, &restored)
		require.NoError(t, err)

		assert.Len(t, restored, 3)

		// Verify user1
		user1Val := restored["user1"]
		user1 := (*user1Val.GetValue()).(map[string]any)
		assert.Equal(t, "Alice", user1["name"])
		assert.Equal(t, float64(30), user1["age"])
		assert.Equal(t, true, user1["active"])

		// Verify user2
		user2Val := restored["user2"]
		user2 := (*user2Val.GetValue()).(map[string]any)
		assert.Equal(t, "Bob", user2["name"])
		assert.Equal(t, float64(25), user2["age"])

		// Verify user3 is null
		user3Val := restored["user3"]
		assert.True(t, user3Val.IsNull())
	})

	t.Run("extreme nesting with 5 levels", func(t *testing.T) {
		// Level 5 (deepest)
		level5 := map[string]any{
			"value":    "deep value",
			"level":    5,
			"isDeep":   true,
			"metadata": []string{"tag1", "tag2", "tag3"},
		}

		// Level 4
		level4 := map[string]any{
			"data":    level5,
			"level":   4,
			"count":   42,
			"numbers": []int{1, 2, 3, 4, 5},
		}

		// Level 3
		level3 := map[string]any{
			"nested": level4,
			"level":  3,
			"active": true,
			"items":  []map[string]any{{"id": 1}, {"id": 2}},
		}

		// Level 2
		level2 := map[string]any{
			"inner":       level3,
			"level":       2,
			"description": "second level",
			"tags":        []string{"a", "b", "c"},
		}

		// Level 1 (top)
		level1 := presence.FromValue[any](map[string]any{
			"root":  level2,
			"level": 1,
			"name":  "top level",
		})

		// Marshal
		data, err := json.Marshal(level1)
		require.NoError(t, err)
		require.NotEmpty(t, data)

		// Unmarshal
		var restored presence.Of[any]
		err = json.Unmarshal(data, &restored)
		require.NoError(t, err)

		// Navigate through all levels
		l1 := (*restored.GetValue()).(map[string]any)
		assert.Equal(t, float64(1), l1["level"])
		assert.Equal(t, "top level", l1["name"])

		l2 := l1["root"].(map[string]any)
		assert.Equal(t, float64(2), l2["level"])
		assert.Equal(t, "second level", l2["description"])

		l3 := l2["inner"].(map[string]any)
		assert.Equal(t, float64(3), l3["level"])
		assert.Equal(t, true, l3["active"])

		l4 := l3["nested"].(map[string]any)
		assert.Equal(t, float64(4), l4["level"])
		assert.Equal(t, float64(42), l4["count"])

		l5 := l4["data"].(map[string]any)
		assert.Equal(t, float64(5), l5["level"])
		assert.Equal(t, "deep value", l5["value"])
		assert.Equal(t, true, l5["isDeep"])

		// Verify deeply nested array
		metadata := l5["metadata"].([]any)
		assert.Len(t, metadata, 3)
		assert.Equal(t, "tag1", metadata[0])
		assert.Equal(t, "tag2", metadata[1])
		assert.Equal(t, "tag3", metadata[2])
	})
}

func TestMarshalJSON_ThreeState(t *testing.T) {
	t.Run("unset marshals as null", func(t *testing.T) {
		n := presence.Of[string]{}
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("null"), data, "unset should marshal as null")
	})

	t.Run("unset with UnsetSkip has IsZero true", func(t *testing.T) {
		n := presence.Of[string]{}
		n.SetMarshalUnset(presence.UnsetSkip)
		assert.True(t, n.IsZero(), "unset with UnsetSkip should be zero for omitempty")
	})

	t.Run("unset with UnsetNull has IsZero false", func(t *testing.T) {
		n := presence.Of[string]{}
		n.SetMarshalUnset(presence.UnsetNull)
		assert.False(t, n.IsZero(), "unset with UnsetNull should not be zero")
	})

	t.Run("explicit null always returns null", func(t *testing.T) {
		n := presence.Null[string]()
		n.SetMarshalUnset(presence.UnsetSkip) // should not affect null
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("null"), data)
	})

	t.Run("value returns value regardless of config", func(t *testing.T) {
		n := presence.FromValue("test")
		n.SetMarshalUnset(presence.UnsetSkip)
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte(`"test"`), data)
	})
}

func TestMarshalJSON_OmitZero(t *testing.T) {
	// Note: omitzero is a Go 1.24+ feature that uses IsZero() to determine if a field should be omitted
	type TestStruct struct {
		Name presence.Of[string] `json:"name,omitzero"`
		Age  presence.Of[int]    `json:"age,omitzero"`
	}

	t.Run("unset fields omitted with omitzero", func(t *testing.T) {
		s := TestStruct{
			Name: presence.FromValue("John"),
			// Age left as unset - should be omitted because IsZero() returns true
		}
		data, err := json.Marshal(s)
		require.NoError(t, err)
		assert.JSONEq(t, `{"name":"John"}`, string(data))
	})

	t.Run("null fields included with omitzero", func(t *testing.T) {
		s := TestStruct{
			Name: presence.FromValue("John"),
			Age:  presence.Null[int](), // explicitly null - IsZero() returns false
		}
		data, err := json.Marshal(s)
		require.NoError(t, err)
		assert.JSONEq(t, `{"name":"John","age":null}`, string(data))
	})

	t.Run("value fields included with omitzero", func(t *testing.T) {
		s := TestStruct{
			Name: presence.FromValue("John"),
			Age:  presence.FromValue(30),
		}
		data, err := json.Marshal(s)
		require.NoError(t, err)
		assert.JSONEq(t, `{"name":"John","age":30}`, string(data))
	})
}

func TestMarshalJSON_OmitEmpty(t *testing.T) {
	// Note: omitempty does NOT use IsZero() - it uses its own rules for "empty"
	// For custom types with MarshalJSON, omitempty checks if the marshaled value is
	// "null", "false", 0, "", or empty array/map
	type TestStruct struct {
		Name presence.Of[string] `json:"name,omitempty"`
		Age  presence.Of[int]    `json:"age,omitempty"`
	}

	t.Run("unset fields marshal as null with omitempty", func(t *testing.T) {
		s := TestStruct{
			Name: presence.FromValue("John"),
			// Age left as unset - marshals as null, included in output
		}
		data, err := json.Marshal(s)
		require.NoError(t, err)
		// Note: omitempty doesn't omit null for custom MarshalJSON types
		assert.JSONEq(t, `{"name":"John","age":null}`, string(data))
	})

	t.Run("null fields included with omitempty", func(t *testing.T) {
		s := TestStruct{
			Name: presence.FromValue("John"),
			Age:  presence.Null[int](),
		}
		data, err := json.Marshal(s)
		require.NoError(t, err)
		assert.JSONEq(t, `{"name":"John","age":null}`, string(data))
	})
}

func TestUnmarshalJSON_ThreeState(t *testing.T) {
	t.Run("explicit null becomes null state", func(t *testing.T) {
		var n presence.Of[string]
		err := n.UnmarshalJSON([]byte("null"))
		require.NoError(t, err)
		assert.True(t, n.IsNull())
		assert.False(t, n.IsUnset())
		assert.True(t, n.IsSet())
	})

	t.Run("missing field stays unset", func(t *testing.T) {
		type TestStruct struct {
			Name presence.Of[string] `json:"name"`
			Age  presence.Of[int]    `json:"age"`
		}
		var s TestStruct
		err := json.Unmarshal([]byte(`{"name":"John"}`), &s)
		require.NoError(t, err)

		assert.False(t, s.Name.IsUnset(), "name should be set")
		assert.Equal(t, "John", *s.Name.GetValue())

		assert.True(t, s.Age.IsUnset(), "age should be unset")
		assert.False(t, s.Age.IsNull(), "age should not be null")
	})

	t.Run("explicit null in JSON becomes null", func(t *testing.T) {
		type TestStruct struct {
			Name presence.Of[string] `json:"name"`
			Age  presence.Of[int]    `json:"age"`
		}
		var s TestStruct
		err := json.Unmarshal([]byte(`{"name":"John","age":null}`), &s)
		require.NoError(t, err)

		assert.False(t, s.Age.IsUnset(), "age should not be unset")
		assert.True(t, s.Age.IsNull(), "age should be null")
	})

	t.Run("value in JSON becomes value", func(t *testing.T) {
		var n presence.Of[int]
		err := n.UnmarshalJSON([]byte("42"))
		require.NoError(t, err)
		assert.False(t, n.IsUnset())
		assert.False(t, n.IsNull())
		assert.True(t, n.IsSet())
		assert.Equal(t, 42, *n.GetValue())
	})
}

func TestMarshalUnmarshal(t *testing.T) {
	obj := getTestObjs(getEmbeddedObj())
	toObj := []testedStruct[embeddedStruct]{{}, {}}

	b, err := json.Marshal(obj)
	t.Run("Marshal nested structs test", func(t *testing.T) {
		require.NoError(t, err, "Marshaling Presence data failed")
	})

	t.Run("Unmarshal tests suite", func(t *testing.T) {
		err = json.Unmarshal(b, &toObj)
		require.NoError(t, err, "Unmarshaling into Presence data failed")

		for i := range 2 {
			t.Run("Simple string matching", func(t *testing.T) {
				assert.Equal(t, obj[i].Name.GetValue(), toObj[i].Name.GetValue(), "Name mismatch at index %d", i)
			})

			t.Run("Simple datetime matching", func(t *testing.T) {
				dte := toObj[i].DateTo.GetValue()
				if dte == nil {
					assert.Nil(t, obj[i].DateTo.GetValue(), "DateTo nil value mismatch at index %d", i)
				} else {
					assert.Equal(t, time.Duration(0), dte.Sub(now), "DateTo value mismatch at index %d", i)
				}
			})

			data := obj[i].Data.GetValue()
			if data == nil {
				t.Run("nil data into nil object checking", func(t *testing.T) {
					assert.Nil(t, toObj[i].Data.GetValue(), "Data nil value mismatch at index %d", i)
				})
			} else {
				t.Run("non nil data into non nil object matching", func(t *testing.T) {
					assert.Equal(t, data.Bool.GetValue(), toObj[i].Data.GetValue().Bool.GetValue(), "Data.Bool mismatch")
					assert.Equal(t, data.Int, toObj[i].Data.GetValue().Int, "Data.Int mismatch")
					assert.Equal(t, data.String, toObj[i].Data.GetValue().String, "Data.String mismatch")
				})
			}
		}
	})
}
