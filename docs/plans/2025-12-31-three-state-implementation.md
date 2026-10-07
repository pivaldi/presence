# Three-State presence Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add 3-state support (unset/null/value) to the presence library for PATCH APIs and partial updates.

**Architecture:** Extend `Of[T]` struct with `isSet bool` flag and configuration fields for marshal/scan behavior. Zero value is unset, `Null[T]()` returns explicit null, `FromValue[T]()` returns set value.

**Tech Stack:** Go 1.24, encoding/json, database/sql

---

## Task 1: Update presenceI Interface

**Files:**
- Modify: `presence.go:17-36`

**Step 1: Write the failing test**

```go
// tests/presence_test.go - add to existing file
func TestThreeStateInterface(t *testing.T) {
	t.Run("interface has IsUnset method", func(t *testing.T) {
		var n presence.presenceI[string] = &presence.Of[string]{}
		assert.True(t, n.IsUnset())
	})

	t.Run("interface has IsSet method", func(t *testing.T) {
		var n presence.presenceI[string] = &presence.Of[string]{}
		assert.False(t, n.IsSet())
	})

	t.Run("interface has Unset method", func(t *testing.T) {
		val := presence.FromValue("test")
		var n presence.presenceI[string] = &val
		n.Unset()
		assert.True(t, n.IsUnset())
	})
}
```

**Step 2: Run test to verify it fails**

Run: `cd tests && go test -run TestThreeStateInterface -v`
Expected: FAIL - interface missing methods

**Step 3: Update interface in presence.go**

Add these methods to `presenceI` interface:

```go
// IsUnset returns true if the value has not been set
IsUnset() bool
// IsSet returns true if the value has been set (null or value)
IsSet() bool
// Unset resets to unset state
Unset()
```

**Step 4: Run test to verify it passes**

Run: `cd tests && go test -run TestThreeStateInterface -v`
Expected: PASS

**Step 5: Commit**

Message: "Add IsUnset, IsSet, Unset methods to presenceI interface"

---

## Task 2: Fix Null Constructor

**Files:**
- Modify: `presence.go:46-49`
- Test: `tests/presence_test.go`

**Step 1: Write the failing test**

```go
// tests/presence_test.go - add new test
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
```

**Step 2: Run test to verify it fails**

Run: `cd tests && go test -run TestNullConstructor -v`
Expected: FAIL - `Null()` returns unset, not null

**Step 3: Update Null constructor**

```go
// Null is a presence constructor with explicit Null value.
func Null[T bool | int | int16 | int32 | int64 | string | uuid.UUID | float64 | JSON]() Of[T] {
	n := Of[T]{}
	n.SetNull()
	return n
}
```

**Step 4: Run test to verify it passes**

Run: `cd tests && go test -run TestNullConstructor -v`
Expected: PASS

**Step 5: Commit**

Message: "Change Null constructor to return explicit null state"

---

## Task 3: Rename UnSet to Unset

**Files:**
- Modify: `of.go:79-87`
- Test: `tests/presence_test.go`

**Step 1: Write the failing test**

```go
// tests/presence_test.go - add new test
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
```

**Step 2: Run test to verify it fails**

Run: `cd tests && go test -run TestUnsetMethod -v`
Expected: FAIL - method `Unset` not found (currently named `UnSet`)

**Step 3: Rename method**

```go
// Unset resets to unset state.
func (n *Of[T]) Unset() {
	if n == nil {
		n = new(Of[T])
	}

	n.isSet = false
	n.val = nil
}
```

**Step 4: Run test to verify it passes**

Run: `cd tests && go test -run TestUnsetMethod -v`
Expected: PASS

**Step 5: Commit**

Message: "Rename UnSet to Unset for Go naming conventions"

---

## Task 4: Add Configuration Types

**Files:**
- Create: `config.go`
- Test: `tests/config_test.go`

**Step 1: Write the failing test**

```go
// tests/config_test.go - new file
package tests

import (
	"testing"

	"github.com/pivaldi/presence"
	"github.com/stretchr/testify/assert"
)

func TestMarshalUnsetBehavior(t *testing.T) {
	t.Run("UnsetSkip is default", func(t *testing.T) {
		assert.Equal(t, presence.MarshalUnsetBehavior(0), presence.UnsetSkip)
	})

	t.Run("UnsetNull is alternative", func(t *testing.T) {
		assert.Equal(t, presence.MarshalUnsetBehavior(1), presence.UnsetNull)
	})
}

func TestScanNullBehavior(t *testing.T) {
	t.Run("ScanNullAsNull is default", func(t *testing.T) {
		assert.Equal(t, presence.ScanNullBehavior(0), presence.ScanNullAsNull)
	})

	t.Run("ScanNullAsUnset is alternative", func(t *testing.T) {
		assert.Equal(t, presence.ScanNullBehavior(1), presence.ScanNullAsUnset)
	})
}

func TestDefaultConfiguration(t *testing.T) {
	t.Run("default marshal unset is skip", func(t *testing.T) {
		assert.Equal(t, presence.UnsetSkip, presence.GetDefaultMarshalUnset())
	})

	t.Run("default scan null is null", func(t *testing.T) {
		assert.Equal(t, presence.ScanNullAsNull, presence.GetDefaultScanNull())
	})
}
```

