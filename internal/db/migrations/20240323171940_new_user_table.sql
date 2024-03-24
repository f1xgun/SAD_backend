-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TYPE user_role as ENUM('student', 'teacher', 'admin');

CREATE TABLE IF NOT EXISTS users (
    uuid VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    login VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'student'
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS user_role;
-- +goose StatementEnd
