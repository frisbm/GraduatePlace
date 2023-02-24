-- +goose Up
-- Document search table, pulled out from documents table, felt wise to split
CREATE TABLE documents_search
(
    id          SERIAL PRIMARY KEY,
    uuid        UUID                                NOT NULL,
    document_id INT                                 NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    ts          TSVECTOR,

    CONSTRAINT fk_document FOREIGN KEY (document_id) REFERENCES documents (id) ON DELETE CASCADE
);

-- GIN index on ts for full text search
CREATE INDEX documents_search_ts
    ON documents_search USING GIN (ts);

-- Add trigger on documents_search to set updated_at on update
CREATE TRIGGER documents_search_set_updated_at
    BEFORE UPDATE
    ON
        documents_search
    FOR EACH ROW
EXECUTE PROCEDURE set_updated_at();

-- Function trigger adding/updating document search row when a document row is updated
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION create_documents_search() RETURNS TRIGGER AS
$$
DECLARE
    document_id_exists INT;
BEGIN
    IF (NEW.content IS NULL OR OLD.content_hash <> NEW.content_hash) THEN
        RETURN NULL;
    END IF;
    SELECT document_id INTO document_id_exists FROM documents_search WHERE documents_search.document_id = NEW.id;
    IF NOT FOUND THEN
        INSERT INTO documents_search (uuid, document_id, ts)
        VALUES (gen_random_uuid(), NEW.id, to_tsvector(
                'english',
                NEW.title || ' ' ||
                NEW.description || ' ' ||
                COALESCE(NEW.content, '')
            ));
    ELSE
        UPDATE documents_search
        SET ts = to_tsvector(
                'english',
                NEW.title || ' ' ||
                NEW.description || ' ' ||
                COALESCE(NEW.content, ''))
        WHERE documents_search.document_id = NEW.id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- Apply trigger to documents
CREATE TRIGGER documents_search
    AFTER INSERT OR UPDATE
    ON documents
    FOR EACH ROW
EXECUTE FUNCTION create_documents_search();

-- +goose Down
DROP TRIGGER documents_search ON documents;
DROP FUNCTION create_documents_search;
DROP TRIGGER documents_search_set_updated_at ON documents_search;
DROP INDEX documents_search_ts;
DROP TABLE documents_search;
