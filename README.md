# Go Presence

[![golangci-lint](https://github.com/pivaldi/presence/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/pivaldi/presence/actions/workflows/golangci-lint.yml)
[![mod-verify](https://github.com/pivaldi/presence/actions/workflows/mod-verify.yml/badge.svg)](https://github.com/pivaldi/presence/actions/workflows/mod-verify.yml)
[![gosec](https://github.com/pivaldi/presence/actions/workflows/gosec.yaml/badge.svg)](https://github.com/pivaldi/presence/actions/workflows/gosec.yaml)
[![staticcheck](https://github.com/pivaldi/presence/actions/workflows/staticcheck.yaml/badge.svg)](https://github.com/pivaldi/presence/actions/workflows/staticcheck.yaml)
[![test](https://github.com/pivaldi/presence/actions/workflows/test.yml/badge.svg)](https://github.com/pivaldi/presence/actions/workflows/test.yml)

A type-safe presence value library for Go using generics, designed for seamless JSON marshaling and database operations.

## Features

- **Type-safe presence values** for any supported type using Go generics
- **3-state model** distinguishing unset, null, and value states for PATCH API support
- **Database-friendly** with built-in `sql.Scanner` and `driver.Valuer` implementations
- **JSON marshaling** that uses standard `null` instead of `{Valid: true, Value: ...}`
- **Configurable behavior** for marshal and scan operations (per-value and package-level)
- **PostgreSQL JSON/JSONB support** for storing complex types
- **UUID support** with `github.com/google/uuid`
- **Zero external dependencies** (except `google/uuid`)
- **Fully tested** with comprehensive unit and integration tests

## Installation

```bash
go get github.com/pivaldi/presence
```

## Quick Start

```go
import "github.com/pivaldi/presence"

// Create presence values
name := presence.FromValue("John Doe")
age := presence.FromValue(30)
email := presence.Null[string]() // Explicitly null

// Check if null
if name.IsNull() {
    // Handle null case
}

// Get value
if !age.IsNull() {
    fmt.Println(*age.GetValue()) // 30
}
```

## Supported Types

The library uses `Of[T any]` which accepts **any type**. Common usage patterns:

- **Primitives**: `int`, `int16`, `int32`, `int64`, `float64`, `bool`, `string`
- **UUID**: `uuid.UUID` (from `github.com/google/uuid`)
- **Time**: `time.Time`
- **Complex types**: structs, slices, maps - stored as JSON in database
- **Custom types**: any type implementing `sql.Scanner`/`driver.Valuer`

For database operations:
- Primitive types (`string`, `int*`, `float64`, `bool`, `time.Time`, `uuid.UUID`) are stored/scanned directly
- Custom types implementing `sql.Scanner` and/or `driver.Valuer` use their custom serialization
- All other types are automatically marshaled to/from JSON for storage

## Three-State Model

The library supports a 3-state model for presence values, enabling PATCH API semantics and partial updates:

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

// Handles all these requests correctly:
// {}                              → no updates
// {"name": "John"}                → update name only
// {"name": "John", "age": null}   → update name, clear age
// {"name": null, "age": null}     → clear both
```

### Configuration

**JSON Marshaling with omitzero (Go 1.24+):**

Unset values can be omitted from JSON output using the `omitzero` struct tag (introduced in Go 1.24). The `IsZero()` method returns `true` for unset values when `UnsetSkip` is configured:

- `UnsetSkip` (default): `IsZero()` returns `true` for unset values, allowing `omitzero` to omit them
- `UnsetNull`: `IsZero()` returns `false`, so unset values are always included as `null`

```go
type Request struct {
    Name presence.Of[string] `json:"name,omitzero"` // omitted when unset (Go 1.24+)
    Age  presence.Of[int]    `json:"age"`           // always included
}

// Package-level default (default: UnsetSkip)
presence.SetDefaultMarshalUnset(presence.UnsetNull)

// Per-value override
val := presence.Of[string]{}
val.SetMarshalUnset(presence.UnsetNull)
```

**Note:** The `omitempty` tag does NOT use `IsZero()` and will include `null` values. Use `omitzero` for proper 3-state omission behavior.

**SQL NULL scanning:**

Control how SQL NULL scans:

```go
// Package-level default (default: ScanNullAsNull)
presence.SetDefaultScanNull(presence.ScanNullAsUnset)

// Per-value override
val := presence.Of[string]{}
val.SetScanNull(presence.ScanNullAsUnset)
```

## Why Use This Library?

### Standard `database/sql` Approach

```go
type User struct {
    Name sql.NullString `json:"name"`
    Age  sql.NullInt64  `json:"age"`
}

// JSON output:
// {"name":{"String":"John","Valid":true},"age":{"Int64":30,"Valid":true}}
```

### With `presence`

```go
type User struct {
    Name presence.Of[string] `json:"name"`
    Age  presence.Of[int]    `json:"age"`
}

// JSON output:
// {"name":"John","age":30}
// or with null values:
// {"name":null,"age":null}
```

## Usage Examples

### Basic Usage

```go
package main

import (
    "encoding/json"
    "fmt"
    "github.com/pivaldi/presence"
)

type User struct {
    ID       presence.Of[int]    `json:"id"`
    Name     presence.Of[string] `json:"name"`
    Email    presence.Of[string] `json:"email"`
    Age      presence.Of[int]    `json:"age"`
    IsActive presence.Of[bool]   `json:"isActive"`
}

func main() {
    // Create user with some null fields
    user := User{
        ID:       presence.FromValue(1),
        Name:     presence.FromValue("John Doe"),
        Email:    presence.Null[string](), // Null email
        Age:      presence.FromValue(30),
        IsActive: presence.FromValue(true),
    }

    // Marshal to JSON
    data, _ := json.Marshal(user)
    fmt.Println(string(data))
    // Output: {"id":1,"name":"John Doe","email":null,"age":30,"isActive":true}

    // Unmarshal from JSON
    jsonStr := `{"id":2,"name":"Jane Doe","email":"jane@example.com","age":null,"isActive":false}`
    var user2 User
    json.Unmarshal([]byte(jsonStr), &user2)

    fmt.Println(*user2.Name.GetValue())    // "Jane Doe"
    fmt.Println(user2.Age.IsNull())        // true
}
```

### Database Operations

#### Insert

```go
import (
    "database/sql"
    "time"
    "github.com/pivaldi/presence"
    _ "github.com/jackc/pgx/v5/stdlib"
)

type Article struct {
    ID          int64                  `db:"id"`
    Title       presence.Of[string]    `db:"title"`
    Content     presence.Of[string]    `db:"content"`
    PublishedAt presence.Of[time.Time] `db:"published_at"`
    AuthorID    presence.Of[int64]     `db:"author_id"`
}

func insertArticle(db *sql.DB) error {
    article := Article{
        Title:       presence.FromValue("My Article"),
        Content:     presence.FromValue("Article content here..."),
        PublishedAt: presence.FromValue(time.Now()),
        AuthorID:    presence.Null[int64](), // Anonymous article
    }

    query := `
        INSERT INTO articles (title, content, published_at, author_id)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

    return db.QueryRow(
        query,
        article.Title,
        article.Content,
        article.PublishedAt,
        article.AuthorID,
    ).Scan(&article.ID)
}
```

#### Query

```go
func getArticle(db *sql.DB, id int64) (*Article, error) {
    var article Article

    query := `
        SELECT id, title, content, published_at, author_id
        FROM articles
        WHERE id = $1
    `

    err := db.QueryRow(query, id).Scan(
        &article.ID,
        &article.Title,
        &article.Content,
        &article.PublishedAt,
        &article.AuthorID,
    )

    if err != nil {
        return nil, err
    }

    return &article, nil
}
```

### Working with JSON/JSONB (PostgreSQL)

Store complex Go types as JSON in PostgreSQL. Simply use the struct type directly - no wrapper needed:

```go
type Metadata struct {
    Tags       []string          `json:"tags"`
    Properties map[string]string `json:"properties"`
    Version    int               `json:"version"`
}

type Document struct {
    ID       int64                   `db:"id"`
    Title    presence.Of[string]     `db:"title"`
    Metadata presence.Of[Metadata]   `db:"metadata"` // Stored as JSONB
}

func insertDocument(db *sql.DB) error {
    meta := Metadata{
        Tags:       []string{"golang", "database"},
        Properties: map[string]string{"type": "article", "lang": "en"},
        Version:    1,
    }

    doc := Document{
        Title:    presence.FromValue("Go Presence Guide"),
        Metadata: presence.FromValue(meta),
    }

    query := `INSERT INTO documents (title, metadata) VALUES ($1, $2) RETURNING id`
    return db.QueryRow(query, doc.Title, doc.Metadata).Scan(&doc.ID)
}
```

### Nested Structures

Use types directly without any wrapper - the library handles them automatically:

```go
type Address struct {
    Street  presence.Of[string] `json:"street"`
    City    presence.Of[string] `json:"city"`
    ZipCode presence.Of[string] `json:"zipCode"`
}

type Profile struct {
    Bio     presence.Of[string]  `json:"bio"`
    Website presence.Of[string]  `json:"website"`
    Address presence.Of[Address] `json:"address"`
}

type User struct {
    Username presence.Of[string]  `json:"username"`
    Email    presence.Of[string]  `json:"email"`
    Profile  presence.Of[Profile] `json:"profile"`
}

func main() {
    user := User{
        Username: presence.FromValue("johndoe"),
        Email:    presence.FromValue("john@example.com"),
        Profile: presence.FromValue(Profile{
            Bio:     presence.FromValue("Software Developer"),
            Website: presence.FromValue("https://johndoe.com"),
            Address: presence.FromValue(Address{
                Street:  presence.FromValue("123 Main St"),
                City:    presence.FromValue("New York"),
                ZipCode: presence.FromValue("10001"),
            }),
        }),
    }

    data, _ := json.MarshalIndent(user, "", "  ")
    fmt.Println(string(data))
}
```

### Custom Types with Scanner/Valuer

For custom primitive types that should be stored as their underlying type (not JSON):

```go
import (
    "database/sql/driver"
    "errors"
    "fmt"
    "strconv"
)

type PhoneNumber string

// Value implements driver.Valuer to store as string in database
func (pn PhoneNumber) Value() (driver.Value, error) {
    return string(pn), nil
}

// Scan implements sql.Scanner to read from database
func (pn *PhoneNumber) Scan(v any) error {
    switch val := v.(type) {
    case int, int64, uint64:
        *pn = PhoneNumber(strconv.Itoa(val.(int)))
    case string:
        *pn = PhoneNumber(val)
    default:
        return errors.New(fmt.Sprintf("cannot scan phone number from type %T", val))
    }
    return nil
}

// Now PhoneNumber will be stored as string, not JSON
type Contact struct {
    Email presence.Of[string]      `db:"email"`
    Phone presence.Of[PhoneNumber] `db:"phone"` // Stored as string, not JSON
}
```

## API Reference

### Creating Presence Values

```go
// From a value
name := presence.FromValue("John")

// Explicitly null
email := presence.Null[string]()

// From a pointer (nil pointer becomes null)
value := presence.FromPtr(ptr) // Returns null if ptr is nil

// From a boolean condition
value := presence.FromBool("value", ok) // Returns null if ok is false

// Using SetValueP
var val presence.Of[string]
val.SetValueP(ptr) // Sets to null if ptr is nil
```

### Checking and Accessing Values

```go
// Check state
if value.IsUnset() {
    // Handle unset (field was never touched)
}
if value.IsNull() {
    // Handle explicit null
}
if value.IsSet() {
    // Field has a value or is explicitly null
}
if value.IsValue() {
    // Field has a concrete value (not null, not unset)
}

// Get value - multiple options
v := value.GetValue()           // Returns *T (nil if null/unset)
v, ok := value.Get()            // Returns (T, bool)
v := value.GetOr("default")     // Returns T or default
v := value.MustGet()            // Returns T or panics
ptr := value.Ptr()              // Returns *T (nil if null/unset)
```

### Setting Values

```go
var value presence.Of[string]

// Set a value
value.SetValue("hello")

// Set from pointer
str := "world"
value.SetValueP(&str)

// Set to null
value.SetNull()

// Reset to unset
value.Unset()
```

### JSON Operations

```go
// Marshal
data, err := json.Marshal(value)

// Unmarshal
var value presence.Of[string]
err := json.Unmarshal([]byte(`"hello"`), &value)

// Unmarshal null
err := json.Unmarshal([]byte(`null`), &value)
// value.IsNull() == true
```

### Functional Operations

```go
// Map - transform the value (package-level function due to Go generics limitations)
age := presence.FromValue(25)
ageStr := presence.Map(age, func(a int) string {
    return fmt.Sprintf("%d years old", a)
})
// ageStr contains "25 years old"

// MapOr - transform or return default
result := presence.MapOr(age, "unknown", func(a int) string {
    return fmt.Sprintf("%d years old", a)
})

// FlatMap - transform to another presence
user := presence.FlatMap(userID, func(id int) presence.Of[User] {
    return fetchUser(id) // returns presence.Of[User]
})

// Filter - keep value only if predicate passes
adult := presence.Filter(age, func(a int) bool {
    return a >= 18
})
// Returns null if age < 18

// Or - return first non-null value
name := presence.Or(preferredName, displayName, defaultName)
```

## Testing

Run all tests including PostgreSQL integration tests using `go test`:

```bash
cd tests
go test -v ./...
```

Or from the root using `gotestsum`:

```bash
make test
```

**Requirements:**
- Docker must be running (testcontainers uses Docker to spin up PostgreSQL)
- No manual database setup needed - testcontainers handles everything

**First run:** Tests will download the PostgreSQL 18 image (~80MB), subsequent runs use cached image.

Run only unit tests (no database required):

```bash
cd tests
go test -run 'TestMarshal|TestUnmarshal|TestPresenceEdgeCases' -v
```

## Examples

### gorm.io/gen Integration

For automatic model generation from database schemas using [gorm.io/gen](https://github.com/go-gorm/gen), see the example in [`examples/gorm-gen/main.go`](examples/gorm-gen/main.go).

The example demonstrates:
- Using `WithDataTypeMap` to wrap nullable columns with `presence.Of[T]`
- Custom type mappings for all PostgreSQL types (strings, integers, floats, booleans, dates, json, uuid)
- Adding required import paths for generated code

Key configuration snippet:
```go
// Helper to wrap nullable columns with presence.Of[T]
func wrapNullable(c gorm.ColumnType, baseType string) string {
    if nullable, _ := c.Nullable(); nullable {
        return fmt.Sprintf("presence.Of[%s]", baseType)
    }
    return baseType
}

// Type mapping functions
var dataTypeMap = map[string]func(gorm.ColumnType) string{
    "varchar": func(c gorm.ColumnType) string { return wrapNullable(c, "string") },
    "int4":    func(c gorm.ColumnType) string { return wrapNullable(c, "int64") },
    "bool":    func(c gorm.ColumnType) string { return wrapNullable(c, "bool") },
    // ... more types
}

config := gen.Config{
    FieldNullable: false, // We handle nullable via WithDataTypeMap
    // ... other config
}
config.WithImportPkgPath("github.com/pivaldi/presence")

g := gen.NewGenerator(config)
g.WithDataTypeMap(dataTypeMap)
```

## Comparison with Alternatives

| Feature | `presence` | `aarondl/opt` | `lomsa-dev/gonull` | `database/sql.Null*` | `guregu/null.v4` |
|---------|-----------|---------------|-------------------|---------------------|------------------|
| Generic (any type) | ✅ `Of[T any]` | ✅ `Val[T any]` | ✅ `Nullable[T]` | ❌ (separate type per kind) | ❌ (separate type per kind) |
| Clean JSON output | ✅ `null` | ✅ `null` | ✅ `null` | ❌ `{"Valid":false}` | ✅ `null` |
| 3-state model | ✅ (unset/null/value) | ✅ (unset/null/value) | ❌ (null/value only) | ❌ | ⚠️ Limited |
| PostgreSQL JSON/JSONB | ✅ Optimized | ✅ Generic | ❌ | ❌ | ⚠️ Limited |
| UUID support | ✅ Built-in | ✅ Any type | ❌ | ❌ | ❌ |
| Custom types | ✅ via Scanner/Valuer | ✅ via Scanner/Valuer | ✅ via Scanner/Valuer | ✅ via Scanner/Valuer | ✅ via Scanner/Valuer |
| Configurable behavior | ✅ Per-value and package-level | ❌ | ❌ | ❌ | ❌ |
| Functional operations | ✅ `Map()`, `Filter()`, etc. | ✅ `Map()`, etc. | ❌ | ❌ | ❌ |
| Package structure | Single type | 3 sub-packages | Single type | N/A | N/A |
| Zero dependencies* | ✅ | ✅ | ✅ | ✅ | ❌ |

*Except `google/uuid` for UUID support

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

See [LICENSE](LICENSE) file for details.
