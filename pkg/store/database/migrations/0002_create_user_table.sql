-- +goose Up
-- Users table
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

-- Users history table
CREATE TABLE users_history
(
    id              SERIAL PRIMARY KEY,
    uuid            UUID      NOT NULL,
    created_at      TIMESTAMP NOT NULL,
    updated_at      TIMESTAMP NOT NULL,
    username        VARCHAR   NOT NULL,
    email           VARCHAR   NOT NULL,
    password        VARCHAR   NOT NULL,
    first_name      VARCHAR   NOT NULL,
    last_name       VARCHAR   NOT NULL,
    is_admin        BOOLEAN   NOT NULL,
--  HISTORY FIELDS
    user_id         INT NOT NULL,
    history_time    TIMESTAMP NOT NULL,
    history_user_id INT,
    operation       VARCHAR
);

-- Unique index on username
CREATE UNIQUE INDEX users_username_key
    ON users (username);

-- Unique index on email
CREATE UNIQUE INDEX users_email_key
    ON users (email);

-- Add trigger on users to set updated_at on update
CREATE TRIGGER user_set_updated_at
    BEFORE UPDATE
    ON
        users
    FOR EACH ROW
EXECUTE PROCEDURE set_updated_at();

-- Function trigger for updating users_history
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION process_users_history() RETURNS TRIGGER AS
$$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        INSERT INTO users_history (uuid, created_at, updated_at, username, email, password, first_name, last_name,
                                   is_admin, user_id, history_time, history_user_id, operation)
        SELECT OLD.uuid,
               OLD.created_at,
               OLD.updated_at,
               OLD.username,
               OLD.email,
               OLD.password,
               OLD.first_name,
               OLD.last_name,
               OLD.is_admin,
               OLD.id,
               CURRENT_TIMESTAMP,
               NULL,
               TG_OP;
    ELSE
        INSERT INTO users_history (uuid, created_at, updated_at, username, email, password, first_name, last_name,
                                   is_admin, user_id, history_time, history_user_id, operation)
        SELECT NEW.uuid,
               NEW.created_at,
               NEW.updated_at,
               NEW.username,
               NEW.email,
               NEW.password,
               NEW.first_name,
               NEW.last_name,
               NEW.is_admin,
               NEW.id,
               NEW.updated_at,
               NULL,
               TG_OP;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- Apply trigger to users
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
