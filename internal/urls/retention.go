package urls

import (
	"context"
	"log/slog"
	"time"

	"github.com/henrywhitaker3/shorturl/internal/config"
	"github.com/henrywhitaker3/shorturl/internal/workers"
)

type Retention struct {
	clicks  *Clicks
	enabled bool
	period  time.Duration
}

type RetentionOpts struct {
	Clicks *Clicks
	Config config.Retention
}

func NewRetention(opts RetentionOpts) *Retention {
	return &Retention{
		clicks:  opts.Clicks,
		enabled: opts.Config.Enabled,
		period:  opts.Config.Period,
	}
}

func (r *Retention) Name() string {
	return "retention"
}

func (r *Retention) Timeout() time.Duration {
	return time.Second * 30
}

func (r *Retention) Interval() workers.Interval {
	return workers.NewInterval(time.Minute)
}

func (r *Retention) Run(ctx context.Context) error {
	if !r.enabled {
		return nil
	}

	deleted, err := r.clicks.Delete(ctx, time.Now().Add(-r.period))
	if err != nil {
		return err
	}

	slog.Info("deleted click data", "count", deleted)

	return nil
}

var _ workers.Worker = &Retention{}
