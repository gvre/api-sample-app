package api

import (
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/gvre/api-sample-app/user"
)

// handlerDefaultTimeout is the timeout the handlers pass to the inner layers.
const handlerDefaultTimeout = 5 * time.Second

// The Server is used as a container for the most important dependencies.
type Server struct {
	Router      *mux.Router
	UserService *user.Service
	Logger      *zap.SugaredLogger
}

// NewServer returns a pointer to a new Server.
func NewServer(userService *user.Service, logger *zap.SugaredLogger) *Server {
	server := &Server{
		Router:      mux.NewRouter().StrictSlash(true),
		UserService: userService,
		Logger:      logger,
	}
	server.CreateRoutes()

	return server
}
