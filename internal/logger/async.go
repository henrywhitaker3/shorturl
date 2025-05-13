package logger

import (
	"context"
	"log/slog"
	"sync/atomic"
)

type entry struct {
	ctx    context.Context
	record slog.Record
}

type AsyncHandler struct {
	handler slog.Handler
	count   *atomic.Int64
	records chan entry
}

func NewAsyncHandler(ctx context.Context, h slog.Handler) (*AsyncHandler, context.CancelFunc) {
	a := &AsyncHandler{
		handler: h,
		records: make(chan entry, 10000),
		count:   &atomic.Int64{},
	}
	go a.run(ctx)
	return a, a.flush
}

func (a *AsyncHandler) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case entry := <-a.records:
			a.handle(entry.ctx, entry.record)
		}
	}
}

func (a *AsyncHandler) flush() {
	for range a.count.Load() {
		entry := <-a.records
		a.handle(entry.ctx, entry.record)
	}
}

func (a *AsyncHandler) handle(ctx context.Context, record slog.Record) {
	a.count.Add(-1)
	if err := a.handler.Handle(ctx, record); err != nil {
		// do nothing
	}
}

func (a *AsyncHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return a.handler.Enabled(ctx, level)
}

func (a *AsyncHandler) Handle(ctx context.Context, record slog.Record) error {
	a.records <- entry{
		ctx:    ctx,
		record: record,
	}
	a.count.Add(1)
	return nil
}

func (a *AsyncHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return a.handler.WithAttrs(attrs)
}

func (a *AsyncHandler) WithGroup(name string) slog.Handler {
	return a.handler.WithGroup(name)
}

var _ slog.Handler = &AsyncHandler{}
