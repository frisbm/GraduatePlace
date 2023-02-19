-- name: CreateDocument :one
INSERT INTO documents (
    uuid, user_id, title, description, filename, filetype, content
) VALUES (
    gen_random_uuid(), $1, $2, $3, $4, $5, $6
)
RETURNING *;
