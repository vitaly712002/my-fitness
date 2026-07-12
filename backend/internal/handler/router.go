package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/vitaly712002/my-fitness/backend/internal/repository"
)

// NewRouter wires up every HTTP route the backend exposes. Called once from
// main.go; new route groups (auth, exercises, workouts, ...) get added here
// as their own Mount/Route calls as later releases add them.
func NewRouter(queries *repository.Queries, logger *slog.Logger) http.Handler {
	r := chi.NewRouter()

	// RequestID: tags each request with a unique ID (visible in logs below),
	// useful once traffic volume makes "which request logged this line?" a
	// real question.
	r.Use(middleware.RequestID)
	// Recoverer: turns a panic in any handler into a 500 instead of killing
	// the whole process — one broken request shouldn't take down the server.
	r.Use(middleware.Recoverer)
	r.Use(requestLogger(logger))

	r.Get("/health", Health(queries, logger))

	return r
}

// requestLogger replaces chi's built-in middleware.Logger (which writes
// plain text to stdout) with one that logs through slog, so every log line
// in the app — from routing down to business logic — shares the same
// structured (JSON in prod) format.
func requestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			logger.Info("http request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"duration_ms", time.Since(start).Milliseconds(),
				"request_id", middleware.GetReqID(r.Context()),
			)
		})
	}
}
