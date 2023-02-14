-- +goose Up
CREATE TABLE users
(
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    username   VARCHAR                             NOT NULL,
    email      VARCHAR                             NOT NULL,
    password   VARCHAR                             NOT NULL,
    first_name VARCHAR                             NOT NULL,
    last_name  VARCHAR                             NOT NULL,
    is_admin   BOOLEAN   DEFAULT FALSE             NOT NULL,
    uuid       UUID                                NOT NULL
);

CREATE UNIQUE INDEX users_username_key
    ON users (username);

CREATE UNIQUE INDEX users_email_key
    ON users (email);

CREATE TRIGGER user_set_updated_at
    BEFORE UPDATE
    ON
        users
    FOR EACH ROW
    EXECUTE PROCEDURE set_updated_at();


-- +goose Down
DROP TRIGGER user_set_updated_at ON users;
DROP INDEX users_email_key;
DROP INDEX users_username_key;
DROP TABLE users;
