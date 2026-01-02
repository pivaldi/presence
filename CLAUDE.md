# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go library (`github.com/pivaldi/presence`) that provides generic presence types for any data type, with special focus on database operations and JSON marshaling/unmarshaling. The library uses Go generics to wrap values in a presence container (`Of[T]`) that can represent SQL NULL values while maintaining type safety.

### Core Architecture

**Main library files (root directory):**
- `presence.go` - Core interface `PresenceI[T]`, helper functions (`FromValue`, `Null`, `FromPtr`, `FromBool`), functional operations (`Map`, `MapOr`, `FlatMap`, `Filter`, `Or`), and type-specific scanning methods
- `of.go` - Generic `Of[T]` struct implementation with methods for SQL scanning (`Scan`), SQL value conversion (`Value`), JSON marshaling/unmarshaling, value access (`Get`, `GetOr`, `MustGet`, `Ptr`), and state management
- `doc.go` - Package documentation

**Key design patterns:**
1. **Generic presence wrapper**: `Of[T any]` wraps any type `T` with an internal pointer `val *T` where `nil` represents NULL, plus `isSet bool` for 3-state support
2. **Type dispatch in Scan/Value**: The `Scan` and `Value` methods use type switches to route to specialized handlers for primitive types, with fallback to JSON for all other types
3. **Custom type support**: Types implementing `sql.Scanner` or `driver.Valuer` interfaces are automatically supported without JSON marshaling
4. **Dual module structure**: Main module at root, separate test module in `tests/` directory with `replace` directive
5. **3-state model**: Distinguishes between unset (zero value), null (explicitly set to null), and value (has a concrete value)
6. **Functional operations**: Package-level functions (`Map`, `FlatMap`, `Filter`, `Or`) for transforming presence values (methods can't have additional type parameters in Go)

### Supported Types

The library uses `Of[T any]` which accepts **any type**. For database operations:

- **Primitive types** (`string`, `int`, `int16`, `int32`, `int64`, `float64`, `bool`, `time.Time`, `uuid.UUID`) are stored/scanned directly
- **Custom types** implementing `sql.Scanner` and/or `driver.Valuer` use their custom serialization (see README.md example with `PhoneNumber`)
- **All other types** (structs, slices, maps, etc.) are automatically marshaled to/from JSON for database storage

## Development Commands

### Running Tests

**All tests (including PostgreSQL):**
```bash
cd tests
go test -v ./...
```

Or from the root:
```bash
make test
```

**Single test:**
```bash
cd tests
go test -run TestAllTypes -v
```

**Requirements:**
- Docker must be running locally
- testcontainers-go automatically manages PostgreSQL container
- First run downloads PostgreSQL 18 image (~80MB), subsequent runs use cached image

### Code Quality

**Lint the code:**
```bash
golangci-lint run
```

The project uses extensive linting (see `.golangci.yml`) with 30+ enabled linters including gosec, govet, errcheck, and revive.

**Tidy dependencies:**
```bash
go mod tidy
cd tests && go mod tidy
```

## Test Organization

The test suite is located in `tests/` directory with its own `go.mod` that uses a `replace` directive to reference the parent module.

**Test files:**
- `marshal_test.go` - Comprehensive tests for JSON marshaling/unmarshaling with complex nested structures
- `nullable_test.go` - Unit tests for presence value operations and edge cases
- `postgres_test.go` - Integration tests with PostgreSQL database using testcontainers
- `setup_test.go` - TestMain setup with testcontainers, database helpers, and cleanup utilities
- `config_test.go` - Tests for configuration options (marshal/scan behaviors)

**Test infrastructure:**
- Uses testcontainers-go to automatically manage PostgreSQL 18 container
- TestMain in setup_test.go starts container once for all tests
- Shared database connection stored in package-level `testDB` variable
- `cleanupTables()` helper truncates tables between tests for isolation
- Container automatically terminates after tests complete

## Working with MarshalJSON/UnmarshalJSON

The `Of[T]` type implements custom JSON marshaling:

**MarshalJSON (of.go:126-137):**
- Returns `[]byte("null")` if value is unset or null
- Otherwise calls `json.Marshal(n.GetValue())` to marshal the value directly

**UnmarshalJSON (of.go:150-172):**
- Handles `null` JSON values by calling `SetNull()`
- For non-null values, unmarshals directly into the wrapped value
- Allocates new `T` if needed before unmarshaling
- Sets `isSet = true` after successful unmarshal

**IsZero (of.go:142-147):**
- Returns `true` for unset values when `UnsetSkip` is configured
- Used by Go 1.24+ `omitzero` struct tag to omit unset fields from JSON output

**Key invariant:** JSON `null` maps to `isSet=true, val=nil`, while missing/unset is `isSet=false, val=nil`.

## Database Integration

The library integrates with `database/sql` through two interfaces:

1. **`driver.Valuer` (of.go:175-207)**: Converts Go values to database values
   - Primitive types (`string`, `int*`, `float64`, `bool`, `time.Time`, `uuid.UUID`) return their dereferenced value directly
   - Other types check for custom `driver.Valuer` first, then marshal to JSON string

2. **`sql.Scanner` (of.go:211-247)**: Converts database values to Go values
   - Routes to type-specific scan methods based on the wrapped type using type switch
   - Primitive types use optimized scanning (e.g., `scanString`, `scanInt`, `scanBool`)
   - Custom types implementing `sql.Scanner` are called directly before JSON fallback
   - All other types fall back to `scanJSON` which unmarshals from JSON
   - Each scan method (in presence.go) handles SQL NULL properly via `handleScanNull()`

## Go Version and Dependencies

- **Go version:** 1.24.10
- **Dependencies:**
  - `github.com/google/uuid` - UUID type support
  - Test dependencies: `pgx/v5`, `sqlx`, `testify`, `testcontainers-go`

## Common Gotchas

1. **Module structure**: Root module (`github.com/pivaldi/presence`) and test module (`github.com/pivaldi/presence/tests`) are separate. Always run `go mod tidy` in both directories after dependency changes.

2. **Test execution**: Integration tests require Docker to be running (testcontainers uses it). Run `cd tests && go test -v ./...` or `make test` for the full suite.

3. **Type handling in Scan/Value**: The library uses type switches in `Scan()` and `Value()` methods to handle primitive types directly. All other types fall back to JSON marshaling/unmarshaling. To add optimized handling for a new primitive type, update the type switch in both methods.

4. **Time precision**: PostgreSQL tests truncate time to seconds (`Truncate(time.Second)`) to match database precision.

5. **3-state assertions in tests**: When testing presence fields in structs, use `.IsNull()` and `.GetValue()` methods instead of checking struct fields directly (e.g., `assert.True(t, field.IsNull())` not `assert.Nil(t, field)`).
