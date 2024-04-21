-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS subjects_teachers (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    subject_id VARCHAR(255) NOT NULL,
    teacher_id VARCHAR(255) NOT NULL,
    CONSTRAINT fk_subject
    FOREIGN KEY (subject_id)
    REFERENCES subjects(id)
    ON DELETE CASCADE,
    CONSTRAINT fk_teacher
    FOREIGN KEY (teacher_id)
    REFERENCES users(uuid)
    ON DELETE CASCADE,
    UNIQUE (subject_id, teacher_id)
);

CREATE TABLE IF NOT EXISTS groups_subjects_new (
   group_id VARCHAR(255) NOT NULL,
   subject_teacher_id VARCHAR(255) NOT NULL,
   CONSTRAINT fk_group
       FOREIGN KEY (group_id)
           REFERENCES groups (id)
           ON DELETE CASCADE,
   CONSTRAINT fk_subject
       FOREIGN KEY (subject_teacher_id)
           REFERENCES subjects_teachers (id)
           ON DELETE CASCADE,
   PRIMARY KEY (group_id, subject_teacher_id),
   UNIQUE (group_id, subject_teacher_id)
);

INSERT INTO groups_subjects_new (group_id, subject_teacher_id)
SELECT gs.group_id, st.id
FROM groups_subjects gs
JOIN subjects_teachers st on gs.subject_id = st.id;

DROP TABLE groups_subjects;

ALTER TABLE groups_subjects_new
RENAME TO groups_subjects;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
CREATE TABLE IF NOT EXISTS groups_subjects_new (
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

INSERT INTO groups_subjects_new (group_id, subject_id)
SELECT gs.group_id, st.subject_id
FROM groups_subjects gs
         JOIN subjects_teachers st on gs.subject_id = st.id;

DROP TABLE groups_subjects;

ALTER TABLE groups_subjects_new
    RENAME TO groups_subjects;

DROP TABLE IF EXISTS subjects_teachers CASCADE;
-- +goose StatementEnd