**Step 2: Run test to verify it fails**

Run: `cd tests && go test -run 'TestMarshalUnsetBehavior|TestScanNullBehavior|TestDefaultConfiguration' -v`
Expected: FAIL - types and constants not defined

**Step 3: Create config.go**

```go
package presence

import "sync"

// MarshalUnsetBehavior controls how unset values are marshaled to JSON.
type MarshalUnsetBehavior int

const (
	// UnsetSkip omits unset fields from JSON output (requires omitempty tag).
	UnsetSkip MarshalUnsetBehavior = iota
	// UnsetNull marshals unset fields as null.
	UnsetNull
)

// ScanNullBehavior controls how SQL NULL values are scanned.
type ScanNullBehavior int

const (
	// ScanNullAsNull interprets SQL NULL as explicit null (isSet=true, val=nil).
	ScanNullAsNull ScanNullBehavior = iota
	// ScanNullAsUnset interprets SQL NULL as unset (isSet=false, val=nil).
	ScanNullAsUnset
)

var (
	defaultMarshalUnset MarshalUnsetBehavior = UnsetSkip
	defaultScanNull     ScanNullBehavior     = ScanNullAsNull
	configMu            sync.RWMutex
)

// SetDefaultMarshalUnset sets the package-level default for marshal unset behavior.
func SetDefaultMarshalUnset(b MarshalUnsetBehavior) {
	configMu.Lock()
	defer configMu.Unlock()
	defaultMarshalUnset = b
}

// GetDefaultMarshalUnset returns the package-level default for marshal unset behavior.
func GetDefaultMarshalUnset() MarshalUnsetBehavior {
	configMu.RLock()
	defer configMu.RUnlock()
	return defaultMarshalUnset
}

// SetDefaultScanNull sets the package-level default for scan null behavior.
func SetDefaultScanNull(b ScanNullBehavior) {
	configMu.Lock()
	defer configMu.Unlock()
	defaultScanNull = b
}

// GetDefaultScanNull returns the package-level default for scan null behavior.
func GetDefaultScanNull() ScanNullBehavior {
	configMu.RLock()
	defer configMu.RUnlock()
	return defaultScanNull
}
```

**Step 4: Run test to verify it passes**

Run: `cd tests && go test -run 'TestMarshalUnsetBehavior|TestScanNullBehavior|TestDefaultConfiguration' -v`
Expected: PASS

**Step 5: Commit**

Message: "Add configuration types for marshal and scan behavior"

---

## Task 5: Add Per-Value Configuration Fields

**Files:**
- Modify: `of.go:12-15`
- Test: `tests/config_test.go`

**Step 1: Write the failing test**

```go
// tests/config_test.go - add to existing file
func TestPerValueConfiguration(t *testing.T) {
	t.Run("SetMarshalUnset configures per-value behavior", func(t *testing.T) {
		n := presence.Of[string]{}
		n.SetMarshalUnset(presence.UnsetNull)
		assert.Equal(t, presence.UnsetNull, n.GetMarshalUnset())
	})

	t.Run("SetScanNull configures per-value behavior", func(t *testing.T) {
		n := presence.Of[string]{}
		n.SetScanNull(presence.ScanNullAsUnset)
		assert.Equal(t, presence.ScanNullAsUnset, n.GetScanNull())
	})

	t.Run("default uses package default for marshal", func(t *testing.T) {
		n := presence.Of[string]{}
		assert.Equal(t, presence.GetDefaultMarshalUnset(), n.GetMarshalUnset())
	})

	t.Run("default uses package default for scan", func(t *testing.T) {
		n := presence.Of[string]{}
		assert.Equal(t, presence.GetDefaultScanNull(), n.GetScanNull())
	})
}
```

**Step 2: Run test to verify it fails**

Run: `cd tests && go test -run TestPerValueConfiguration -v`
Expected: FAIL - methods not defined

**Step 3: Add fields and methods to Of struct in of.go**

Update struct:
```go
type Of[T bool | int | int16 | int32 | int64 | string | uuid.UUID | float64 | JSON] struct {
	val          *T
	isSet        bool
	marshalUnset *MarshalUnsetBehavior
	scanNull     *ScanNullBehavior
}
```

