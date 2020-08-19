package api

// CreateRoutes creates all app routes.
func (s *Server) CreateRoutes() {
	s.Router.Handle(
		"/users",
		s.EnableCORS(s.Compress(s.HandleGetUsers())),
	).Methods("GET")

	s.Router.Handle(
		"/users",
		s.EnableCORS(s.LimitBody(s.HandleAddUser())),
	).Methods("POST")

	s.Router.Handle(
		"/users/{id}",
		s.EnableCORS(s.HandleGetUser()),
	).Methods("GET")
}
