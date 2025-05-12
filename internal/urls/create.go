package urls

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/henrywhitaker3/shorturl/internal/queue"
	"github.com/hibiken/asynq"
)

type CreateJobHandler struct {
	svc Urls
}

func NewCreateJobHandler(svc Urls) *CreateJobHandler {
	return &CreateJobHandler{
		svc: svc,
	}
}

func (c *CreateJobHandler) Handle(ctx context.Context, payload []byte) error {
	job := queue.CreateJob{}
	if err := json.Unmarshal(payload, &job); err != nil {
		return fmt.Errorf("unmarhsal job: %w %w", err, asynq.SkipRetry)
	}

	_, err := c.svc.Create(ctx, CreateParams{
		ID:     job.ID,
		Url:    job.Url,
		Domain: job.Domain,
	})

	return err
}
