package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hibiken/asynq"
	"github.com/pkg/errors"

	"github.com/MatthewFrisby/thesis-pieces/pkg/store"
	"github.com/MatthewFrisby/thesis-pieces/pkg/tasks"
)

type Store interface {
	GetDocument(ctx context.Context, id int32) (*store.Document, error)
	SetDocumentContent(ctx context.Context, id int32, content string) (*store.Document, error)
}

type S3 interface {
	GetObject(ctx context.Context, bucketName, filename string) ([]byte, error)
}

type DocumentProcessor struct {
	store Store
	s3    S3
}

func NewDocumentProcessor(store Store, s3 S3) *DocumentProcessor {
	return &DocumentProcessor{
		store: store,
		s3:    s3,
	}
}

func (d *DocumentProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p tasks.ProcessDocumentPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	document, err := d.store.GetDocument(ctx, p.DocumentId)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error getting document with id: [%v]", p.DocumentId))
	}

	object, err := d.s3.GetObject(ctx, p.Bucket, document.Filename)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error getting document from bucket with id: [%v]", p.DocumentId))
	}

	var content string
	if document.Filetype == "TXT" {
		content = strings.ToValidUTF8(string(object), "")
	}

	_, err = d.store.SetDocumentContent(ctx, p.DocumentId, content)
	if err != nil {
		return err
	}

	return nil
}