Add methods after Unset method:
```go
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
```

**Step 4: Run test to verify it passes**

Run: `cd tests && go test -run TestPerValueConfiguration -v`
Expected: PASS

**Step 5: Commit**

Message: "Add per-value configuration for marshal and scan behavior"

---

## Task 6: Update MarshalJSON for 3-State

**Files:**
- Modify: `of.go:89-96` (MarshalJSON method)
- Test: `tests/marshal_test.go`

**Step 1: Write the failing test**

```go
// tests/marshal_test.go - add new test section
func TestMarshalJSON_ThreeState(t *testing.T) {
	t.Run("unset with UnsetSkip returns nil", func(t *testing.T) {
		n := presence.Of[string]{}
		n.SetMarshalUnset(presence.UnsetSkip)
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Nil(t, data, "unset with UnsetSkip should return nil for omitempty")
	})

	t.Run("unset with UnsetNull returns null", func(t *testing.T) {
		n := presence.Of[string]{}
		n.SetMarshalUnset(presence.UnsetNull)
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("null"), data)
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

func TestMarshalJSON_OmitEmpty(t *testing.T) {
	type TestStruct struct {
		Name presence.Of[string] `json:"name,omitempty"`
		Age  presence.Of[int]    `json:"age,omitempty"`
	}

	t.Run("unset fields omitted with omitempty", func(t *testing.T) {
		s := TestStruct{
			Name: presence.FromValue("John"),
			// Age left as unset
		}
		data, err := json.Marshal(s)
		require.NoError(t, err)
		assert.JSONEq(t, `{"name":"John"}`, string(data))
	})

	t.Run("null fields included even with omitempty", func(t *testing.T) {
		s := TestStruct{
			Name: presence.FromValue("John"),
			Age:  presence.Null[int](),
		}
		data, err := json.Marshal(s)
		require.NoError(t, err)
		assert.JSONEq(t, `{"name":"John","age":null}`, string(data))
	})
}
```

**Step 2: Run test to verify it fails**

Run: `cd tests && go test -run 'TestMarshalJSON_ThreeState|TestMarshalJSON_OmitEmpty' -v`
Expected: FAIL - unset values not handled correctly

**Step 3: Update MarshalJSON**

```go
// MarshalJSON implements the encoding json interface.
func (n Of[T]) MarshalJSON() ([]byte, error) {
	if n.IsUnset() {
		if n.GetMarshalUnset() == UnsetSkip {
			return nil, nil
		}
		return []byte("null"), nil
	}

	if n.IsNull() {
		return []byte("null"), nil
	}

	return marshalJSON(&n)
}
```

**Step 4: Run test to verify it passes**

Run: `cd tests && go test -run 'TestMarshalJSON_ThreeState|TestMarshalJSON_OmitEmpty' -v`
Expected: PASS

**Step 5: Commit**

Message: "Update MarshalJSON to handle 3-state with configurable unset behavior"

---

## Task 7: Update UnmarshalJSON for 3-State

**Files:**
- Modify: `of.go:98-120` (UnmarshalJSON method)
- Test: `tests/marshal_test.go`

**Step 1: Write the failing test**

```go
// tests/marshal_test.go - add new test
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
```

**Step 2: Run test to verify it fails**

Run: `cd tests && go test -run TestUnmarshalJSON_ThreeState -v`
Expected: FAIL - isSet not set on unmarshal value

**Step 3: Update UnmarshalJSON**

```go
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
		return fmt.Errorf("presence Unmarshal Error : %w", err)
	}

	n.isSet = true
	return nil
}
```

**Step 4: Run test to verify it passes**

Run: `cd tests && go test -run TestUnmarshalJSON_ThreeState -v`
Expected: PASS

**Step 5: Commit**

Message: "Update UnmarshalJSON to set isSet flag on value unmarshal"

---

## Task 8: Update Scan Methods for 3-State

**Files:**
- Modify: `presence.go` (all scan methods)
- Test: `tests/postgres_test.go`

**Step 1: Write the failing test**

```go
// tests/presence_test.go - add new test
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
```

**Step 2: Run test to verify it fails**

Run: `cd tests && go test -run TestScanNullBehavior -v`
Expected: FAIL - ScanNullAsUnset not respected

**Step 3: Create helper method and update scan methods**

Add helper to presence.go:
```go
// handleScanNull handles null scanning based on configuration.
func (n *Of[T]) handleScanNull() {
	if n.GetScanNull() == ScanNullAsUnset {
		n.Unset()
	} else {
		n.SetNull()
	}
}
```

Update each scan method to use `n.handleScanNull()` instead of `n.SetNull()`.

