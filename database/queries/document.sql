-- name: CreateDocument :one
INSERT INTO documents (
    uuid, user_id, title, description, filename, filetype, content
) VALUES (
    gen_random_uuid(), $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: SetDocumentHistoryUserId :one
UPDATE documents_history
SET history_user_id = $3
WHERE id=$1 AND history_time=$2
RETURNING *;

-- name: SetDocumentContent :one
UPDATE documents
SET content = $2
WHERE id=$1 RETURNING *;

-- name: GetDocument :one
SELECT * FROM documents
WHERE id=$1;

-- name: SearchDocuments :many
WITH matching_search_results AS (
    SELECT
        document_id,
        ts_rank_cd(ts, query, 32) AS document_rank
    FROM documents_search, websearch_to_tsquery('english', @query) query
    WHERE query @@ ts
    group by 1, 2
    ORDER BY document_rank DESC
    LIMIT 1000
)
SELECT documents.uuid,
       documents.created_at,
       documents.updated_at,
       documents.title,
       documents.description,
       documents.filename,
       documents.filetype,
       users.username,
       matching_search_results.document_rank AS rank,
       (SELECT COUNT('') FROM matching_search_results) AS count
FROM documents
JOIN matching_search_results ON documents.id = matching_search_results.document_id
JOIN users ON users.id = documents.user_id
ORDER BY rank DESC
LIMIT $1 OFFSET $2;
