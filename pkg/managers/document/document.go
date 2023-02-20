package document

import (
	"context"

	"github.com/MatthewFrisby/thesis-pieces/pkg/store"

	"github.com/MatthewFrisby/thesis-pieces/pkg/models/document"
	"github.com/MatthewFrisby/thesis-pieces/pkg/utils"

	"github.com/pkg/errors"
)

type Store interface {
	CreateDocument(ctx context.Context, uploadDocument document.UploadDocument) (*store.Document, error)
}

type S3 interface {
	UploadFile(ctx context.Context, bucketName, filename string, file []byte) error
}

type Manager struct {
	store Store
	s3    S3
}

func NewManager(store Store, s3 S3) *Manager {
	return &Manager{
		store: store,
		s3:    s3,
	}
}

func (m *Manager) UploadDocument(ctx context.Context, uploadDocument document.UploadDocument) error {
	userCtx, err := utils.GetUserFromContext(ctx)
	if err != nil {
		return err
	}
	uploadDocument.UserID = userCtx.ID

	if len(uploadDocument.File) == 0 {
		return errors.New("file bytes not found")
	}
	err = m.s3.UploadFile(ctx, userCtx.Username, uploadDocument.FileName, uploadDocument.File)
	if err != nil {
		return err
	}
	_, err = m.store.CreateDocument(ctx, uploadDocument)

	return err
}