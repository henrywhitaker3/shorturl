package logger

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/henrywhitaker3/shorturl/internal/http/common"
)

func Setup(ctx context.Context, level slog.Level, outputs ...io.Writer) context.CancelFunc {
	if len(outputs) == 0 {
		outputs = []io.Writer{os.Stdout}
	}
	handler := slog.NewJSONHandler(outputs[0], &slog.HandlerOptions{
		Level:     level,
		AddSource: false,
	})
	async, cancel := NewAsyncHandler(ctx, handler)
	logger := slog.New(async)
	slog.SetDefault(logger)
	return cancel
}

func Logger(ctx context.Context) *slog.Logger {
	log := slog.Default()
	if trace := common.TraceID(ctx); trace != "" {
		log = log.With("trace_id", trace)
	}
	if req := common.ContextID(ctx); req != "" {
		log = log.With("request_id", req)
	}
	return log
}
