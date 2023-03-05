package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"code.sajari.com/docconv"
	"github.com/ledongthuc/pdf"

	"github.com/hibiken/asynq"
	"github.com/pkg/errors"

	"github.com/frisbm/graduateplace/pkg/utils"

	"github.com/frisbm/graduateplace/pkg/store"
	"github.com/frisbm/graduateplace/pkg/tasks"
)

type Store interface {
	GetDocument(ctx context.Context, id int32) (*store.Document, error)
	SetDocumentContent(ctx context.Context, id int32, content *string) (*store.Document, error)
}

type S3 interface {
	GetObject(ctx context.Context, bucketName, filename string) (io.ReadCloser, error)
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

	reader, err := d.s3.GetObject(ctx, p.Bucket, document.Filename)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error getting document from bucket with id: [%v]", p.DocumentId))
	}

	content, err := d.process(document.Filetype, reader)
	if err != nil {
		return err
	}

	_, err = d.store.SetDocumentContent(ctx, p.DocumentId, content)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error setting document content with id: [%v]", p.DocumentId))
	}

	return nil
}

func (d *DocumentProcessor) process(filetype string, file io.ReadCloser) (*string, error) {
	buff := bytes.NewBuffer([]byte{})
	size, err := io.Copy(buff, file)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(buff.Bytes())

	switch filetype {
	case "PDF":
		return d.parsePdf(reader, size)
	case "TXT":
		return d.parseTxt(reader)
	case "DOCX":
		return d.parseDocx(reader)
	default:
		return nil, errors.Errorf("filetype not allowed: [%v]", filetype)
	}
}

func (d *DocumentProcessor) cleanContent(content string) *string {
	noNullBytes := strings.ReplaceAll(content, "\u0000", "")
	validUtf8 := strings.ToValidUTF8(noNullBytes, "")
	return utils.Ptr(validUtf8)
}

func (d *DocumentProcessor) parseDocx(reader io.Reader) (*string, error) {
	docx, _, err := docconv.ConvertDocx(reader)
	if err != nil {
		return nil, err
	}
	return d.cleanContent(docx), nil
}

func (d *DocumentProcessor) parsePdf(reader io.ReaderAt, size int64) (content *string, err error) {
	defer func() {
		if r := recover(); r != nil {
			content = nil
			err = errors.New("panic when trying to extract text")
		}
	}()
	pdfReader, err := pdf.NewReader(reader, size)
	if err != nil {
		return nil, err
	}

	b, err := pdfReader.GetPlainText()
	if err != nil {
		return nil, err
	}

	read, err := io.ReadAll(b)
	if err != nil {
		return nil, err
	}

	return d.cleanContent(string(read)), nil
}

func (d *DocumentProcessor) parseTxt(reader io.Reader) (*string, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return d.cleanContent(string(content)), nil
}
