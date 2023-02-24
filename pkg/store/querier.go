// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0

package store

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateDocument(ctx context.Context, arg CreateDocumentParams) (*Document, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (*User, error)
	GetDocument(ctx context.Context, id int32) (*Document, error)
	GetUserFromEmail(ctx context.Context, email string) (*User, error)
	GetUserFromUUID(ctx context.Context, uuid uuid.UUID) (*User, error)
	GetUsers(ctx context.Context) ([]*User, error)
	SearchDocuments(ctx context.Context, arg SearchDocumentsParams) ([]*SearchDocumentsRow, error)
	SetDocumentContent(ctx context.Context, arg SetDocumentContentParams) (*Document, error)
	SetDocumentHistoryUserId(ctx context.Context, arg SetDocumentHistoryUserIdParams) (*DocumentsHistory, error)
	SetUserHistoryUserId(ctx context.Context, arg SetUserHistoryUserIdParams) (*UsersHistory, error)
}

var _ Querier = (*Queries)(nil)
