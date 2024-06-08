-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS grades_teachers (
   grade_id VARCHAR(255) NOT NULL,
   teacher_id VARCHAR(255) NOT NULL,
   CONSTRAINT fk_grade
   FOREIGN KEY (grade_id)
   REFERENCES grades(id)
   ON DELETE CASCADE,
   CONSTRAINT fk_teacher
   FOREIGN KEY (teacher_id)
   REFERENCES users(uuid)
   ON DELETE CASCADE,
   UNIQUE (grade_id, teacher_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS grades_teachers;
-- +goose StatementEnd
