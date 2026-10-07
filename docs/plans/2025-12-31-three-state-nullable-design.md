# Three-State presence Design

## Overview

This design adds a 3-state model to the presence library, distinguishing between:
- **Unset**: Zero value, field was never touched
- **Null**: Explicitly set to null
- **Value**: Has a concrete value

This enables PATCH API semantics, partial database updates, and form handling where distinguishing "not provided" from "explicitly null" is important.

## Core State Model

**Struct Definition:**

```go
type Of[T ...] struct {
    val          *T
    isSet        bool
    marshalUnset *MarshalUnsetBehavior
    scanNull     *ScanNullBehavior
}
```

**State Definitions:**
| State | isSet | val |
|-------|-------|-----|
| Unset | false | nil |
| Null | true | nil |
| Value | true | &v |

## Constructors

| Function | Result |
|----------|--------|
| `Of[T]{}` | Unset |
| `Null[T]()` | Null |
| `FromValue[T](v)` | Value |

**Note:** `Null[T]()` behavior changes from returning zero value to returning explicit null (`isSet: true`). This is a breaking change for v2.

## State Query Methods

| Method | Unset | Null | Value |
|--------|-------|------|-------|
| `IsUnset()` | true | false | false |
| `IsNull()` | false | true | false |
| `IsSet()` | false | true | true |
| `GetValue()` | nil | nil | *T |

## State Mutation Methods

| Method | Result |
|--------|--------|
| `SetValue(T)` | Value state |
| `SetValueP(*T)` | Value if non-nil, Null if nil |
| `SetNull()` | Null state |
| `Unset()` | Unset state |

## JSON Marshaling

**Behavior with configurable unset handling:**

| State | Default (`UnsetSkip`) | Alternative (`UnsetNull`) |
|-------|----------------------|---------------------------|
| Unset | Field omitted | `null` |
| Null | `null` | `null` |
| Value | The value | The value |

**Configuration:**

```go
type MarshalUnsetBehavior int

const (
    UnsetSkip MarshalUnsetBehavior = iota  // default: omit field
    UnsetNull                               // marshal as null
)

// Package-level default
presence.SetDefaultMarshalUnset(presence.UnsetSkip)  // default
presence.SetDefaultMarshalUnset(presence.UnsetNull)  // alternative

// Per-value override
val := presence.Of[string]{}
val.SetMarshalUnset(presence.UnsetNull)
```

**Implementation:**

```go
func (n Of[T]) MarshalJSON() ([]byte, error) {
    if n.IsUnset() {
        behavior := n.marshalUnset
        if behavior == nil {
            behavior = &defaultMarshalUnset
        }
        if *behavior == UnsetSkip {
            return nil, nil  // signals omitempty to skip
        }
        return []byte("null"), nil
    }
    if n.IsNull() {
        return []byte("null"), nil
    }
    return json.Marshal(n.val)
}
```

**Note:** Users need `omitempty` struct tag for field omission to work.

## JSON Unmarshaling

| JSON Input | Result State |
|------------|--------------|
| Field missing | Unset (`isSet: false`) |
| Field is `null` | Null (`isSet: true, val: nil`) |
| Field has value | Value (`isSet: true, val: &v`) |

**Implementation:**

```go
func (n *Of[T]) UnmarshalJSON(data []byte) error {
    if n == nil {
        n = new(Of[T])
    }

    // Explicit null in JSON
    if data == nil || string(data) == "null" {
        n.SetNull()
        return nil
    }

    // Has a value
    if n.val == nil {
        n.val = new(T)
    }

    err := json.Unmarshal(data, n.val)
    if err != nil {
        return fmt.Errorf("presence Unmarshal Error: %w", err)
    }

    n.isSet = true
    return nil
}
```

Missing fields never call `UnmarshalJSON` - they remain at zero value (unset).

**Usage example:**

```go
type Request struct {
    Name presence.Of[string] `json:"name"`
    Age  presence.Of[int]    `json:"age"`
}

// Input: {"name": "John"}
// Result: Name = Value("John"), Age = Unset

// Input: {"name": "John", "age": null}
// Result: Name = Value("John"), Age = Null

// Input: {"name": null, "age": 30}
// Result: Name = Null, Age = Value(30)
```

## Database Operations

### Scan (reading from database)

| DB Value | Default Result | With `ScanNullAsUnset` option |
|----------|----------------|-------------------------------|
| SQL NULL | Null (`isSet: true, val: nil`) | Unset (`isSet: false`) |
| Value | Value (`isSet: true, val: &v`) | Value (`isSet: true, val: &v`) |

**Configuration:**

```go
type ScanNullBehavior int

const (
    ScanNullAsNull  ScanNullBehavior = iota  // default
    ScanNullAsUnset
)

// Package-level default
presence.SetDefaultScanNull(presence.ScanNullAsNull)   // default
presence.SetDefaultScanNull(presence.ScanNullAsUnset)  // alternative

// Per-value override
val := presence.Of[string]{}
val.SetScanNull(presence.ScanNullAsUnset)
```

### Value (writing to database)

| State | DB Value |
|-------|----------|
| Unset | SQL NULL |
| Null | SQL NULL |
| Value | The value |

No configuration needed - both unset and null write as SQL NULL. Users who need to skip unset fields should check `IsUnset()` before building their query.

## Configuration Methods

**Per-value:**

| Method | Purpose |
|--------|---------|
| `SetMarshalUnset(MarshalUnsetBehavior)` | Per-value marshal behavior |
| `SetScanNull(ScanNullBehavior)` | Per-value scan behavior |

**Package-level defaults:**

| Function | Purpose |
|----------|---------|
| `SetDefaultMarshalUnset(MarshalUnsetBehavior)` | Default: `UnsetSkip` |
| `SetDefaultScanNull(ScanNullBehavior)` | Default: `ScanNullAsNull` |

## Breaking Changes from v1

1. `Null[T]()` now returns explicit null (`isSet: true`) instead of zero value
2. `IsNull()` returns `false` for unset values (previously returned `true`)
3. Zero value `Of[T]{}` is now "unset" state, not "null" state

## Migration Guide

| v1 Pattern | v2 Equivalent |
|------------|---------------|
| `presence.Null[T]()` for "no value" | `presence.Of[T]{}` for unset, `presence.Null[T]()` for explicit null |
| `IsNull()` to check "no value" | `!IsSet()` or `IsUnset()` to check unset, `IsNull()` for explicit null |
