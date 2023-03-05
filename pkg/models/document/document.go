package document

import (
	"time"

	"github.com/pkg/errors"

	"github.com/google/uuid"
)

type Filetype string

const (
	FiletypePDF  Filetype = "PDF"
	FiletypeDOCX Filetype = "DOCX"
	FiletypeTXT  Filetype = "TXT"
)

var (
	InvalidFiletypeError = errors.New("invalid filetype, must be a pdf, txt, or docx")
)

func (e Filetype) Valid() bool {
	switch e {
	case FiletypePDF,
		FiletypeDOCX,
		FiletypeTXT:
		return true
	}
	return false
}

type UploadDocument struct {
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	FileName    string   `json:"filename,omitempty"`
	File        []byte   `json:"file,omitempty"`
	FileType    Filetype `json:"filetype,omitempty"`
	UserID      int32    `json:"userId,omitempty"`
	Content     *string  `json:"content,omitempty"`
}

func (u *UploadDocument) SetFile(file []byte) {
	u.File = file
}

func (u *UploadDocument) SetFileName(filename string) {
	u.FileName = filename
}

func (u *UploadDocument) SetFileType(filetype string) error {
	f := Filetype(filetype)
	if !f.Valid() {
		return InvalidFiletypeError
	}
	u.FileType = f
	return nil
}

type SearchDocuments struct {
	Query  string `schema:"query,required"`
	Limit  int32  `schema:"limit"`
	Offset int32  `schema:"offset"`
}

type SearchDocumentsResult struct {
	Uuid        uuid.UUID `json:"uuid,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"UpdatedAt,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Filename    string    `json:"filename,omitempty"`
	Filetype    Filetype  `json:"filetype,omitempty"`
	Username    string    `json:"username,omitempty"`
	Rank        float32   `json:"rank,omitempty"`
}
