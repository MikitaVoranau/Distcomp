-- +goose Up
-- +goose StatementBegin

-- 1. Создаем схему (если её нет) и настраиваем пути
CREATE SCHEMA IF NOT EXISTS distcomp;
SET search_path TO distcomp, public;

-- 2. Таблица пользователей
CREATE TABLE IF NOT EXISTS tbl_user (
                                        id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                                        login VARCHAR(64) NOT NULL UNIQUE,
                                        password VARCHAR(128) NOT NULL,
                                        firstname VARCHAR(64) NOT NULL,
                                        lastname VARCHAR(64) NOT NULL
);

-- 3. Таблица задач (СВЯЗЬ С CASCADE)
CREATE TABLE IF NOT EXISTS tbl_issue (
                                         id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                                         user_id BIGINT NOT NULL REFERENCES tbl_user(id) ON DELETE CASCADE,
                                         title VARCHAR(64) NOT NULL,
                                         content VARCHAR(2048) NOT NULL,
                                         created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                         modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 4. Таблица меток
CREATE TABLE IF NOT EXISTS tbl_label (
                                         id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                                         name VARCHAR(32) NOT NULL UNIQUE
);

-- 5. Таблица реакций (СВЯЗЬ С CASCADE)
CREATE TABLE IF NOT EXISTS tbl_reaction (
                                            id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                                            issue_id BIGINT NOT NULL REFERENCES tbl_issue(id) ON DELETE CASCADE,
                                            content VARCHAR(2048) NOT NULL
);

-- 6. Таблица связей (Many-to-Many)
CREATE TABLE IF NOT EXISTS tbl_issue_label (
                                               issue_id BIGINT NOT NULL REFERENCES tbl_issue(id) ON DELETE CASCADE,
                                               label_id BIGINT NOT NULL REFERENCES tbl_label(id) ON DELETE CASCADE,
                                               PRIMARY KEY (issue_id, label_id)
);

-- 7. Настройка для тестов
ALTER ROLE postgres SET search_path TO distcomp, public;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_issue_label;
DROP TABLE IF EXISTS tbl_reaction;
DROP TABLE IF EXISTS tbl_label;
DROP TABLE IF EXISTS tbl_issue;
DROP TABLE IF EXISTS tbl_user;
-- +goose StatementEnd