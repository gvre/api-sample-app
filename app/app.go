package app

import "context"

// A User represents a user entity.
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// UserRepository should be implemented to get access to the data store.
type UserRepository interface {
	FetchAll(ctx context.Context) ([]User, error)
	FetchByID(ctx context.Context, userID int) (*User, error)
	Add(ctx context.Context, name string) (int, error)
}
