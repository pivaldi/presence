# gqlgen Example Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Create a GraphQL example demonstrating `presence.Of[T]` for PATCH mutations with proper 3-state handling.

**Architecture:** gqlgen generates GraphQL server code from schema; we provide custom `UpdateUserInput` model using `presence.Of[T]` instead of pointers; resolver demonstrates checking `IsSet()`, `IsNull()`, `IsValue()` for partial updates.

**Tech Stack:** Go 1.25, gqlgen, presence library, in-memory storage

---

### Task 1: Create Directory Structure and go.mod

**Files:**
- Create: `examples/gqlgen/go.mod`

**Step 1: Create the directory**

```bash
mkdir -p examples/gqlgen/graph/model
```

**Step 2: Create go.mod**

Create `examples/gqlgen/go.mod`:

```go
module github.com/pivaldi/presence/examples/gqlgen

go 1.25.0

require (
	github.com/99designs/gqlgen v0.17.66
	github.com/pivaldi/presence v0.0.0
	github.com/vektah/gqlparser/v2 v2.5.22
)

replace github.com/pivaldi/presence => ../..
```

**Step 3: Update go.work**

Add to `go.work`:

```go
use (
	.
	./examples/gorm-gen
	./examples/gqlgen
	./tests
)
```

**Step 4: Verify**

Run: `go mod tidy` from workspace root
Expected: No errors

---

### Task 2: Create GraphQL Schema

**Files:**
- Create: `examples/gqlgen/graph/schema.graphqls`

**Step 1: Create the schema file**

Create `examples/gqlgen/graph/schema.graphqls`:

```graphql
type User {
  id: ID!
  username: String!
  email: String
  bio: String
  website: String
  age: Int
}

input UpdateUserInput {
  username: String
  email: String
  bio: String
  website: String
  age: Int
}

type Query {
  user(id: ID!): User
  users: [User!]!
}

type Mutation {
  updateUser(id: ID!, input: UpdateUserInput!): User
}
```

**Step 2: Verify file exists**

Run: `cat examples/gqlgen/graph/schema.graphqls`
Expected: Schema content displayed

---

### Task 3: Create Custom Models

**Files:**
- Create: `examples/gqlgen/graph/model/models.go`

**Step 1: Create the models file**

Create `examples/gqlgen/graph/model/models.go`:

```go
package model

import "github.com/pivaldi/presence"

// User represents a user in the system.
// Nullable fields use pointers for GraphQL compatibility.
type User struct {
	ID       string
	Username string
	Email    *string
	Bio      *string
	Website  *string
	Age      *int
}

// UpdateUserInput uses presence.Of[T] to distinguish:
// - Field not sent (IsUnset)
// - Field explicitly set to null (IsNull)
// - Field has a value (IsValue)
type UpdateUserInput struct {
	Username presence.Of[string] `json:"username"`
	Email    presence.Of[string] `json:"email"`
	Bio      presence.Of[string] `json:"bio"`
	Website  presence.Of[string] `json:"website"`
	Age      presence.Of[int]    `json:"age"`
}
```

**Step 2: Verify syntax**

Run: `cd examples/gqlgen && go build ./graph/model/`
Expected: No errors

---

### Task 4: Create gqlgen Configuration

**Files:**
- Create: `examples/gqlgen/gqlgen.yml`

**Step 1: Create gqlgen.yml**

Create `examples/gqlgen/gqlgen.yml`:

```yaml
schema:
  - graph/*.graphqls

exec:
  filename: graph/generated/generated.go
  package: generated

model:
  filename: graph/model/models_gen.go
  package: model

resolver:
  layout: follow-schema
  dir: graph
  package: graph
  filename_template: "{name}.resolvers.go"

models:
  ID:
    model:
      - github.com/99designs/gqlgen/graphql.ID
  Int:
    model:
      - github.com/99designs/gqlgen/graphql.Int
  UpdateUserInput:
    model: github.com/pivaldi/presence/examples/gqlgen/graph/model.UpdateUserInput
  User:
    model: github.com/pivaldi/presence/examples/gqlgen/graph/model.User
```

**Step 2: Verify YAML syntax**

Run: `cat examples/gqlgen/gqlgen.yml`
Expected: Valid YAML displayed

---

### Task 5: Generate gqlgen Code

