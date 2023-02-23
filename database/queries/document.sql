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
