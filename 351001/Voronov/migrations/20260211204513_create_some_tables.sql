-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS users (
    user_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    login varchar(64) NOT NULL UNIQUE,
    password varchar(128) NOT NULL,
    first_name varchar(64),
    last_name varchar(64)
);

CREATE TABLE IF NOT EXISTS issue (
    issue_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id int REFERENCES users(user_id) ON DELETE CASCADE,
    title varchar(64),
    content text,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS reaction (
    reaction_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    issue_id INT REFERENCES issue(issue_id) ON DELETE CASCADE,
    content text
);

CREATE TABLE IF NOT EXISTS label (
    label_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name varchar(32) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS issue_label (
    issue_id INT REFERENCES issue(issue_id),
    label_id INT REFERENCES label(label_id),
    PRIMARY KEY (issue_id, label_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users, reaction, issue, label, issue_label;
-- +goose StatementEnd