**Files:**
- Creates: `examples/gqlgen/graph/generated/generated.go`
- Creates: `examples/gqlgen/graph/schema.resolvers.go`
- Creates: `examples/gqlgen/graph/resolver.go`

**Step 1: Install gqlgen tool**

Run: `cd examples/gqlgen && go get github.com/99designs/gqlgen`

**Step 2: Run gqlgen generate**

Run: `cd examples/gqlgen && go run github.com/99designs/gqlgen generate`
Expected: Files generated in `graph/` directory

**Step 3: Verify generation**

Run: `ls examples/gqlgen/graph/generated/`
Expected: `generated.go` exists

---

### Task 6: Implement Resolver with Seed Data

**Files:**
- Modify: `examples/gqlgen/graph/resolver.go`

**Step 1: Replace resolver.go with seed data**

Replace `examples/gqlgen/graph/resolver.go` with:

```go
package graph

import "github.com/pivaldi/presence/examples/gqlgen/graph/model"

// Resolver is the root resolver with in-memory user storage.
type Resolver struct {
	users map[string]*model.User
}

// NewResolver creates a resolver with seed data.
func NewResolver() *Resolver {
	return &Resolver{
		users: map[string]*model.User{
			"1": {
				ID:       "1",
				Username: "alice",
				Email:    ptr("alice@example.com"),
				Bio:      ptr("Software developer"),
				Website:  nil,
				Age:      ptr(30),
			},
			"2": {
				ID:       "2",
				Username: "bob",
				Email:    nil,
				Bio:      nil,
				Website:  ptr("https://bob.dev"),
				Age:      nil,
			},
		},
	}
}

func ptr[T any](v T) *T {
	return &v
}
```

**Step 2: Verify syntax**

Run: `cd examples/gqlgen && go build ./graph/`
Expected: No errors (resolvers not implemented yet, but should compile)

---

### Task 7: Implement Query Resolvers

**Files:**
- Modify: `examples/gqlgen/graph/schema.resolvers.go`

**Step 1: Implement User query**

In `schema.resolvers.go`, implement the `User` method:

```go
// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, nil // Return nil for not found (GraphQL nullable)
	}
	return user, nil
}
```

**Step 2: Implement Users query**

In `schema.resolvers.go`, implement the `Users` method:

```go
// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	users := make([]*model.User, 0, len(r.users))
	for _, u := range r.users {
		users = append(users, u)
	}
	return users, nil
}
```

**Step 3: Verify syntax**

Run: `cd examples/gqlgen && go build ./graph/`
Expected: No errors

---

### Task 8: Implement UpdateUser Mutation

**Files:**
- Modify: `examples/gqlgen/graph/schema.resolvers.go`

**Step 1: Implement UpdateUser mutation**

In `schema.resolvers.go`, implement the `UpdateUser` method:

```go
// UpdateUser is the resolver for the updateUser field.
// Demonstrates presence.Of[T] 3-state handling for PATCH semantics.
func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input model.UpdateUserInput) (*model.User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, fmt.Errorf("user not found: %s", id)
	}

	// Username: required field, cannot be set to null
	if input.Username.IsSet() {
		if input.Username.IsNull() {
			return nil, fmt.Errorf("username cannot be null")
		}
		user.Username = input.Username.MustGet()
	}

	// Email: optional, can be set to null to clear
	if input.Email.IsSet() {
		user.Email = input.Email.Ptr()
	}

	// Bio: optional, can be set to null to clear
	if input.Bio.IsSet() {
		user.Bio = input.Bio.Ptr()
	}

	// Website: optional, can be set to null to clear
	if input.Website.IsSet() {
		user.Website = input.Website.Ptr()
	}

	// Age: optional, can be set to null to clear
	if input.Age.IsSet() {
		user.Age = input.Age.Ptr()
	}

	return user, nil
}
```

**Step 2: Add fmt import if needed**

Ensure `"fmt"` is in the imports.

**Step 3: Verify syntax**

Run: `cd examples/gqlgen && go build ./graph/`
Expected: No errors

---

### Task 9: Create Server Entry Point

**Files:**
- Create: `examples/gqlgen/server.go`

**Step 1: Create server.go**

Create `examples/gqlgen/server.go`:

```go
package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/pivaldi/presence/examples/gqlgen/graph"
	"github.com/pivaldi/presence/examples/gqlgen/graph/generated"
)

func main() {
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{
			Resolvers: graph.NewResolver(),
		}),
	)

	http.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	http.Handle("/query", srv)

	log.Println("GraphQL server running at http://localhost:8181/")
	log.Println("Open http://localhost:8181/ for GraphQL Playground")
	log.Fatal(http.ListenAndServe(":8181", nil))
}
```

