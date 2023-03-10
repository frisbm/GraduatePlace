-- name: GetUserFromEmail :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;

-- name: GetUserFromUUID :one
SELECT *
FROM users
WHERE uuid = $1
LIMIT 1;

-- name: GetUsers :many
SELECT *
FROM users;

-- name: CreateUser :one
INSERT INTO users (uuid, username, email, password, first_name, last_name, is_admin)
VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, FALSE)
RETURNING *;

-- name: SetUserHistoryUserId :one
UPDATE users_history
SET history_user_id = $3
WHERE user_id = $1
  AND history_time = $2
RETURNING *;
