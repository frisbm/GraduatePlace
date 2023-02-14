-- name: GetUserFromEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserFromUUID :one
SELECT * FROM users
WHERE uuid = $1 LIMIT 1;

-- name: GetUsers :many
SELECT * FROM users;

-- name: CreateUser :one
INSERT INTO users (
    username, email, password, first_name, last_name, is_admin, uuid
) VALUES (
    $1, $2, $3, $4, $5, FALSE, gen_random_uuid()
)
RETURNING *;
