-- +goose Up
-- Documents table
CREATE TABLE documents
(
    id           SERIAL PRIMARY KEY,
    uuid         UUID                                NOT NULL,
    user_id      INT                                 NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    title        VARCHAR                             NOT NULL,
    description  VARCHAR                             NOT NULL,
    filename     VARCHAR                             NOT NULL,
    filetype     VARCHAR                             NOT NULL,
    content      VARCHAR,
    content_hash VARCHAR GENERATED ALWAYS AS (MD5(content)) STORED,

    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Documents history table
CREATE TABLE documents_history
(
    id              INT       NOT NULL,
    uuid            UUID      NOT NULL,
    user_id         INT       NOT NULL,
    created_at      TIMESTAMP NOT NULL,
    updated_at      TIMESTAMP NOT NULL,
    title           VARCHAR   NOT NULL,
    description     VARCHAR   NOT NULL,
    filename        VARCHAR   NOT NULL,
    filetype        VARCHAR   NOT NULL,
    content         VARCHAR,
    content_hash    VARCHAR,
    history_time    TIMESTAMP NOT NULL,
    history_user_id INT,
    operation       VARCHAR
);

-- Unique index on document & user
CREATE UNIQUE INDEX documents_per_user
    ON documents (id, user_id);

-- Index on document filename
CREATE INDEX documents_filename
    ON documents (filename);

-- Add trigger on documents to set updated_at on update
CREATE TRIGGER documents_set_updated_at
    BEFORE UPDATE
    ON
        documents
    FOR EACH ROW
EXECUTE PROCEDURE set_updated_at();

-- Function trigger for updating documents_history
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION process_documents_history() RETURNS TRIGGER AS
$$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        INSERT INTO documents_history SELECT OLD.*, CURRENT_TIMESTAMP, NULL, TG_OP;
    ELSE
        INSERT INTO documents_history SELECT NEW.*, NEW.updated_at, NULL, TG_OP;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- Apply trigger to documents
CREATE TRIGGER documents_history
    AFTER INSERT OR UPDATE OR DELETE
    ON documents
    FOR EACH ROW
EXECUTE FUNCTION process_documents_history();

-- +goose Down
DROP TRIGGER documents_history ON documents;
DROP FUNCTION process_documents_history;
DROP TRIGGER documents_set_updated_at ON documents;
DROP INDEX documents_per_user;
DROP INDEX documents_filename;
DROP TABLE documents_history;
DROP TABLE documents;
