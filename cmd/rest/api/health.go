package api

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gvre/api-sample-app/app"
)

// HandleCheckLive is used for checking if the service is up.
func (s *Server) HandleCheckLive() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := s.Logger.With(actionKey, "check_live")
		logger.With("success", true).Info("")
		Ok(w, nil, http.StatusOK)
	}
}

// HandleCheckHealth is used for checking if the service is healthy.
func (s *Server) HandleCheckHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := s.Logger.With(actionKey, "check_ready")

		ctx, cancel := context.WithTimeout(r.Context(), handlerDefaultTimeout)
		defer cancel()

		// If more checks are added, they should be called concurrently
		// to reduce the total execution time.
		var checks []app.Health
		checks = append(checks, s.UserService.Health(ctx))

		status := app.HealthStatusOK
		for _, h := range checks {
			// Only core services affect the healthcheck status.
			if !h.Core {
				continue
			}

			switch h.Status {
			case app.HealthStatusError:
				status = app.HealthStatusError
			case app.HealthStatusWarning:
				if status != app.HealthStatusError {
					status = app.HealthStatusWarning
				}
			}
		}

		success := status == app.HealthStatusOK
		logger.With("success", success).Info("")

		httpStatus := http.StatusOK
		if status == app.HealthStatusError {
			httpStatus = http.StatusServiceUnavailable
		}

		hostname, _ := os.Hostname()
		result := struct {
			Datetime string       `json:"datetime"`
			Status   string       `json:"status"`
			Hostname string       `json:"hostname"`
			Checks   []app.Health `json:"checks"`
		}{
			Datetime: time.Now().UTC().Format("2006-01-02T15:04:05Z"),
			Status:   status,
			Hostname: hostname,
			Checks:   checks,
		}

		Ok(w, result, httpStatus)
	}
}
