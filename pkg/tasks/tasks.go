package tasks

import (
	"encoding/json"
	"log"
	"time"

	"github.com/hibiken/asynq"

	"github.com/MatthewFrisby/thesis-pieces/pkg/constants"
	"github.com/MatthewFrisby/thesis-pieces/pkg/utils"
)

type Options struct {
	MaxRetry  *int
	Timeout   *time.Duration
	Retention *time.Duration
	Queue     *string
	Deadline  *time.Time
	ProcessAt *time.Time
	ProcessIn *time.Duration
}

var DefaultOptions = Options{
	MaxRetry:  utils.Ptr(2),
	Timeout:   utils.Ptr(5 * time.Minute),
	Retention: utils.Ptr(72 * time.Hour),
	Queue:     utils.Ptr(constants.LOW_PRIORITY_QUEUE),
	Deadline:  nil,
	ProcessAt: nil,
	ProcessIn: nil,
}

func (o *Options) toAsync() []asynq.Option {
	var opts []asynq.Option
	if o.MaxRetry != nil {
		opts = append(opts, asynq.MaxRetry(*o.MaxRetry))
	}
	if o.Timeout != nil {
		opts = append(opts, asynq.Timeout(*o.Timeout))
	}
	if o.Retention != nil {
		opts = append(opts, asynq.Retention(*o.Retention))
	}
	if o.Queue != nil {
		opts = append(opts, asynq.Queue(*o.Queue))
	}
	if o.Deadline != nil {
		opts = append(opts, asynq.Deadline(*o.Deadline))
	}
	if o.ProcessAt != nil {
		opts = append(opts, asynq.ProcessAt(*o.ProcessAt))
	}
	if o.ProcessIn != nil {
		opts = append(opts, asynq.ProcessIn(*o.ProcessIn))
	}
	return opts
}

const (
	TypeEmailDelivery = "email:deliver"
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

func (t *TaskManager) NewTask(typename string, payload any, opts Options) error {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(typename, bytes, opts.toAsync()...)
	info, err := t.client.Enqueue(task)
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
		return err
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
	return nil
}
