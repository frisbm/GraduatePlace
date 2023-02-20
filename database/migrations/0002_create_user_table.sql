-- +goose Up
CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    uuid       UUID                                NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    username   VARCHAR                             NOT NULL,
    email      VARCHAR                             NOT NULL,
    password   VARCHAR                             NOT NULL,
    first_name VARCHAR                             NOT NULL,
    last_name  VARCHAR                             NOT NULL,
    is_admin   BOOLEAN   DEFAULT FALSE             NOT NULL
);

CREATE TABLE users_history
(
    id              INT                   NOT NULL,
    uuid            UUID                  NOT NULL,
    created_at      TIMESTAMP             NOT NULL,
    updated_at      TIMESTAMP             NOT NULL,
    username        VARCHAR               NOT NULL,
    email           VARCHAR               NOT NULL,
    password        VARCHAR               NOT NULL,
    first_name      VARCHAR               NOT NULL,
    last_name       VARCHAR               NOT NULL,
    is_admin        BOOLEAN NOT NULL,
    history_time    TIMESTAMP             NOT NULL,
    history_user_id INT,
    operation       VARCHAR
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

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION process_users_history() RETURNS TRIGGER AS
$$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        INSERT INTO users_history SELECT OLD.*, CURRENT_TIMESTAMP, NULL, TG_OP;
    ELSE
        INSERT INTO users_history SELECT NEW.*, NEW.updated_at, NULL, TG_OP;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER users_history
    AFTER INSERT OR UPDATE OR DELETE
    ON users
    FOR EACH ROW
EXECUTE FUNCTION process_users_history();

-- +goose Down
DROP TRIGGER users_history ON users;
DROP FUNCTION process_users_history;
DROP TRIGGER user_set_updated_at ON users;
DROP INDEX users_email_key;
DROP INDEX users_username_key;
DROP TABLE users_history;
DROP TABLE users;
