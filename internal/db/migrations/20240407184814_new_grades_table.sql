-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS grades (
      id VARCHAR(255) NOT NULL UNIQUE,
      evaluation INT NOT NULL CHECK (evaluation >= 0 AND evaluation <= 5),
      created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
      subject_id VARCHAR(255) NOT NULL,
      student_id VARCHAR(255) NOT NULL,
      CONSTRAINT fk_subject
      FOREIGN KEY (subject_id)
      REFERENCES subjects(id)
      ON DELETE CASCADE,
      CONSTRAINT fk_student
      FOREIGN KEY (student_id)
      REFERENCES users(uuid)
      ON DELETE CASCADE,
      PRIMARY KEY (id, subject_id, student_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE IF EXISTS grades;
-- +goose StatementEnd
