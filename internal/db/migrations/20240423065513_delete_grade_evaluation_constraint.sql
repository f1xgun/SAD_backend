-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE IF EXISTS grades
DROP CONSTRAINT grades_evaluation_check;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE IF EXISTS grades
ADD CONSTRAINT grades_evaluation_check CHECK (evaluation >= 0 AND evaluation <= 5);
-- +goose StatementEnd
