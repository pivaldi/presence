# Go Nullable

[![golangci-lint](https://github.com/pivaldi/nullable/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/pivaldi/nullable/actions/workflows/golangci-lint.yml)
[![mod-verify](https://github.com/pivaldi/nullable/actions/workflows/mod-verify.yml/badge.svg)](https://github.com/pivaldi/nullable/actions/workflows/mod-verify.yml)
[![gosec](https://github.com/pivaldi/nullable/actions/workflows/gosec.yaml/badge.svg)](https://github.com/pivaldi/nullable/actions/workflows/gosec.yaml)
[![staticcheck](https://github.com/pivaldi/nullable/actions/workflows/staticcheck.yaml/badge.svg)](https://github.com/pivaldi/nullable/actions/workflows/staticcheck.yaml)
[![test](https://github.com/pivaldi/nullable/actions/workflows/test.yml/badge.svg)](https://github.com/pivaldi/nullable/actions/workflows/test.yml)

A type-safe nullable value library for Go using generics, designed for seamless JSON marshaling and database operations.

## Features

- **Type-safe nullable values** for any supported type using Go generics
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
go get github.com/pivaldi/nullable
```

## Quick Start

```go
import "github.com/pivaldi/nullable"

// Create nullable values
name := nullable.FromValue("John Doe")
age := nullable.FromValue(30)
email := nullable.Null[string]() // Explicitly null

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

The library supports a 3-state model for nullable values, enabling PATCH API semantics and partial updates:

| State | Description | Creation |
|-------|-------------|----------|
| Unset | Field was never touched | `nullable.Of[T]{}` or `var x nullable.Of[T]` |
| Null | Explicitly set to null | `nullable.Null[T]()` |
| Value | Has a concrete value | `nullable.FromValue(x)` |

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
    Name  nullable.Of[string] `json:"name,omitempty"`
    Email nullable.Of[string] `json:"email,omitempty"`
    Age   nullable.Of[int]    `json:"age,omitempty"`
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
    Name nullable.Of[string] `json:"name,omitzero"` // omitted when unset (Go 1.24+)
    Age  nullable.Of[int]    `json:"age"`           // always included
}

// Package-level default (default: UnsetSkip)
nullable.SetDefaultMarshalUnset(nullable.UnsetNull)

// Per-value override
val := nullable.Of[string]{}
val.SetMarshalUnset(nullable.UnsetNull)
```

**Note:** The `omitempty` tag does NOT use `IsZero()` and will include `null` values. Use `omitzero` for proper 3-state omission behavior.

**SQL NULL scanning:**

Control how SQL NULL scans:

```go
// Package-level default (default: ScanNullAsNull)
nullable.SetDefaultScanNull(nullable.ScanNullAsUnset)

// Per-value override
val := nullable.Of[string]{}
val.SetScanNull(nullable.ScanNullAsUnset)
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

### With `nullable`

```go
type User struct {
    Name nullable.Of[string] `json:"name"`
    Age  nullable.Of[int]    `json:"age"`
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
    "github.com/pivaldi/nullable"
)

type User struct {
    ID       nullable.Of[int]    `json:"id"`
    Name     nullable.Of[string] `json:"name"`
    Email    nullable.Of[string] `json:"email"`
    Age      nullable.Of[int]    `json:"age"`
    IsActive nullable.Of[bool]   `json:"isActive"`
}

func main() {
    // Create user with some null fields
    user := User{
        ID:       nullable.FromValue(1),
        Name:     nullable.FromValue("John Doe"),
        Email:    nullable.Null[string](), // Null email
        Age:      nullable.FromValue(30),
        IsActive: nullable.FromValue(true),
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
    "github.com/pivaldi/nullable"
    _ "github.com/jackc/pgx/v5/stdlib"
)

type Article struct {
    ID          int64                  `db:"id"`
    Title       nullable.Of[string]    `db:"title"`
    Content     nullable.Of[string]    `db:"content"`
    PublishedAt nullable.Of[time.Time] `db:"published_at"`
    AuthorID    nullable.Of[int64]     `db:"author_id"`
}

func insertArticle(db *sql.DB) error {
    article := Article{
        Title:       nullable.FromValue("My Article"),
        Content:     nullable.FromValue("Article content here..."),
        PublishedAt: nullable.FromValue(time.Now()),
        AuthorID:    nullable.Null[int64](), // Anonymous article
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
    Title    nullable.Of[string]     `db:"title"`
    Metadata nullable.Of[Metadata]   `db:"metadata"` // Stored as JSONB
}

func insertDocument(db *sql.DB) error {
    meta := Metadata{
        Tags:       []string{"golang", "database"},
        Properties: map[string]string{"type": "article", "lang": "en"},
        Version:    1,
    }

    doc := Document{
        Title:    nullable.FromValue("Go Nullable Guide"),
        Metadata: nullable.FromValue(meta),
    }

    query := `INSERT INTO documents (title, metadata) VALUES ($1, $2) RETURNING id`
    return db.QueryRow(query, doc.Title, doc.Metadata).Scan(&doc.ID)
}
```

### Nested Structures

Use types directly without any wrapper - the library handles them automatically:

```go
type Address struct {
    Street  nullable.Of[string] `json:"street"`
    City    nullable.Of[string] `json:"city"`
    ZipCode nullable.Of[string] `json:"zipCode"`
}

type Profile struct {
    Bio     nullable.Of[string]  `json:"bio"`
    Website nullable.Of[string]  `json:"website"`
    Address nullable.Of[Address] `json:"address"`
}

type User struct {
    Username nullable.Of[string]  `json:"username"`
    Email    nullable.Of[string]  `json:"email"`
    Profile  nullable.Of[Profile] `json:"profile"`
}

func main() {
    user := User{
        Username: nullable.FromValue("johndoe"),
        Email:    nullable.FromValue("john@example.com"),
        Profile: nullable.FromValue(Profile{
            Bio:     nullable.FromValue("Software Developer"),
            Website: nullable.FromValue("https://johndoe.com"),
            Address: nullable.FromValue(Address{
                Street:  nullable.FromValue("123 Main St"),
                City:    nullable.FromValue("New York"),
                ZipCode: nullable.FromValue("10001"),
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
    Email nullable.Of[string]      `db:"email"`
    Phone nullable.Of[PhoneNumber] `db:"phone"` // Stored as string, not JSON
}
```

## API Reference

### Creating Nullable Values

```go
// From a value
name := nullable.FromValue("John")

// Explicitly null
email := nullable.Null[string]()

// From a pointer (nil pointer becomes null)
value := nullable.FromPtr(ptr) // Returns null if ptr is nil

// From a boolean condition
value := nullable.FromBool("value", ok) // Returns null if ok is false

// Using SetValueP
var val nullable.Of[string]
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
var value nullable.Of[string]

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
var value nullable.Of[string]
err := json.Unmarshal([]byte(`"hello"`), &value)

// Unmarshal null
err := json.Unmarshal([]byte(`null`), &value)
// value.IsNull() == true
```

### Functional Operations

```go
// Map - transform the value (package-level function due to Go generics limitations)
age := nullable.FromValue(25)
ageStr := nullable.Map(age, func(a int) string {
    return fmt.Sprintf("%d years old", a)
})
// ageStr contains "25 years old"

// MapOr - transform or return default
result := nullable.MapOr(age, "unknown", func(a int) string {
    return fmt.Sprintf("%d years old", a)
})

// FlatMap - transform to another nullable
user := nullable.FlatMap(userID, func(id int) nullable.Of[User] {
    return fetchUser(id) // returns nullable.Of[User]
})

// Filter - keep value only if predicate passes
adult := nullable.Filter(age, func(a int) bool {
    return a >= 18
})
// Returns null if age < 18

// Or - return first non-null value
name := nullable.Or(preferredName, displayName, defaultName)
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
go test -run 'TestMarshal|TestUnmarshal|TestNullableEdgeCases' -v
```

## Comparison with Alternatives

| Feature | `nullable` | `aarondl/opt` | `lomsa-dev/gonull` | `database/sql.Null*` | `guregu/null.v4` |
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
