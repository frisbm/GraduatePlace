package tasks

type ProcessDocumentPayload struct {
	DocumentId int32
	Bucket     string
}

func (t *TaskManager) ProcessDocumentTask(documentId int32, bucket string) error {
	payload := ProcessDocumentPayload{
		DocumentId: documentId,
		Bucket:     bucket,
	}
	return t.NewTask(ProcessDocumentTask, &payload, DefaultOptions)
}
