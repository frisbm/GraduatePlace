package document

import (
	"time"

	"github.com/google/uuid"
)

type UploadDocument struct {
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	FileName    string  `json:"filename,omitempty"`
	File        []byte  `json:"file,omitempty"`
	FileType    string  `json:"filetype,omitempty"`
	UserID      int32   `json:"userId,omitempty"`
	Content     *string `json:"content,omitempty"`
}

func (u *UploadDocument) SetFile(file []byte) {
	u.File = file
}

func (u *UploadDocument) SetFileName(filename string) {
	u.FileName = filename
}

func (u *UploadDocument) SetFileType(filetype string) {
	u.FileType = filetype
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
	Filetype    string    `json:"filetype,omitempty"`
	Username    string    `json:"username,omitempty"`
	Rank        float32   `json:"rank,omitempty"`
}
