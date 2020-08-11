package api_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.uber.org/zap"

	"github.com/gvre/api-sample-app/cmd/rest/api"
	"github.com/gvre/api-sample-app/user"

	"github.com/DATA-DOG/go-txdb"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

type testcase struct {
	name        string
	method      string
	url         string
	params      []byte
	status      int
	contentType string
	expected    interface{} // nil, []byte OR func(t *testing.T, b []byte)
}

func init() {
	// Register an sql driver named "txdb".
	txdb.Register("txdb", "postgres", "postgres://")
}

func setup() (*api.Server, *pgxpool.Pool) {
	db, err := pgxpool.Connect(context.Background(), "postgres://")
	if err != nil {
		panic(err)
	}

	// Services
	userService := user.NewService(user.NewDatabaseRepository(db))

	// Logger
	// NewNop returns a no-op Logger. It never writes out logs or internal errors,
	// and it never runs user-defined hooks.
	logger := zap.NewNop().Sugar()

	// Rest server
	return api.NewServer(userService, logger), db
}

func run(t *testing.T, tt []testcase) {
	server, db := setup()
	defer db.Close()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			if tc.params != nil {
				buf = bytes.NewBuffer(tc.params)
			}

			req, err := http.NewRequest(tc.method, tc.url, buf)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			server.Router.ServeHTTP(rr, req)
			res := rr.Result()
			defer res.Body.Close()

			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode != tc.status {
				t.Errorf("Expected status %d, got %d. Response: %q", tc.status, res.StatusCode, string(body))
			}

			contentType := res.Header.Get("Content-Type")
			if !strings.Contains(contentType, tc.contentType) {
				t.Errorf("Expected content type %q, got %q", tc.contentType, contentType)
			}

			if tc.expected != nil {
				if expected, ok := tc.expected.([]byte); ok {
					if !bytes.Equal(expected, body) {
						t.Errorf("Expected \n%q\ngot\n%q", expected, body)
					}
				} else if fn, ok := tc.expected.(func(t *testing.T, b []byte)); ok {
					fn(t, body)
				} else {
					panic("Unsupported type")
				}
			}

		})
	}
}
