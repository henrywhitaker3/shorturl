package urls

import (
	"context"
	"fmt"
	"log/slog"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/henrywhitaker3/shorturl/internal/uuid"
	"github.com/prometheus/client_golang/prometheus"
)

type Cache struct {
	svc   *Service
	cache *lru.Cache[string, *Url]

	keys   prometheus.Gauge
	hits   prometheus.Counter
	misses prometheus.Counter
}

type CacheOpts struct {
	Service  *Service
	Size     int
	Registry prometheus.Registerer
}

func NewCache(opts CacheOpts) (*Cache, error) {
	cache, err := lru.New[string, *Url](opts.Size)
	if err != nil {
		return nil, fmt.Errorf("create lru: %w", err)
	}
	c := &Cache{
		svc:   opts.Service,
		cache: cache,
		keys: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "url_cache_keys",
		}),
		hits: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "url_cache_hits_total",
		}),
		misses: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "url_cache_misses_total",
		}),
	}

	if opts.Registry != nil {
		if err := opts.Registry.Register(c.hits); err != nil {
			slog.Error("failed to register cache metric", "metric", "hits")
		}
		if err := opts.Registry.Register(c.keys); err != nil {
			slog.Error("failed to register cache metric", "metric", "keys")
		}
		if err := opts.Registry.Register(c.misses); err != nil {
			slog.Error("failed to register cache metric", "metric", "misses")
		}
	}

	return c, nil
}

func (c *Cache) Get(ctx context.Context, id uuid.UUID) (*Url, error) {
	url, ok := c.cache.Get(id.String())
	if !ok {
		c.misses.Inc()
		var err error
		url, err = c.svc.Get(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("hydrate cache url by id: %w", err)
		}
	} else {
		c.hits.Inc()
	}
	c.cache.Add(id.String(), url)
	c.keys.Set(float64(c.cache.Len()))
	return url, nil
}

func (c *Cache) GetAlias(ctx context.Context, alias string) (*Url, error) {
	url, ok := c.cache.Get(alias)
	if !ok {
		c.misses.Inc()
		var err error
		url, err = c.svc.GetAlias(ctx, alias)
		if err != nil {
			return nil, fmt.Errorf("hydrate cache url by alias: %w", err)
		}
	} else {
		c.hits.Inc()
	}
	c.cache.Add(alias, url)
	c.keys.Set(float64(c.cache.Len()))
	return url, nil
}

func (c *Cache) Count(ctx context.Context) (int, error) {
	return c.svc.Count(ctx)
}

func (c *Cache) Create(ctx context.Context, params CreateParams) (*Url, error) {
	return c.svc.Create(ctx, params)
}

var _ Urls = &Cache{}
