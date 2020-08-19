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
		"/users/{id}",
		s.HandleGetUser(),
	).Methods("GET")
}
