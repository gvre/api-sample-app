package api

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// CreateRoutes creates all app routes.
func (s *Server) CreateRoutes() {
	s.Router.Handle(
		"GET /-/live",
		s.EnableCORS(s.HandleCheckLive()),
	)

	s.Router.Handle(
		"GET /-/health",
		s.EnableCORS(s.HandleCheckHealth()),
	)

	s.Router.Handle(
		"GET /-/metrics",
		s.EnableCORS(promhttp.Handler()),
	)

	s.Router.Handle(
		"GET /users",
		s.EnableCORS(s.HandleGetUsers()),
	)

	s.Router.Handle(
		"POST /users",
		s.EnableCORS(s.LimitBody(s.HandleAddUser())),
	)

	s.Router.Handle(
		"GET /users/{id}",
		s.EnableCORS(s.HandleGetUser()),
	)
}
