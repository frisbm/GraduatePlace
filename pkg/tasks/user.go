package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/MatthewFrisby/thesis-pieces/pkg/constants"

	"github.com/hibiken/asynq"
)

type EmailDeliveryPayload struct {
	UserID     int32
	TemplateID string
}

func (t *TaskManager) SendUserEmailTask(userID int32, tmplID string) error {
	payload, err := json.Marshal(EmailDeliveryPayload{UserID: userID, TemplateID: tmplID})
	if err != nil {
		return err
	}
	task := asynq.NewTask(
		TypeEmailDelivery,
		payload,
		asynq.MaxRetry(1),
		asynq.Timeout(time.Minute),
		asynq.Retention(24*time.Hour),
		asynq.Queue(constants.HIGH_PRIORITY_QUEUE),
	)
	info, err := t.client.Enqueue(task)
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
	return nil
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
