package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// See https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
func unmarshalBody(r *http.Request, dst interface{}) error {
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequestError{status: http.StatusBadRequest, msg: msg, err: err}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return &malformedRequestError{status: http.StatusBadRequest, msg: msg, err: err}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequestError{status: http.StatusBadRequest, msg: msg, err: err}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequestError{status: http.StatusBadRequest, msg: msg, err: err}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &malformedRequestError{status: http.StatusBadRequest, msg: msg, err: err}

		case err.Error() == "http: request body too large":
			msg := "Request body is too large"
			return &malformedRequestError{status: http.StatusRequestEntityTooLarge, msg: msg, err: err}

		default:
			return &malformedRequestError{status: http.StatusBadRequest, msg: err.Error(), err: err}
		}
	}

	if dec.More() {
		msg := "Request body must only contain a single JSON object"
		return &malformedRequestError{status: http.StatusBadRequest, msg: msg, err: err}
	}

	return nil
}
