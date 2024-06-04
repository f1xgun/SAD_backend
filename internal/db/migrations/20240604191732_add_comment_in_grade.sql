-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE IF EXISTS grades
ADD comment VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE IF EXISTS grades
DROP COLUMN comment;
-- +goose StatementEnd
