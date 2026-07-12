-- +goose Up
-- pgcrypto gives us gen_random_uuid(), which Release 1's users table (and
-- most tables after it) will use as the default for its id column instead of
-- relying on the application to generate UUIDs.
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- +goose Down
DROP EXTENSION IF EXISTS pgcrypto;
