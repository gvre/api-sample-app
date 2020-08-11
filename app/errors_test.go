package app_test

import (
	"errors"
	"net/url"
	"syscall"
	"testing"

	"github.com/gvre/api-sample-app/app"
)

func TestNewErrorFunc(t *testing.T) {
	msg := "message"
	err := app.NewError(msg, nil)
	if app.ErrorMessage(err) != msg {
		t.Errorf("Expected %q, got %q", msg, app.ErrorMessage(err))
	}
}

func TestErrorWithNilAppError(t *testing.T) {
	var err *app.Error
	if err.Error() != "" {
		t.Errorf("Expected empty error message, got %q", err.Error())
	}
}

func TestErrorMessageWithNilError(t *testing.T) {
	var err error
	if app.ErrorMessage(err) != "" {
		t.Errorf("Expected empty error message, got %q", app.ErrorMessage(err))
	}
}

func TestErrorMessage(t *testing.T) {
	err := &app.Error{
		Msg: "message",
	}

	if app.ErrorMessage(err) != err.Msg {
		t.Errorf("Expected %q, got %q", err.Msg, app.ErrorMessage(err))
	}
}

func TestErrorWithWrappedError(t *testing.T) {
	msg := "message"
	err := &app.Error{
		Err: errors.New(msg),
	}

	if err.Error() != msg {
		t.Errorf("Expected %q, got %q", msg, err.Error())
	}
}

func TestErrorWithoutWrappedError(t *testing.T) {
	msg := "message"
	err := &app.Error{
		Msg: "message",
	}

	if err.Error() != msg {
		t.Errorf("Expected %q, got %q", msg, err.Error())
	}
}

func TestErrorMessageWithErrorsPackage(t *testing.T) {
	err := errors.New("test")

	if app.ErrorMessage(err) != "test" {
		t.Errorf("Expected %q error message, got %q", "test", app.ErrorMessage(err))
	}
}

func TestErrorsWithDeepNesting(t *testing.T) {
	err := &app.Error{
		Err: &url.Error{
			Err: syscall.ECONNRESET,
		},
	}

	// app.Error
	var appError *app.Error
	if !errors.As(err, &appError) {
		t.Errorf("Expected app.Error, got %#v", err)
	}

	// url.Error
	var urlError *url.Error
	if !errors.As(err, &urlError) {
		t.Errorf("Expected url.Error, got %#v", err)
	}

	// syscall.ECONNRESET
	if !errors.Is(err, syscall.ECONNRESET) {
		t.Errorf("Expected syscall.ECONNRESET, got %#v", err)
	}
}

func TestFuncWithErrorInterface(t *testing.T) {
	err := func() error {
		return &app.Error{}
	}()

	var appError *app.Error
	if !errors.As(err, &appError) {
		t.Errorf("Expected %#v, got %#v", appError, err)
	}
}

func TestNilErrorFuncWithErrorInterface(t *testing.T) {
	err := func() error {
		var err *app.Error
		return err
	}()

	var appError *app.Error
	if !errors.As(err, &appError) {
		t.Errorf("Expected %#v, got %#v", appError, err)
	}
}
