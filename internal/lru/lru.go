package lru

import (
	"fmt"
	"sync"
	"sync/atomic"

	lru "github.com/hashicorp/golang-lru/v2"
)

type LRU[T comparable, U any] struct {
	pool []*lru.Cache[T, U]
	next *atomic.Int64
	mu   *sync.Mutex
}

func New[T comparable, U any](poolSize, cacheSize int) (*LRU[T, U], error) {
	if poolSize < 1 {
		return nil, fmt.Errorf("pool size must be positive integer")
	}
	if cacheSize < 1 {
		return nil, fmt.Errorf("cache size must be positive integer")
	}
	cache := &LRU[T, U]{
		pool: []*lru.Cache[T, U]{},
		next: &atomic.Int64{},
		mu:   &sync.Mutex{},
	}

	for range poolSize {
		l, err := lru.New[T, U](cacheSize)
		if err != nil {
			return nil, fmt.Errorf("create lru: %w", err)
		}
		cache.pool = append(cache.pool, l)
	}

	return cache, nil
}

func (l *LRU[T, U]) Get(key T) (U, bool) {
	return l.read().Get(key)
}

func (l *LRU[T, U]) Add(key T, val U) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	evicted := false
	for _, c := range l.pool {
		if c.Add(key, val) {
			evicted = true
		}
	}
	return evicted
}

func (l *LRU[T, U]) Len() int {
	return l.read().Len()
}

func (l *LRU[T, U]) read() *lru.Cache[T, U] {
	n := l.next.Add(1)
	return l.pool[int(n)%len(l.pool)]
}