For example in scanString:
```go
func (n *Of[T]) scanString(v any) error {
	if n == nil {
		return errors.New("calling scanString on nil receiver")
	}

	null := sql.NullString{}
	err := null.Scan(v)
	if err != nil {
		return fmt.Errorf("presence database scanning string : %w", err)
	}

	if null.Valid {
		n.SetValue(any(null.String).(T))
	} else {
		n.handleScanNull()
	}

	return nil
}
```

Apply same pattern to: scanJSON, scanUUID, scanInt, scanFloat, scanBool, scanTime.

**Step 4: Run test to verify it passes**

Run: `cd tests && go test -run TestScanNullBehavior -v`
Expected: PASS

**Step 5: Commit**

Message: "Update scan methods to respect ScanNullBehavior configuration"

---

## Task 9: Update Existing Tests

**Files:**
- Modify: `tests/presence_test.go`
- Modify: `tests/marshal_test.go`

**Step 1: Identify failing tests**

Run: `cd tests && go test -v ./...`
Expected: Some existing tests may fail due to behavior changes

**Step 2: Update tests for new behavior**

Key changes:
- `presence.Null[T]()` now returns `IsNull()=true`, `IsUnset()=false`
- Zero value `Of[T]{}` now returns `IsNull()=false`, `IsUnset()=true`
- Tests using `Null()` expecting it to equal zero value need updating

In `presence_test.go`, update:
```go
t.Run("IsNull on zero value", func(t *testing.T) {
	var test testedStruct[embeddedStruct]
	assert.False(t, test.Name.IsNull(), "Zero value should be unset, not null")
	assert.True(t, test.Name.IsUnset(), "Zero value should be unset")
})
```

In `marshal_test.go`, tests using `Null()` that expect zero-value behavior need review. The marshal tests should mostly work since `Null()` now marshals to `"null"` which is correct.

**Step 3: Run all tests**

Run: `cd tests && go test -v ./...`
Expected: PASS

**Step 4: Commit**

Message: "Update existing tests for 3-state behavior"

---

## Task 10: Run Full Test Suite and Lint

**Files:** None (verification only)

**Step 1: Run all tests**

Run: `cd tests && go test -v ./...`
Expected: All PASS

**Step 2: Run linter**

Run: `golangci-lint run`
Expected: No errors

**Step 3: Tidy modules**

Run: `go mod tidy && cd tests && go mod tidy`
Expected: No changes or clean tidy

**Step 4: Commit** (if any fixes needed)

Message: "Fix lint issues and tidy modules"

---

## Task 11: Update README Documentation

**Files:**
- Modify: `README.md`

**Step 1: Add 3-state documentation section**

Add after "Supported Types" section:

```markdown
## Three-State Model (v2)

The library supports a 3-state model for presence values:

| State | Description | Creation |
|-------|-------------|----------|
| Unset | Field was never touched | `presence.Of[T]{}` or `var x presence.Of[T]` |
| Null | Explicitly set to null | `presence.Null[T]()` |
| Value | Has a concrete value | `presence.FromValue(x)` |

### Checking State

```go
if value.IsUnset() {
    // Field was never touched
}
if value.IsNull() {
    // Field was explicitly set to null
}
if value.IsSet() {
    // Field has null or value (not unset)
}
```

### PATCH API Example

```go
type UpdateUserRequest struct {
    Name  presence.Of[string] `json:"name,omitempty"`
    Email presence.Of[string] `json:"email,omitempty"`
    Age   presence.Of[int]    `json:"age,omitempty"`
}

func UpdateUser(req UpdateUserRequest) {
    if req.Name.IsSet() {
        if req.Name.IsNull() {
            // Clear the name
        } else {
            // Update with new name
        }
    }
    // else: don't touch name field
}
```

### Configuration

Control how unset values marshal to JSON:

```go
// Package-level default (default: UnsetSkip)
presence.SetDefaultMarshalUnset(presence.UnsetNull)

// Per-value override
val := presence.Of[string]{}
val.SetMarshalUnset(presence.UnsetNull)
```

Control how SQL NULL scans:

```go
// Package-level default (default: ScanNullAsNull)
presence.SetDefaultScanNull(presence.ScanNullAsUnset)

// Per-value override
val := presence.Of[string]{}
val.SetScanNull(presence.ScanNullAsUnset)
```
```

**Step 2: Update "Why Use This Library" section**

Add comparison showing 3-state advantage.

**Step 3: Commit**

Message: "Update README with 3-state documentation"

---

Plan complete and saved to `docs/plans/2025-12-31-three-state-implementation.md`. Two execution options:

**1. Subagent-Driven (this session)** - I dispatch fresh subagent per task, review between tasks, fast iteration

**2. Parallel Session (separate)** - Open new session with executing-plans, batch execution with checkpoints

**Which approach?**
