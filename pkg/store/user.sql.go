// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0
// source: user.sql

package store

import (
	"context"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    username, email, password, first_name, last_name, is_admin, uuid
) VALUES (
    $1, $2, $3, $4, $5, FALSE, gen_random_uuid()
)
RETURNING id, uuid, created_at, updated_at, username, email, password, first_name, last_name, is_admin
`

type CreateUserParams struct {
	Username  string
	Email     string
	Password  string
	FirstName string
	LastName  string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (*User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Username,
		arg.Email,
		arg.Password,
		arg.FirstName,
		arg.LastName,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.FirstName,
		&i.LastName,
		&i.IsAdmin,
	)
	return &i, err
}

const getUserFromEmail = `-- name: GetUserFromEmail :one
SELECT id, uuid, created_at, updated_at, username, email, password, first_name, last_name, is_admin FROM users
WHERE email = $1 LIMIT 1
`

func (q *Queries) GetUserFromEmail(ctx context.Context, email string) (*User, error) {
	row := q.db.QueryRowContext(ctx, getUserFromEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.FirstName,
		&i.LastName,
		&i.IsAdmin,
	)
	return &i, err
}

const getUserFromUUID = `-- name: GetUserFromUUID :one
SELECT id, uuid, created_at, updated_at, username, email, password, first_name, last_name, is_admin FROM users
WHERE uuid = $1 LIMIT 1
`

func (q *Queries) GetUserFromUUID(ctx context.Context, uuid uuid.UUID) (*User, error) {
	row := q.db.QueryRowContext(ctx, getUserFromUUID, uuid)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.FirstName,
		&i.LastName,
		&i.IsAdmin,
	)
	return &i, err
}

const getUsers = `-- name: GetUsers :many
SELECT id, uuid, created_at, updated_at, username, email, password, first_name, last_name, is_admin FROM users
`

func (q *Queries) GetUsers(ctx context.Context) ([]*User, error) {
	rows, err := q.db.QueryContext(ctx, getUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Uuid,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Username,
			&i.Email,
			&i.Password,
			&i.FirstName,
			&i.LastName,
			&i.IsAdmin,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
