package api

import (
	"github.com/gorilla/mux"
	"log/slog"
	"time"

	"github.com/gvre/api-sample-app/user"
)

// handlerDefaultTimeout is the timeout the handlers pass to the inner layers.
const handlerDefaultTimeout = 5 * time.Second

// The Server is used as a container for the most important dependencies.
type Server struct {
	Router      *mux.Router
	UserService *user.Service
	Logger      *slog.Logger
}

// NewServer returns a pointer to a new Server.
func NewServer(userService *user.Service, logger *slog.Logger) *Server {
	server := &Server{
		Router:      mux.NewRouter().StrictSlash(true),
		UserService: userService,
		Logger:      logger,
	}
	server.CreateRoutes()

	return server
}
