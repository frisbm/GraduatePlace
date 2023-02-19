package document

import (
	"context"

	"github.com/MatthewFrisby/thesis-pieces/pkg/models/document"

	"github.com/MatthewFrisby/thesis-pieces/pkg/store"
)

type Store struct {
	db store.Querier
}

func NewStore(db store.Querier) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) CreateDocument(ctx context.Context, uploadDocument document.UploadDocument) (*store.Document, error) {
	return s.db.CreateDocument(ctx, store.CreateDocumentParams{
		UserID:      uploadDocument.UserID,
		Title:       uploadDocument.Title,
		Description: uploadDocument.Description,
		Filename:    uploadDocument.FileName,
		Filetype:    uploadDocument.FileType,
		Content:     uploadDocument.Content,
	})
}
