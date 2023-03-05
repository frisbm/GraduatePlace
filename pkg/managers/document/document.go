package document

import (
	"context"

	"github.com/frisbm/graduateplace/pkg/models/pagination"

	"github.com/frisbm/graduateplace/pkg/store"

	"github.com/frisbm/graduateplace/pkg/models/document"
	"github.com/frisbm/graduateplace/pkg/utils"

	"github.com/pkg/errors"
)

type Store interface {
	CreateDocument(ctx context.Context, uploadDocument document.UploadDocument) (*store.Document, error)
	SearchDocuments(ctx context.Context, searchDocuments document.SearchDocuments) ([]*store.SearchDocumentsRow, error)
}

type S3 interface {
	UploadFile(ctx context.Context, bucketName, filename string, file []byte) error
}

type Tasks interface {
	ProcessDocumentTask(documentId int32, bucket string) error
}

type Manager struct {
	store Store
	s3    S3
	tasks Tasks
}

func NewManager(store Store, s3 S3, tasks Tasks) *Manager {
	return &Manager{
		store: store,
		s3:    s3,
		tasks: tasks,
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

	doc, err := m.store.CreateDocument(ctx, uploadDocument)
	if err != nil {
		return err
	}

	return m.tasks.ProcessDocumentTask(doc.ID, userCtx.Username)
}

func (m *Manager) SearchDocuments(ctx context.Context, searchDocuments document.SearchDocuments) (*pagination.Pagination[document.SearchDocumentsResult], error) {
	results, err := m.store.SearchDocuments(ctx, searchDocuments)
	if err != nil {
		return nil, err
	}

	data := make([]*document.SearchDocumentsResult, len(results))
	for i, result := range results {
		resultData := document.SearchDocumentsResult{
			Uuid:        result.Uuid,
			CreatedAt:   result.CreatedAt,
			UpdatedAt:   result.UpdatedAt,
			Title:       result.Title,
			Description: result.Description,
			Filename:    result.Filename,
			Filetype:    result.Filetype,
			Username:    result.Username,
			Rank:        result.Rank,
		}
		data[i] = &resultData
	}
	searchResults := pagination.Pagination[document.SearchDocumentsResult]{
		Data:   data,
		Limit:  searchDocuments.Limit,
		Offset: searchDocuments.Offset,
		Count:  0,
	}

	if len(results) > 0 {
		searchResults.Count = int32(results[0].Count)
	}
	return &searchResults, nil
}
