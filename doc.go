/*
Package presence provides type-safe presence values for Go using generics, designed for seamless JSON marshaling and database operations.

# Overview

The presence package offers a generic way to handle nullable values in Go, with built-in support
for JSON marshaling/unmarshaling and database operations. Unlike standard library sql.Null* types,
presence values marshal to clean JSON (null instead of {"Valid": false, "Value": ...}) and support
any type through Go generics.

# Key Features

  - Type-safe presence values for any type using Of[T any]
  - 3-state model distinguishing unset, null, and value states for PATCH API support
  - Database-friendly with built-in sql.Scanner and driver.Valuer implementations
  - JSON marshaling that uses standard null instead of complex objects
  - PostgreSQL JSON/JSONB support for storing complex types
  - UUID support with github.com/google/uuid
  - Configurable behavior for marshal and scan operations
  - Functional operations (Map, Filter, FlatMap, Or)
  - Zero external dependencies except google/uuid

# Quick Start

Basic usage with presence values:

	import "github.com/pivaldi/presence"

	// Create presence values
	name := presence.FromValue("John Doe")
	age := presence.FromValue(30)
	email := presence.Null[string]()     // Explicitly null
	var phone presence.Of[string]        // Unset (not touched)

	// Check states
	if name.IsNull() {
	    // Handle null case
	}
	if phone.IsUnset() {
	    // Handle unset case
	}

	// Get value
	if !age.IsNull() {
	    fmt.Println(*age.GetValue()) // 30
	}

# Three-State Model

The library supports a 3-state model for presence values, enabling PATCH API semantics:

  - Unset: Field was never touched - presence.Of[T]{} or var x presence.Of[T]
  - Null: Explicitly set to null - presence.Null[T]()
  - Value: Has a concrete value - presence.FromValue(x)

Example PATCH API handling:

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

# Database Integration

The library integrates seamlessly with database/sql through driver.Valuer and sql.Scanner interfaces:

  - Primitive types (string, int*, float64, bool, time.Time, uuid.UUID) are stored directly
  - Custom types implementing sql.Scanner/driver.Valuer use their custom serialization
  - All other types are automatically marshaled to/from JSON for storage

Example database usage:

	type User struct {
	    ID    int64               `db:"id"`
	    Name  presence.Of[string] `db:"name"`
	    Email presence.Of[string] `db:"email"`
	    Age   presence.Of[int]    `db:"age"`
	}

	user := User{
	    Name:  presence.FromValue("John Doe"),
	    Email: presence.Null[string](), // NULL in database
	    Age:   presence.FromValue(30),
	}

	query := `INSERT INTO users (name, email, age) VALUES ($1, $2, $3) RETURNING id`
	err := db.QueryRow(query, user.Name, user.Email, user.Age).Scan(&user.ID)

# JSON Marshaling

Presence values marshal to clean JSON, using standard null for null values:

	type User struct {
	    Name  presence.Of[string] `json:"name"`
	    Email presence.Of[string] `json:"email"`
	    Age   presence.Of[int]    `json:"age"`
	}

	user := User{
	    Name:  presence.FromValue("John"),
	    Email: presence.Null[string](),
	    Age:   presence.FromValue(30),
	}

	data, _ := json.Marshal(user)
	// Output: {"name":"John","email":null,"age":30}

For Go 1.24+, use the omitzero struct tag with UnsetSkip configuration to omit unset fields:

	type Request struct {
	    Name presence.Of[string] `json:"name,omitzero"` // omitted when unset
	    Age  presence.Of[int]    `json:"age"`           // always included
	}

# PostgreSQL JSON/JSONB Support

Store complex Go types as JSON/JSONB in PostgreSQL without any special configuration:

	type Metadata struct {
	    Tags       []string          `json:"tags"`
	    Properties map[string]string `json:"properties"`
	}

	type Document struct {
	    ID       int64                 `db:"id"`
	    Title    presence.Of[string]   `db:"title"`
	    Metadata presence.Of[Metadata] `db:"metadata"` // Stored as JSONB
	}

	meta := Metadata{
	    Tags:       []string{"golang", "database"},
	    Properties: map[string]string{"type": "article"},
	}

	doc := Document{
	    Title:    presence.FromValue("Guide"),
	    Metadata: presence.FromValue(meta),
	}

# Configuration Options

Control marshaling and scanning behavior:

	// Package-level defaults
	presence.SetDefaultMarshalUnset(presence.UnsetNull) // Marshal unset as null
	presence.SetDefaultScanNull(presence.ScanNullAsUnset) // Scan NULL as unset

	// Per-value overrides
	var val presence.Of[string]
	val.SetMarshalUnset(presence.UnsetNull)
	val.SetScanNull(presence.ScanNullAsUnset)

# Functional Operations

Transform and combine presence values with functional operations:

	// Map - transform the value
	age := presence.FromValue(25)
	ageStr := presence.Map(age, func(a int) string {
	    return fmt.Sprintf("%d years old", a)
	})

	// Filter - keep value only if predicate passes
	adult := presence.Filter(age, func(a int) bool {
	    return a >= 18
	})

	// Or - return first non-null value
	name := presence.Or(preferredName, displayName, defaultName)

# API Reference

Creating presence values:

	FromValue(v T) Of[T]           // Create from value
	Null[T]() Of[T]                // Create explicit null
	FromPtr(ptr *T) Of[T]          // Create from pointer (nil becomes null)
	FromBool(v T, ok bool) Of[T]   // Create null if ok is false

Checking state:

	IsUnset() bool    // Field was never touched
	IsNull() bool     // Field is explicitly null
	IsSet() bool      // Field has a value or is null (not unset)
	IsValue() bool    // Field has a concrete value (not null or unset)

Getting values:

	GetValue() *T              // Returns pointer to value (nil if null/unset)
	Get() (T, bool)            // Returns value and presence indicator
	GetOr(defaultVal T) T      // Returns value or default
	MustGet() T                // Returns value or panics
	Ptr() *T                   // Returns pointer to value (nil if null/unset)

Setting values:

	SetValue(v T)              // Set to value
	SetValueP(ptr *T)          // Set from pointer (nil becomes null)
	SetNull()                  // Set to explicit null
	Unset()                    // Reset to unset state

# Custom Types

Custom types implementing sql.Scanner and driver.Valuer are automatically supported:

	type PhoneNumber string

	func (pn PhoneNumber) Value() (driver.Value, error) {
	    return string(pn), nil
	}

	func (pn *PhoneNumber) Scan(v any) error {
	    // Implementation
	    return nil
	}

	// PhoneNumber will be stored using its Value() method, not as JSON
	type Contact struct {
	    Phone presence.Of[PhoneNumber] `db:"phone"`
	}

# Supported Types

The library uses Of[T any] which accepts any type:

  - Primitives: int, int16, int32, int64, float64, bool, string
  - UUID: uuid.UUID (from github.com/google/uuid)
  - Time: time.Time
  - Complex types: structs, slices, maps (stored as JSON in database)
  - Custom types: any type implementing sql.Scanner/driver.Valuer

For more examples and detailed usage, see the package examples and README at:
https://github.com/pivaldi/presence
*/
package presence
