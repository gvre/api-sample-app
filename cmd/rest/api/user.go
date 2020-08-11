package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/gvre/api-sample-app/app"
)

// HandleGetUsers handles the "GET /users" endpoint.
func (s *Server) HandleGetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := s.Logger.With(actionKey, "get_users")

		ctx, cancel := context.WithTimeout(r.Context(), handlerDefaultTimeout)
		defer cancel()

		users, err := s.UserService.Users(ctx)
		switch err {
		case nil:
			logger.With("success", true).Info("")
			Ok(w, users, http.StatusOK)
		default:
			logger.With("success", false).Error(err)
			ServerError(w, err)
		}
	}
}

// HandleGetUser handles the "GET /users/{id}" endpoint.
func (s *Server) HandleGetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := s.Logger.With(actionKey, "get_user")
		vars := mux.Vars(r)

		id, err := toInt(vars, "id")
		if err != nil {
			logger.With("success", false).Error(err)
			BadRequestError(w, err)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), handlerDefaultTimeout)
		defer cancel()

		user, err := s.UserService.User(ctx, id)
		switch {
		case err == nil:
			logger.With("success", true).Info("")
			Ok(w, user, http.StatusOK)
		case app.IsNotFoundError(err):
			logger.With("success", false).Warn(err)
			NotFoundError(w, err)
		case app.IsValidationError(err):
			logger.With("success", false).Warn(err)
			ValidationError(w, err)
		default:
			logger.With("success", false).Error(err)
			ServerError(w, err)
		}
	}
}

// HandleAddUser handles the "POST /users" endpoint.
func (s *Server) HandleAddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := s.Logger.With(actionKey, "add_user")

		var fields struct {
			Name string `json:"name"`
		}

		if err := unmarshalBody(r, &fields); err != nil {
			// unmarshalBody always returns a "malformedRequestError" on failure.
			// Let the "Error" function handle it, in order to simplify the logic in controllers.
			// The http status code we pass to "Error" does not matter when a malformedRequestError is returned,
			// so pass the "httpStatusIgnore" constant to make that clear to the future reader of this code.
			// TODO improve this
			logger.With("success", false).Error(err)
			Error(w, err, httpStatusIgnore)
			return
		}
		name := fields.Name

		ctx, cancel := context.WithTimeout(r.Context(), handlerDefaultTimeout)
		defer cancel()

		id, err := s.UserService.Add(ctx, name)
		switch err {
		case nil:
			response := struct {
				ID int `json:"id"`
			}{ID: id}

			logger.With("success", true).Info("")
			Ok(w, response, http.StatusCreated)
		default:
			logger.With("success", false).Error(err)
			ServerError(w, err)
		}
	}
}
