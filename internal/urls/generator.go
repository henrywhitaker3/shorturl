package urls

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"math/rand"
	"time"

	"github.com/henrywhitaker3/shorturl/database/queries"
	"github.com/henrywhitaker3/shorturl/internal/workers"
	"github.com/prometheus/client_golang/prometheus"
)

type AliasGeneratorOpts struct {
	DB *queries.Queries
	// The number of free aliases to keep in memory (default: 10000)
	BufferSize int
	// The interval the generator fills the buffer
	Interval time.Duration

	Registry prometheus.Registerer
}

type AliasGenerator struct {
	db         *queries.Queries
	size       int
	interval   time.Duration
	logger     *slog.Logger
	bufferSize prometheus.Gauge
	collisions prometheus.Counter
}

func NewAliasGenerator(opts AliasGeneratorOpts) *AliasGenerator {
	if opts.BufferSize == 0 {
		opts.BufferSize = 10000
	}
	gen := &AliasGenerator{
		db:       opts.DB,
		size:     opts.BufferSize,
		interval: opts.Interval,
		logger:   slog.Default().With("subsystem", "generator"),
		bufferSize: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "generator_buffer_size",
		}),
		collisions: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "generator_collisiont_total",
		}),
	}

	if opts.Registry != nil {
		if err := opts.Registry.Register(gen.bufferSize); err != nil {
			gen.logger.Error("failed to register metric", "metric", "bufferSize")
		}
		if err := opts.Registry.Register(gen.collisions); err != nil {
			gen.logger.Error("failed to register metric", "metric", "collisions")
		}
	}

	return gen
}

func (a *AliasGenerator) Name() string {
	return "generator"
}

func (a *AliasGenerator) Interval() workers.Interval {
	return workers.NewInterval(a.interval)
}

func (a *AliasGenerator) Timeout() time.Duration {
	return time.Second * 30
}

var (
	chars = []string{
		"a",
		"b",
		"c",
		"d",
		"e",
		"f",
		"g",
		"h",
		"i",
		"j",
		"k",
		"l",
		"m",
		"n",
		"o",
		"p",
		"q",
		"r",
		"s",
		"t",
		"u",
		"v",
		"w",
		"x",
		"y",
		"z",
		"0",
		"1",
		"2",
		"3",
		"4",
		"5",
		"6",
		"7",
		"8",
		"9",
		"0",
	}
	fourChar = int(math.Pow(float64(len(chars)), 4))
	fiveChar = int(math.Pow(float64(len(chars)), 5))
	sixChar  = int(math.Pow(float64(len(chars)), 6))
)

func (a *AliasGenerator) Run(ctx context.Context) error {
	count, err := a.db.CountUrls(ctx)
	if err != nil {
		return fmt.Errorf("could not count urls: %w", err)
	}

	a.logger.Debug("filling up url buffer")

	inBuffer, err := a.db.CountAliasBuffer(ctx)
	if err != nil {
		return fmt.Errorf("could not count alis buffer: %w", err)
	}

	total := inBuffer + count
	length := 4
	if total >= int64(fourChar) {
		length = 5
	}
	if total >= int64(fiveChar) {
		length = 6
	}
	if total > int64(sixChar) {
		length = 7
	}

	toGenerate := a.size - int(inBuffer)
	if toGenerate == 0 {
		a.logger.Debug("buffer already full", "size", inBuffer)
		return nil
	}

	a.logger.Debug("generating aliases", "length", length, "count", toGenerate)
	generated := 0
	for range toGenerate {
		alias := generateAlias(length)
		inserted, err := a.db.InsertAliasBuffer(ctx, alias)
		if err != nil {
			return fmt.Errorf("insert alias into buffer: %w", err)
		}
		if inserted != 1 {
			a.logger.Debug("got a conflict inserting into buffer", "alias", alias)
			continue
		}
		generated++
	}

	a.logger.Info("filled up buffer", "count", generated)
	a.bufferSize.Set(float64(inBuffer + int64(generated)))

	return nil
}

func generateAlias(length int) string {
	out := ""
	for range length {
		out += chars[rand.Intn(len(chars))]
	}
	return out
}

var _ workers.Worker = &AliasGenerator{}
