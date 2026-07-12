# Shortcuts for the recurring dev commands so nobody has to remember the
# full goose/sqlc invocations (in particular the Postgres connection string).
# Assumes `goose`/`sqlc` are on PATH (see docs/backend/01-getting-started.md)
# and this is run from inside backend/.

# Pulls POSTGRES_* values from the repo-root .env — same file docker-compose
# uses, so credentials only live in one place.
include ../.env

DB_URL := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

.PHONY: run build migrate-up migrate-down migrate-status sqlc-generate

run: ## Start the API server (native, hot-reloads only via manual restart).
	go run ./cmd/api

build: ## Compile a binary into bin/api (bin/ is gitignored).
	go build -o bin/api ./cmd/api

migrate-up: ## Apply all pending migrations.
	goose -dir db/migrations postgres "$(DB_URL)" up

migrate-down: ## Roll back the most recently applied migration.
	goose -dir db/migrations postgres "$(DB_URL)" down

migrate-status: ## Show which migrations are applied vs pending.
	goose -dir db/migrations postgres "$(DB_URL)" status

sqlc-generate: ## Regenerate internal/repository from db/queries + db/migrations.
	sqlc generate
