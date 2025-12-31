package nullable

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Of[T bool | int | int16 | int32 | int64 | string | uuid.UUID | float64 | JSON] struct {
	val          *T
	isSet        bool
	marshalUnset *MarshalUnsetBehavior
	scanNull     *ScanNullBehavior
}

// IsNull returns true iff the value is nil and it is set
func (n *Of[T]) IsNull() bool {
	return n != nil && n.val == nil && n.isSet
}

// IsUnset returns true iff it is not set
func (n *Of[T]) IsUnset() bool {
	return n == nil || !n.isSet
}

// IsSet returns true iff it is set
func (n *Of[T]) IsSet() bool {
	return n != nil && n.isSet
}

// GetValue implements the getter.
func (n *Of[T]) GetValue() *T {
	if n == nil {
		return nil
	}

	return n.val
}

// SetValue implements the setter.
func (n *Of[T]) SetValue(b T) {
	if n == nil {
		n = new(Of[T])
		n.SetValue(b)

		return
	}

	n.isSet = true
	n.val = &b
}

// SetValueP implements the setter by pointer.
// If ref is not nil, calls SetValue(*ref)
// If ref is nil, calls SetNull()
func (n *Of[T]) SetValueP(ref *T) {
	if n == nil {
		n = new(Of[T])
	}

	if ref != nil {
		n.SetValue(*ref)
	} else {
		n.SetNull()
	}
}

// SetNull set to null.
func (n *Of[T]) SetNull() {
	if n == nil {
		n = new(Of[T])
	}

	n.isSet = true
	n.val = nil
}

// Unset resets to unset state.
func (n *Of[T]) Unset() {
	if n == nil {
		n = new(Of[T])
	}

	n.isSet = false
	n.val = nil
}

// SetMarshalUnset sets per-value marshal unset behavior.
func (n *Of[T]) SetMarshalUnset(b MarshalUnsetBehavior) {
	if n == nil {
		return
	}
	n.marshalUnset = &b
}

// GetMarshalUnset returns the effective marshal unset behavior.
func (n *Of[T]) GetMarshalUnset() MarshalUnsetBehavior {
	if n == nil || n.marshalUnset == nil {
		return GetDefaultMarshalUnset()
	}
	return *n.marshalUnset
}

// SetScanNull sets per-value scan null behavior.
func (n *Of[T]) SetScanNull(b ScanNullBehavior) {
	if n == nil {
		return
	}
	n.scanNull = &b
}

// GetScanNull returns the effective scan null behavior.
func (n *Of[T]) GetScanNull() ScanNullBehavior {
	if n == nil || n.scanNull == nil {
		return GetDefaultScanNull()
	}
	return *n.scanNull
}

// MarshalJSON implements the encoding json interface.
// Note: UnsetSkip behavior requires the struct field to have the `omitempty` tag.
// When marshaling directly (not as a struct field), unset values marshal as null.
func (n Of[T]) MarshalJSON() ([]byte, error) {
	if n.IsUnset() || n.IsNull() {
		return []byte("null"), nil
	}

	return marshalJSON(&n)
}

// IsZero implements the interface used by encoding/json's omitempty.
// Returns true for unset values when UnsetSkip is configured,
// allowing struct fields with `json:",omitempty"` to be omitted.
func (n Of[T]) IsZero() bool {
	if n.IsUnset() && n.GetMarshalUnset() == UnsetSkip {
		return true
	}
	return false
}

// UnmarshalJSON implements the decoding json interface.
func (n *Of[T]) UnmarshalJSON(data []byte) error {
	if n == nil {
		n = new(Of[T])
	}

	if data == nil || string(data) == "null" {
		n.SetNull()

		return nil
	}

	if n.val == nil {
		n.val = new(T)
	}

	err := json.Unmarshal(data, n.val)
	if err != nil {
		return fmt.Errorf("nullable Unmarshal Error : %w", err)
	}

	n.isSet = true
	return nil
}

// Value implements the driver.Valuer interface.
func (n Of[T]) Value() (driver.Value, error) {
	if n.val == nil {
		return nil, nil
	}

	switch value := any(n.val).(type) {
	case *string, *int16, *int32, *int, *int64, *float64, *bool, *time.Time, *uuid.UUID, string,
		int16, int32, int, int64, float64, bool, time.Time, uuid.UUID:
		return *n.val, nil
	case JSON:
		if value == nil {
			return nil, nil
		}

		if valuer, ok := value.(driver.Valuer); ok {
			v, err := valuer.Value()
			if err != nil {
				return nil, fmt.Errorf("custom valuer error on nullable : %w", err)
			}

			return v, nil
		}

		b, err := json.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("nullable database value error : %w", err)
		}

		return string(b), nil
	}

	return nil, fmt.Errorf("type %T is not supported for value %v", *n.val, *n.val)
}

// Scan implements the sql.Scanner interface.
// This method decodes a JSON-encoded value into the struct.
func (n *Of[T]) Scan(v any) error {
	if n == nil {
		n = new(Of[T])
	}

	// Use a zero value of T to determine the type, since n.val may be nil
	switch any(new(T)).(type) {
	case *string:
		return n.scanString(v)
	case *uuid.UUID:
		return n.scanUUID(v)
	case *int16, *int32, *int, *int64:
		return n.scanInt(v)
	case *float64:
		return n.scanFloat(v)
	case *bool:
		return n.scanBool(v)
	case *time.Time:
		return n.scanTime(v)
	case *JSON, JSON:
		return n.scanJSON(v)
	}

	return fmt.Errorf("type %T is not handled as nullable", v)
}
