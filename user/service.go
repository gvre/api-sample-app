package user

import (
	"context"

	multierror "github.com/hashicorp/go-multierror"

	"github.com/gvre/rest-api-sample-app/app"
)

// Service wraps the user repository.
type Service struct {
	repo app.UserRepository
}

// NewService returns a new Service.
func NewService(repo app.UserRepository) *Service {
	return &Service{repo: repo}
}

// Users returns all application users.
// Pagination should be implemented when the number of users grows.
func (s *Service) Users(ctx context.Context) ([]app.User, error) {
	return s.repo.FetchAll(ctx)
}

// User returns the user with the provided id.
func (s *Service) User(ctx context.Context, id int) (*app.User, error) {
	var err error

	if id < 1 {
		err = multierror.Append(err, &app.ValidationError{
			Field:   "id",
			Message: "Field cannot be less than 1",
		})
		return nil, err
	}

	return s.repo.FetchByID(ctx, id)
}

// Add inserts a new user and returns its id.
func (s *Service) Add(ctx context.Context, name string) (int, error) {
	return s.repo.Add(ctx, name)
}
