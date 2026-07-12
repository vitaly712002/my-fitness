package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/vitaly712002/my-fitness/backend/internal/repository"
)

// healthResponse is the JSON body /health returns. "db" reflects whether the
// Ping query round-tripped through Postgres successfully — this is what
// makes /health a real end-to-end check rather than just "the process is up".
type healthResponse struct {
	Status string `json:"status"`
	DB     string `json:"db"`
}

// Health returns a handler for GET /health. It depends on *repository.Queries
// (not a raw pgxpool.Pool) specifically so the check exercises the same
// sqlc-generated code path every other feature will use.
func Health(queries *repository.Queries, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := healthResponse{Status: "ok", DB: "up"}
		statusCode := http.StatusOK

		if _, err := queries.Ping(r.Context()); err != nil {
			logger.Error("health check: db ping failed", "error", err)
			resp.Status = "degraded"
			resp.DB = "down"
			statusCode = http.StatusServiceUnavailable
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(resp)
	}
}
