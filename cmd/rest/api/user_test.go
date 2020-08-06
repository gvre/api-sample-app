// +build api

package api_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gvre/rest-api-sample-app/app"
)

func TestUsers(t *testing.T) {
	tt := []testcase{
		{
			"GetAllUsers",
			"GET",
			"/users",
			nil,
			http.StatusOK,
			"application/json",
			func(t *testing.T, b []byte) {
				var users []app.User
				if err := json.Unmarshal(b, &users); err != nil {
					t.Errorf("Expected nil error, got %q", err)
				}

				if len(users) == 0 {
					t.Error("Expected some users, got none")
				}
			},
		},

		{
			"GetUser",
			"GET",
			"/users/1",
			nil,
			http.StatusOK,
			"application/json",
			func(t *testing.T, b []byte) {
				var user app.User
				if err := json.Unmarshal(b, &user); err != nil {
					t.Errorf("Expected nil error, got %q", err)
				}

				if user.ID != 1 {
					t.Errorf("Expected ID 1, got %d", user.ID)
				}
			},
		},

		{
			"GetUserDoesNotExist",
			"GET",
			"/users/1234567890",
			nil,
			http.StatusInternalServerError,
			"text/plain",
			nil,
		},

		{
			"AddUser",
			"POST",
			"/users",
			[]byte(`{"name":"test"}`),
			http.StatusCreated,
			"application/json",
			func(t *testing.T, b []byte) {
				var res struct {
					ID int `json:"id"`
				}
				if err := json.Unmarshal(b, &res); err != nil {
					t.Errorf("Expected nil error, got %q", err)
				}

				if res.ID < 1 {
					t.Errorf("Expected ID, got %d", res.ID)
				}
			},
		},

		{
			"AddUserMalformedRequest",
			"POST",
			"/users",
			[]byte(`{"name":"`),
			http.StatusBadRequest,
			"application/json",
			nil,
		},

		{
			"AddUserUnknownField",
			"POST",
			"/users",
			[]byte(`{"name":"test","nonexistentfield":""`),
			http.StatusBadRequest,
			"application/json",
			nil,
		},

		{
			"AddUserMultipleJSON",
			"POST",
			"/users",
			[]byte(`{"name":"test"}{"name":"test"}`),
			http.StatusBadRequest,
			"application/json",
			nil,
		},
	}

	run(t, tt)
}
