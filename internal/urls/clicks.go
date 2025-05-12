package urls

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/henrywhitaker3/shorturl/database/queries"
	"github.com/henrywhitaker3/shorturl/internal/queue"
	"github.com/henrywhitaker3/shorturl/internal/uuid"
	"github.com/hibiken/asynq"
)

type Clicks struct {
	db *queries.Queries
}

type ClickOpts struct {
	DB *queries.Queries
}

func NewClicks(opts ClickOpts) *Clicks {
	return &Clicks{
		db: opts.DB,
	}
}

type StoreClick struct {
	ID   uuid.UUID
	IP   string
	Time time.Time
}

func (c *Clicks) Click(ctx context.Context, params StoreClick) error {
	id, err := uuid.Ordered()
	if err != nil {
		return fmt.Errorf("generate click id: %w", err)
	}

	if err := c.db.StoreClick(ctx, queries.StoreClickParams{
		ID:        id.UUID(),
		UrlID:     params.ID.UUID(),
		Ip:        params.IP,
		ClickedAt: params.Time.Unix(),
	}); err != nil {
		return fmt.Errorf("store click: %w", err)
	}

	return nil
}

type Stats struct {
	Clicks int `json:"clicks"`
}

func (c *Clicks) Stats(ctx context.Context, id uuid.UUID) (*Stats, error) {
	clicks, err := c.db.CountClicks(ctx, id.UUID())
	if err != nil {
		return nil, fmt.Errorf("could not count clicks: %w", err)
	}
	return &Stats{
		Clicks: int(clicks),
	}, nil
}

type ClickJobHandler struct {
	svc *Clicks
}

func NewClickJobHandler(svc *Clicks) *ClickJobHandler {
	return &ClickJobHandler{svc: svc}
}

func (c *ClickJobHandler) Handle(ctx context.Context, payload []byte) error {
	job := queue.ClickJob{}
	if err := json.Unmarshal(payload, &job); err != nil {
		return fmt.Errorf("unmarshal click job: %w %w", err, asynq.SkipRetry)
	}

	if err := c.svc.Click(ctx, StoreClick{
		ID:   job.ID,
		IP:   job.IP,
		Time: job.Time,
	}); err != nil {
		return fmt.Errorf("store click: %w", err)
	}
	return nil
}
