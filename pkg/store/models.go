// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0

package store

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string
	Email     string
	Password  string
	FirstName string
	LastName  string
	IsAdmin   bool
	Uuid      uuid.UUID
}
