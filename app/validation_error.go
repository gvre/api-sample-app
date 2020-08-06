package app

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
)

// ValidationError defines a standard application error.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error returns the error field and message as a combined string.
func (err *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", err.Field, err.Message)
}

// IsValidationError checks if the provided error is related to validation.
func IsValidationError(err error) bool {
	_, ok := err.(*multierror.Error)

	return ok
}

// WrappedErrors returns the wrapped validation errors.
// If err is not a validation error, nil is returned.
func WrappedErrors(err error) []error {
	if merr, ok := err.(*multierror.Error); ok {
		return merr.Errors
	}

	return nil
}
