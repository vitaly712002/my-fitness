// Package config loads application configuration from environment variables.
package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

// Config holds every value the backend needs at startup. Kept as one flat
// struct so main.go (and later, other packages) has a single place to read
// settings from instead of calling os.Getenv scattered across the codebase.
type Config struct {
	ServerPort string // e.g. "8080" — net/http wants ":8080", callers add the colon themselves
	LogLevel   string // "debug" | "info" | "warn" | "error"

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
}

// Load reads .env (if present) into the process environment, then builds a
// Config from environment variables, filling in defaults where that's safe.
//
// We try two .env locations because `go run ./cmd/api` gets invoked either
// from backend/ (the common case) or from the repo root, depending on which
// terminal tab is open — both should work without extra flags.
func Load() Config {
	for _, path := range []string{".env", "../.env"} {
		if err := godotenv.Load(path); err == nil {
			break
		}
	}

	return Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		LogLevel:   getEnv("LOG_LEVEL", "info"),

		PostgresHost: getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort: getEnv("POSTGRES_PORT", "5432"),
		PostgresUser: getEnv("POSTGRES_USER", "fitness"),
		// No fallback for the password: a secret shouldn't have a guessable default.
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:       getEnv("POSTGRES_DB", "fitness"),
	}
}

// DatabaseURL builds the libpq-style connection string that both pgx (at
// runtime) and goose (from the CLI) understand.
//
// sslmode=disable is fine here because Postgres only runs on localhost in
// dev (via docker-compose); a managed prod database will need a different
// value, set through its own env at deploy time.
func (c Config) DatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.PostgresUser, c.PostgresPassword, c.PostgresHost, c.PostgresPort, c.PostgresDB,
	)
}

// ParseLogLevel converts LogLevel into an slog.Level, defaulting to Info for
// anything unrecognised rather than failing startup over a typo in .env.
func (c Config) ParseLogLevel() slog.Level {
	switch c.LogLevel {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
