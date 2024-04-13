package api

type malformedRequestError struct {
	status int
	msg    string
	err    error
}

func (mr *malformedRequestError) Error() string {
	return mr.msg
}

// ApiError defines the standard error the API returns.
type ApiError struct {
	Message string  `json:"message"`
	Errors  []error `json:"errors,omitempty"`
}

// ErrorWrapper wraps an ApiError, so consumers can check for it easier.
type ErrorWrapper struct {
	Error ApiError `json:"error"`
}
