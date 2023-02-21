package tasks

import (
	"github.com/hibiken/asynq"
)

type Client interface {
	Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error)
}

type TaskManager struct {
	client Client
}

func NewTaskManager(client Client) *TaskManager {
	return &TaskManager{
		client: client,
	}
}

const (
	TypeEmailDelivery = "email:deliver"
)
