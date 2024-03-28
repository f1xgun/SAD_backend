-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS users_groups (
    user_id VARCHAR(255) NOT NULL,
    group_id VARCHAR(255) NOT NULL,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users (uuid)
        ON DELETE CASCADE,
    CONSTRAINT fk_group
        FOREIGN KEY (group_id)
        REFERENCES groups (id)
        ON DELETE CASCADE,
    PRIMARY KEY (user_id, group_id),
    UNIQUE (user_id, group_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS users_groups;
-- +goose StatementEnd
