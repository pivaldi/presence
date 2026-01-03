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
