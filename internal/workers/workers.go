package workers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/henrywhitaker3/shorturl/internal/logger"
	"github.com/henrywhitaker3/shorturl/internal/metrics"
	"github.com/redis/rueidis"
)

type Kind int

const (
	ScheduleInterval Kind = iota
	ScheduleCron
)

type Interval struct {
	raw any
}

func NewInterval(input any) Interval {
	return Interval{raw: input}
}

func (i Interval) Kind() Kind {
	_, ok := i.raw.(time.Duration)
	if ok {
		return ScheduleInterval
	}
	_, ok = i.raw.(string)
	if ok {
		return ScheduleCron
	}
	panic("interval input must be a string or time.Duration")
}

func (i Interval) Interval() time.Duration {
	inter, ok := i.raw.(time.Duration)
	if !ok {
		panic("could not cast Interval to time.Duration")
	}
	return inter
}

func (i Interval) Cron() string {
	cron, ok := i.raw.(string)
	if !ok {
		panic("interval input was not a string")
	}
	return cron
}

type Worker interface {
	Name() string
	Run(ctx context.Context) error
	Interval() Interval
	Timeout() time.Duration
}

type Runner struct {
	sched  gocron.Scheduler
	locker *Locker
	ctx    context.Context
}

func NewRunner(ctx context.Context, redis rueidis.Client) (*Runner, error) {
	locker, err := NewLocker(LockerOpts{
		Redis: redis,
		Topic: "workers",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialise locker: %w", err)
	}

	sched, err := gocron.NewScheduler(
		gocron.WithDistributedLocker(locker),
	)
	if err != nil {
		return nil, fmt.Errorf("created scheduler: %w", err)
	}

	return &Runner{
		sched:  sched,
		locker: locker,
		ctx:    ctx,
	}, nil
}

func (r *Runner) Register(w Worker) error {
	logger := logger.Logger(r.ctx).With("subsystem", "runner")

	logger.Info("registering worker", "name", w.Name())

	var at gocron.JobDefinition
	switch w.Interval().Kind() {
	case ScheduleCron:
		at = gocron.CronJob(w.Interval().Cron(), false)
	case ScheduleInterval:
		at = gocron.DurationJob(w.Interval().Interval())
	default:
		return errors.New("invalid schedule kind")
	}

	_, err := r.sched.NewJob(
		at,
		gocron.NewTask(
			func() {
				ctx, cancel := context.WithTimeout(context.Background(), w.Timeout())
				defer cancel()
				metrics.WorkerExecutions.WithLabelValues(w.Name()).Inc()
				if err := w.Run(ctx); err != nil {
					logger.Error("worker run failed", "name", w.Name(), "error", err)
					metrics.WorkerExecutionErrors.WithLabelValues(w.Name()).Inc()
				}
			},
		),
		gocron.WithName(w.Name()),
	)
	if err != nil {
		return fmt.Errorf("failed to register worker: %w", err)
	}
	return nil
}

func (r *Runner) Run() {
	go r.locker.Run(r.ctx)
	<-r.locker.Initialised()
	r.sched.Start()
}

func (r *Runner) Stop() error {
	return r.sched.Shutdown()
}