**Step 2: Verify build**

Run: `cd examples/gqlgen && go build .`
Expected: Binary created successfully

---

### Task 10: Create README

**Files:**
- Create: `examples/gqlgen/README.md`

**Step 1: Create README.md**

Create `examples/gqlgen/README.md`:

```markdown
# gqlgen Example: PATCH Mutations with presence.Of[T]

This example demonstrates how `presence.Of[T]` enables proper 3-state handling for GraphQL PATCH mutations.

## The Problem

GraphQL mutations with optional input fields have ambiguous semantics:
- Field not sent → don't update
- Field sent as `null` → clear the value
- Field sent with value → update to new value

Standard pointer-based approaches (`*string`) cannot distinguish "not sent" from "sent as null".

## The Solution

Use `presence.Of[T]` for input fields:

```go
type UpdateUserInput struct {
    Username presence.Of[string] `json:"username"`
    Email    presence.Of[string] `json:"email"`
    // ...
}
```

Then in the resolver:

```go
if input.Email.IsSet() {      // Was the field sent?
    if input.Email.IsNull() { // Was it explicitly null?
        user.Email = nil
    } else {
        user.Email = input.Email.Ptr()
    }
}
// If not IsSet(), don't touch the field
```

## Running the Example

```bash
cd examples/gqlgen
go run .
```

Open http://localhost:8181/ for GraphQL Playground.

## Example Queries

### Get all users
```graphql
query {
  users {
    id
    username
    email
    bio
    website
    age
  }
}
```

### Get single user
```graphql
query {
  user(id: "1") {
    id
    username
    email
    bio
  }
}
```

### Update only username (other fields untouched)
```graphql
mutation {
  updateUser(id: "1", input: {username: "alice_new"}) {
    id
    username
    bio
  }
}
```

### Clear bio explicitly (set to null)
```graphql
mutation {
  updateUser(id: "1", input: {bio: null}) {
    id
    username
    bio
  }
}
```

### Update multiple fields, clear website
```graphql
mutation {
  updateUser(id: "2", input: {
    username: "bobby",
    email: "bob@newmail.com",
    website: null
  }) {
    id
    username
    email
    website
  }
}
```

## Key Methods

| Method | Use Case |
|--------|----------|
| `IsSet()` | Check if field was sent at all |
| `IsNull()` | Check if explicitly set to null |
| `IsValue()` | Check if has concrete value |
| `MustGet()` | Get value (panics if null/unset) |
| `Ptr()` | Get pointer (nil if null/unset) |
```

**Step 2: Verify file**

Run: `cat examples/gqlgen/README.md`
Expected: README content displayed

---

### Task 11: Run go mod tidy and Verify

**Step 1: Tidy dependencies**

Run: `cd examples/gqlgen && go mod tidy`
Expected: Dependencies resolved

**Step 2: Build final binary**

Run: `cd examples/gqlgen && go build .`
Expected: No errors

**Step 3: Verify workspace build**

Run: `go build ./...` (from workspace root)
Expected: All modules build successfully

---

### Task 12: Manual Testing

**Step 1: Start the server**

Run: `cd examples/gqlgen && go run .`
Expected: "GraphQL server running at http://localhost:8181/"

**Step 2: Test in playground**

Open http://localhost:8181/ and run:

```graphql
query { users { id username email bio website age } }
```

Expected: Two users returned (alice and bob)

**Step 3: Test PATCH mutation**

```graphql
mutation { updateUser(id: "1", input: {bio: null}) { id username bio } }
```

Expected: bio is now null, username unchanged

**Step 4: Stop server**

Press Ctrl+C to stop

---

## Summary

| Task | Description |
|------|-------------|
| 1 | Create directory and go.mod |
| 2 | Create GraphQL schema |
| 3 | Create custom models with presence.Of[T] |
| 4 | Create gqlgen.yml configuration |
| 5 | Generate gqlgen code |
| 6 | Implement resolver with seed data |
| 7 | Implement query resolvers |
| 8 | Implement UpdateUser mutation |
| 9 | Create server entry point |
| 10 | Create README |
| 11 | Tidy and verify build |
| 12 | Manual testing |
