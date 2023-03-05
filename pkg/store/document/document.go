package document

import (
	"context"

	"github.com/frisbm/graduateplace/pkg/models/document"

	"github.com/frisbm/graduateplace/pkg/store"
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
	doc, err := s.db.CreateDocument(ctx, store.CreateDocumentParams{
		UserID:      uploadDocument.UserID,
		Title:       uploadDocument.Title,
		Description: uploadDocument.Description,
		Filename:    uploadDocument.FileName,
		Filetype:    uploadDocument.FileType,
		Content:     uploadDocument.Content,
	})
	if err != nil {
		return nil, err
	}
	_, err = s.db.SetDocumentHistoryUserId(ctx, store.SetDocumentHistoryUserIdParams{
		DocumentID:    doc.ID,
		HistoryTime:   doc.UpdatedAt,
		HistoryUserID: &uploadDocument.UserID,
	})
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (s *Store) SetDocumentContent(ctx context.Context, id int32, content *string) (*store.Document, error) {
	doc, err := s.db.SetDocumentContent(ctx, store.SetDocumentContentParams{
		ID:      id,
		Content: content,
	})
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (s *Store) GetDocument(ctx context.Context, id int32) (*store.Document, error) {
	doc, err := s.db.GetDocument(ctx, id)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (s *Store) SearchDocuments(ctx context.Context, searchDocuments document.SearchDocuments) ([]*store.SearchDocumentsRow, error) {
	docs, err := s.db.SearchDocuments(ctx, store.SearchDocumentsParams{
		Limit:              searchDocuments.Limit,
		Offset:             searchDocuments.Offset,
		WebsearchToTsquery: searchDocuments.Query,
	})
	if err != nil {
		return nil, err
	}
	return docs, nil
}
