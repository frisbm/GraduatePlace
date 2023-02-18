-- +goose Up
CREATE TABLE documents
(
    id          SERIAL PRIMARY KEY,
    uuid        UUID                                NOT NULL,
    user_id     INT NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    title       VARCHAR                             NOT NULL,
    description VARCHAR                             NOT NULL,
    filepath    VARCHAR                             NOT NULL,
    filetype    VARCHAR                             NOT NULL,
    content     VARCHAR,

    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,

    ts TSVECTOR GENERATED ALWAYS AS (
        to_tsvector('english', title || ' ' || description || ' ' ||  content)
    ) STORED
);

CREATE UNIQUE INDEX documents_per_user
    ON documents (id, user_id);

CREATE UNIQUE INDEX documents_file_path
    ON documents (filepath);

CREATE INDEX documents_ts
    ON documents USING GIN (ts);

CREATE TRIGGER documents_set_updated_at
    BEFORE UPDATE
    ON
        documents
    FOR EACH ROW
EXECUTE PROCEDURE set_updated_at();

-- +goose Down
DROP TRIGGER documents_set_updated_at ON documents;
DROP INDEX documents_ts;
DROP INDEX documents_per_user;
DROP INDEX documents_file_path;
DROP TABLE documents;
