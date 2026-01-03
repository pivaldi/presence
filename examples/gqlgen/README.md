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
