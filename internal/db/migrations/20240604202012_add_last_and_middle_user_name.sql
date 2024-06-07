-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE IF EXISTS users
ADD COLUMN last_name VARCHAR(255) DEFAULT 'UNKNOWN' NOT NULL,
ADD COLUMN middle_name VARCHAR(255);
ALTER TABLE IF EXISTS users
ALTER COLUMN last_name DROP DEFAULT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE IF EXISTS users
DROP COLUMN last_name,
DROP COLUMN middle_name;
-- +goose StatementEnd
