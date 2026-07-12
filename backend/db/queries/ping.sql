-- name: Ping :one
-- Smoke-test query with no table dependency: proves a real round trip
-- through pgx to Postgres and back, via sqlc-generated code, without
-- assuming any business schema exists yet.
SELECT 1 AS ok;
