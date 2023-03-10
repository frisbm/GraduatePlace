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
    id              SERIAL PRIMARY KEY,
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
--  HISTORY FIELDS
    document_id     INT NOT NULL,
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
        INSERT INTO documents_history (uuid, user_id, created_at, updated_at, title, description, filename, filetype,
                                       content, content_hash, document_id, history_time, history_user_id, operation)
        SELECT OLD.uuid,
               OLD.user_id,
               OLD.created_at,
               OLD.updated_at,
               OLD.title,
               OLD.description,
               OLD.filename,
               OLD.filetype,
               OLD.content,
               OLD.content_hash,
               OLD.id,
               CURRENT_TIMESTAMP,
               NULL,
               TG_OP;
    ELSE
        INSERT INTO documents_history (uuid, user_id, created_at, updated_at, title, description, filename, filetype,
                                       content, content_hash, document_id, history_time, history_user_id, operation)
        SELECT NEW.uuid,
               NEW.user_id,
               NEW.created_at,
               NEW.updated_at,
               NEW.title,
               NEW.description,
               NEW.filename,
               NEW.filetype,
               NEW.content,
               NEW.content_hash,
               NEW.id,
               NEW.updated_at,
               NULL,
               TG_OP;
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
