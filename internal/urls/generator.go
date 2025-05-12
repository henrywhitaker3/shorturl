package urls

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"math/rand"
	"time"

	"github.com/henrywhitaker3/shorturl/internal/workers"
	"github.com/prometheus/client_golang/prometheus"
)

type AliasGeneratorOpts struct {
	Alias *Alias

	// The number of free aliases to keep in memory (default: 10000)
	BufferSize int
	// The interval the generator fills the buffer
	Interval time.Duration

	// The length of the alias
	Length int

	Registry prometheus.Registerer
}

type AliasGenerator struct {
	alias      *Alias
	size       int
	interval   time.Duration
	length     int
	logger     *slog.Logger
	generated  prometheus.Counter
	collisions prometheus.Counter
}

func NewAliasGenerator(opts AliasGeneratorOpts) *AliasGenerator {
	if opts.BufferSize == 0 {
		opts.BufferSize = 10000
	}
	gen := &AliasGenerator{
		alias:    opts.Alias,
		size:     opts.BufferSize,
		interval: opts.Interval,
		length:   opts.Length,
		logger:   slog.Default().With("subsystem", "generator"),
		generated: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "generator_generated",
		}),
		collisions: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "generator_collisiont_total",
		}),
	}

	if opts.Registry != nil {
		if err := opts.Registry.Register(gen.generated); err != nil {
			gen.logger.Error("failed to register metric", "metric", "generated")
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

func (a *AliasGenerator) Run(ctx context.Context) error {
	a.logger.Debug("filling up url buffer")

	aliases, err := a.generateAliases(ctx)
	if err != nil {
		return fmt.Errorf("generate aliases: %w", err)
	}

	if len(aliases) == 0 {
		a.logger.Info("no aliases to generate")
		return nil
	}

	generated := 0
	for _, al := range aliases {
		err := a.alias.Create(ctx, al)
		if err != nil {
			a.logger.Error("could not store alias", "alias", al, "error", err)
			continue
		}
		generated++
	}

	a.logger.Info("filled up buffer", "count", generated)

	return nil
}

func (a *AliasGenerator) generateAliases(ctx context.Context) ([]string, error) {
	free, err := a.alias.CountFree(ctx)
	if err != nil {
		return nil, err
	}

	toGenerate := a.size - free
	if toGenerate == 0 {
		return []string{}, nil
	}

	a.logger.Debug("generating aliases", "length", a.length, "count", toGenerate)
	aliases := []string{}
	for range toGenerate {
		aliases = append(aliases, generateAlias(a.length))
	}

	filtered, err := a.alias.FilterOutExisting(ctx, aliases)
	if err != nil {
		return nil, fmt.Errorf("could not filter existing alises out of generated: %w", err)
	}

	return filtered, nil
}

func generateAlias(length int) string {
	out := ""
	for range length {
		out += chars[rand.Intn(len(chars))]
	}
	return out
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

var _ workers.Worker = &AliasGenerator{}
