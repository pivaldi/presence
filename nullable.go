package nullable

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type NullableI[T any] interface {
	// IsNull returns true if itself is nil or the value is nil/null
	IsNull() bool
	// IsUnset returns true if the value has not been set
	IsUnset() bool
	// IsSet returns true if the value has been set (null or value)
	IsSet() bool
	// GetValue implements the getter.
	GetValue() *T
	// SetValue implements the setter.
	SetValue(T)
	// SetValueP implements the setter by pointer.
	SetValueP(*T)
	// SetNull set to null.
	SetNull()
	// Unset resets to unset state
	Unset()
	// MarshalJSON implements the encoding json interface.
	MarshalJSON() ([]byte, error)
	// UnmarshalJSON implements the decoding json interface.
	UnmarshalJSON([]byte) error
	// Value implements the driver.Valuer interface.
	Value() (driver.Value, error)
	// Scan implements the sql.Scanner interface.
	Scan(v any) error
}

// FromValue is a Nullable constructor from the given value thanks to Go generics' inference.
func FromValue[T any](b T) Of[T] {
	out := Of[T]{}
	out.SetValue(b)

	return out
}

// Null is a Nullable constructor with explicit Null value.
func Null[T any]() Of[T] {
	n := Of[T]{}
	n.SetNull()

	return n
}

func (n *Of[T]) scanJSON(v any) error {
	if n == nil {
		return errors.New("calling scanJSON on nil receiver")
	}

	null := sql.NullString{}
	err := null.Scan(v)
	if err != nil {
		return fmt.Errorf("nullable database scanning json : %w", err)
	}

	if null.Valid {
		value := new(T)

		if scanner, ok := any(value).(sql.Scanner); ok {
			err := scanner.Scan(v)
			if err != nil {
				return fmt.Errorf("custom scanner error on nullable : %w", err)
			}
		} else {
			err := json.Unmarshal([]byte(null.String), value)
			if err != nil {
				return fmt.Errorf("nullable database unmarshaling json : %w", err)
			}
		}

		n.SetValue(*value)
	} else {
		n.handleScanNull()
	}

	return nil
}

func (n *Of[T]) scanString(v any) error {
	if n == nil {
		return errors.New("calling scanString on nil receiver")
	}

	null := sql.NullString{}
	err := null.Scan(v)
	if err != nil {
		return fmt.Errorf("nullable database scanning string : %w", err)
	}

	if null.Valid {
		n.SetValue(any(null.String).(T))
	} else {
		n.handleScanNull()
	}

	return nil
}

func (n *Of[T]) scanUUID(v any) error {
	if n == nil {
		return errors.New("calling scanUUID on nil receiver")
	}

	null := sql.NullString{}
	err := null.Scan(v)
	if err != nil {
		return fmt.Errorf("nullable database scanning string : %w", err)
	}

	if null.Valid {
		uid, err := uuid.Parse(null.String)
		if err != nil {
			return fmt.Errorf("UUID parsing failed : %w", err)
		}

		n.SetValue(any(uid).(T))
	} else {
		n.handleScanNull()
	}

	return nil
}

func (n *Of[T]) scanInt(v any) error {
	switch any(new(T)).(type) {
	case int16, *int16:
		null := new(sql.NullInt16)
		err := null.Scan(v)
		if err != nil {
			return fmt.Errorf("nullable database scanning int16 : %w", err)
		}

		if null.Valid {
			n.SetValue(any(null.Int16).(T))
		} else {
			n.handleScanNull()
		}

		return nil
	case int32, *int32:
		null := new(sql.NullInt32)
		err := null.Scan(v)
		if err != nil {
			return fmt.Errorf("nullable database scanning int32 : %w", err)
		}

		if null.Valid {
			n.SetValue(any(null.Int32).(T))
		} else {
			n.handleScanNull()
		}

		return nil
	case int, *int:
		null := new(sql.NullInt64)
		err := null.Scan(v)
		if err != nil {
			return fmt.Errorf("nullable database scanning int : %w", err)
		}

		if null.Valid {
			n.SetValue(any(int(null.Int64)).(T))
		} else {
			n.handleScanNull()
		}

		return nil
	case int64, *int64:
		null := new(sql.NullInt64)
		err := null.Scan(v)
		if err != nil {
			return fmt.Errorf("nullable database scanning int64 : %w", err)
		}

		if null.Valid {
			n.SetValue(any(null.Int64).(T))
		} else {
			n.handleScanNull()
		}

		return nil
	}

	return fmt.Errorf("type %T is not supported", *new(T))
}

func (n *Of[T]) scanFloat(v any) error {
	null := new(sql.NullFloat64)
	err := null.Scan(v)
	if err != nil {
		return fmt.Errorf("nullable database scanning float64 : %w", err)
	}

	if null.Valid {
		n.SetValue(any(null.Float64).(T))
	} else {
		n.handleScanNull()
	}

	return nil
}

func (n *Of[T]) scanBool(v any) error {
	null := new(sql.NullBool)
	err := null.Scan(v)
	if err != nil {
		return fmt.Errorf("nullable database scanning bool : %w", err)
	}

	if null.Valid {
		n.SetValue(any(null.Bool).(T))
	} else {
		n.handleScanNull()
	}

	return nil
}

func (n *Of[T]) scanTime(v any) error {
	if v == nil {
		n.handleScanNull()

		return nil
	}

	null := new(sql.NullTime)

	switch t := v.(type) {
	case string:
		var err error
		null.Time, err = time.Parse(t, t)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	case time.Time:
		err := null.Scan(v)
		if err != nil {
			return fmt.Errorf("nullable database scanning Time : %w", err)
		}
	default:
		return fmt.Errorf("canot parse type \"%T\" with value \"%v\" to time", t, t)
	}

	if null.Valid {
		n.SetValue(any(null.Time).(T))
	} else {
		n.handleScanNull()
	}

	return nil
}

// handleScanNull handles null scanning based on configuration.
func (n *Of[T]) handleScanNull() {
	if n.GetScanNull() == ScanNullAsUnset {
		n.Unset()
	} else {
		n.SetNull()
	}
}
