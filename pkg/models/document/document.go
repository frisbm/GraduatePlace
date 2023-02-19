package document

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
