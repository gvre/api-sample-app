package api

// CreateRoutes creates all app routes.
func (s *Server) CreateRoutes() {
	s.Router.Handle(
		"/users",
		s.Compress(s.HandleGetUsers()),
	).Methods("GET")

	s.Router.Handle(
		"/users",
		s.LimitBody(s.HandleAddUser()),
	).Methods("POST")

	s.Router.Handle(
		"/users/{id:[0-9]+}",
		s.HandleGetUser(),
	).Methods("GET")
}
