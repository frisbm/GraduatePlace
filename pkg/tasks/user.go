package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

type EmailDeliveryPayload struct {
	UserID     int32
	TemplateID string
}

func (t *TaskManager) SendUserEmailTask(userID int32, tmplID string) error {
	payload := EmailDeliveryPayload{
		UserID:     userID,
		TemplateID: tmplID,
	}
	return t.NewTask(TypeEmailDelivery, &payload, DefaultOptions)
}

func HandleSendUserEmailTask(ctx context.Context, task *asynq.Task) error {
	var p EmailDeliveryPayload
	err := json.Unmarshal(task.Payload(), &p)
	if err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("Sending Email to User: user_id=%d, template_id=%s", p.UserID, p.TemplateID)
	return nil
}
