# gqlgen Example Design

## Overview

A GraphQL example demonstrating how `presence.Of[T]` handles the 3-state model (unset/null/value) for PATCH mutations using gqlgen.

## Problem Statement

GraphQL mutations with optional input fields suffer from ambiguity:
- Field not sent → don't update
- Field sent as `null` → clear the value
- Field sent with value → update to new value

Standard pointer-based approaches (`*string`) cannot distinguish "not sent" from "sent as null" - both result in `nil`.

## Solution

Use `presence.Of[T]` for input fields to get proper 3-state handling:
- `IsUnset()` → field not sent
- `IsNull()` → explicitly set to null
- `IsValue()` → has concrete value

## Schema

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
}

type Mutation {
  updateUser(id: ID!, input: UpdateUserInput!): User
}
```

## Directory Structure

```
examples/gqlgen/
├── go.mod
├── gqlgen.yml
├── graph/
│   ├── schema.graphqls
│   ├── model/
│   │   └── models.go
│   ├── resolver.go
│   └── schema.resolvers.go
├── server.go
└── README.md
```

## Custom Model

```go
// graph/model/models.go
package model

import "github.com/pivaldi/presence"

type User struct {
    ID       string
    Username string
    Email    *string
    Bio      *string
    Website  *string
    Age      *int
}

type UpdateUserInput struct {
    Username presence.Of[string] `json:"username"`
    Email    presence.Of[string] `json:"email"`
    Bio      presence.Of[string] `json:"bio"`
    Website  presence.Of[string] `json:"website"`
    Age      presence.Of[int]    `json:"age"`
}
```

## gqlgen Configuration

```yaml
# gqlgen.yml
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

models:
  UpdateUserInput:
    model: github.com/pivaldi/presence/examples/gqlgen/graph/model.UpdateUserInput
  User:
    model: github.com/pivaldi/presence/examples/gqlgen/graph/model.User
```

## Resolver Logic

```go
func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input model.UpdateUserInput) (*model.User, error) {
    user, ok := r.users[id]
    if !ok {
        return nil, fmt.Errorf("user not found: %s", id)
    }

    if input.Username.IsSet() {
        if input.Username.IsNull() {
            return nil, fmt.Errorf("username cannot be null")
        }
        user.Username = input.Username.MustGet()
    }

    if input.Email.IsSet() {
        user.Email = input.Email.Ptr()
    }

    if input.Bio.IsSet() {
        user.Bio = input.Bio.Ptr()
    }

    if input.Website.IsSet() {
        user.Website = input.Website.Ptr()
    }

    if input.Age.IsSet() {
        user.Age = input.Age.Ptr()
    }

    r.users[id] = user
    return user, nil
}
```

## Seed Data

```go
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
}
```

## Example Queries

```graphql
# Get user
query { user(id: "1") { id username email bio } }

# Update only username (other fields untouched)
mutation { updateUser(id: "1", input: {username: "newname"}) { id username bio } }

# Clear bio explicitly (set to null)
mutation { updateUser(id: "1", input: {bio: null}) { id username bio } }

# Update multiple fields, clear website
mutation { updateUser(id: "1", input: {username: "john", website: null}) { id username website } }
```

## Key Patterns Demonstrated

| Method | Use Case |
|--------|----------|
| `IsSet()` | Check if field was sent at all |
| `IsNull()` | Check if explicitly set to null |
| `IsValue()` | Check if has concrete value |
| `MustGet()` | Get value for required fields |
| `Ptr()` | Get pointer (nil for null) for optional fields |

## Server

Simple HTTP server with GraphQL Playground at root for interactive testing:

```go
func main() {
    srv := handler.NewDefaultServer(
        generated.NewExecutableSchema(generated.Config{
            Resolvers: &graph.Resolver{},
        }),
    )

    http.Handle("/", playground.Handler("GraphQL", "/query"))
    http.Handle("/query", srv)

    log.Println("Server running at http://localhost:8181/")
    log.Fatal(http.ListenAndServe(":8181", nil))
}
```

## Dependencies

- `github.com/99designs/gqlgen`
- `github.com/pivaldi/presence`
