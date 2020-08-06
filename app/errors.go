package app

import (
	"errors"
)

var (
	ErrNoRecords = errors.New("no records")
)

// Error defines a standard application error.
type Error struct {
	Msg string // Human-readable message.
	Err error  // Nested error.
}

// NewError returns a pointer to app.Error.
func NewError(msg string, err error) *Error {
	return &Error{msg, err}
}

// Error returns the error string of the first wrapped error.
// It returns the human-readable message if the wrapped error is nil.
// If receiver is nil, an empty string is returned.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	if e := e.Unwrap(); e != nil {
		return e.Error()
	}

	return e.Message()
}

// Message returns the human-readable message.
func (e *Error) Message() string {
	return e.Msg
}

// Unwrap returns the wrapped error.
// Go 1.13+ added two new functions, "errors.Is" and "errors.As"
// that call "Unwrap", if exists, to check whether any error in Err's chain matches target.
func (e *Error) Unwrap() error {
	return e.Err
}

// ErrorMessage returns the human-readable message of the error, if available.
// Otherwise returns a generic error message.
func ErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	if err, ok := err.(*Error); ok {
		return err.Msg
	}

	return err.Error()
}

// IsNotFoundError returns if the provided error
// is related to no records found.
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrNoRecords)
}
