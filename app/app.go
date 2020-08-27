package app

import "context"

// A User represents a user entity.
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

const (
	HealthStatusOK = "ok"
	HealthStatusWarning = "warning"
	HealthStatusError = "error"
)

type Health struct {
	Name      string `json:"name"`
	Status    string `json:"status"` 	// ok|warning|error
	Core      bool   `json:"core"` 		// true for core dependencies
	LatencyMs int64  `json:"latency_ms"`
	Data      struct {
		Message string `json:"message,omitempty"`
		Code int `json:"code,omitempty"` // remote HTTP status code
	} `json:"data"`
}

// UserRepository should be implemented to get access to the data store.
type UserRepository interface {
	Ping(ctx context.Context) error
	FetchAll(ctx context.Context) ([]User, error)
	FetchByID(ctx context.Context, userID int) (*User, error)
	Add(ctx context.Context, name string) (int, error)
}
