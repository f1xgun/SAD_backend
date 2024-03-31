-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS subjects (
    id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS groups_subjects (
    group_id VARCHAR(255) NOT NULL,
    subject_id VARCHAR(255) NOT NULL,
    CONSTRAINT fk_group
    FOREIGN KEY (group_id)
    REFERENCES groups (id)
    ON DELETE CASCADE,
    CONSTRAINT fk_subject
    FOREIGN KEY (subject_id)
    REFERENCES subjects (id)
    ON DELETE CASCADE,
    PRIMARY KEY (group_id, subject_id),
    UNIQUE (group_id, subject_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS groups_subjects;

DROP TABLE IF EXISTS subjects;
-- +goose StatementEnd
