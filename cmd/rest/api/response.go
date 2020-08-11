package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gvre/api-sample-app/app"
)

// actionKey holds a standard name to be used by the logger.
const actionKey = "action"

// httpStatusIgnore is used when the status the caller passes should be ignored.
const httpStatusIgnore = 0

// Ok encodes to JSON and writes the provided response (if any) along with the httpStatus.
func Ok(w http.ResponseWriter, response interface{}, httpStatus int) {
	// String values encode as JSON strings coerced to valid UTF-8,
	// replacing invalid bytes with the Unicode replacement rune.
	// So that the JSON will be safe to embed inside HTML <script> tags,
	// the string is encoded using HTMLEscape,
	// which replaces "<", ">", "&", U+2028, and U+2029 are escaped
	// to "\u003c","\u003e", "\u0026", "\u2028", and "\u2029".
	// This replacement can be disabled when using an Encoder,
	// by calling SetEscapeHTML(false).
	// See https://github.com/golang/go/blob/release-branch.go1.14/src/encoding/json/encode.go#L46
	buf := new(bytes.Buffer)
	if response != nil {
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(response)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	if buf.Len() > 0 {
		// Remove the extra newline json.Encoder.Encode() adds.
		w.Write(bytes.TrimRight(buf.Bytes(), "\n"))
	}
}

// BadRequestError writes the provided error along with a 400 http status.
func BadRequestError(w http.ResponseWriter, err error) {
	Error(w, err, http.StatusBadRequest)
}

// NotFoundError writes the provided error along with a 404 http status.
func NotFoundError(w http.ResponseWriter, err error) {
	Error(w, err, http.StatusNotFound)
}

// ServerError writes the provided error along with a 500 http status.
func ServerError(w http.ResponseWriter, err error) {
	Error(w, err, http.StatusInternalServerError)
}

// Error writes the provided error along with the provided http status.
func Error(w http.ResponseWriter, err error, httpStatus int) {
	// Malformed request error.
	if e, ok := err.(*malformedRequestError); ok {
		err = &app.Error{
			Msg: e.msg,
			Err: e.err,
		}
		httpStatus = e.status
	}

	// 5xx
	if httpStatus >= http.StatusInternalServerError {
		http.Error(w, http.StatusText(httpStatus), httpStatus)
		return
	}

	apiError := ErrorWrapper{
		Error: ApiError{
			Message: app.ErrorMessage(err),
		},
	}

	res, e := json.Marshal(apiError)
	if e != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	w.Write(res)
}

// ValidationError writes the provided error along with a 422 http status.
// If the error cannot be marshalled, a 500 error is returned instead.
func ValidationError(w http.ResponseWriter, err error) {
	msg := "Validation error"
	apiError := ErrorWrapper{
		Error: ApiError{
			Message: msg,
			Errors:  app.WrappedErrors(err),
		},
	}

	res, e := json.Marshal(apiError)
	if e != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Write(res)
}
