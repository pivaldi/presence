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
