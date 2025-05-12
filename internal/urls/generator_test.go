package urls_test

import (
	"context"
	"testing"
	"time"

	"github.com/henrywhitaker3/boiler"
	"github.com/henrywhitaker3/shorturl/database/queries"
	"github.com/henrywhitaker3/shorturl/internal/test"
	"github.com/henrywhitaker3/shorturl/internal/urls"
)

func TestItGeneratesAliases(t *testing.T) {
	b := test.Boiler(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	gen := urls.NewAliasGenerator(urls.AliasGeneratorOpts{
		DB:         boiler.MustResolve[*queries.Queries](b),
		BufferSize: 100000,
		Interval:   time.Millisecond * 250,
	})
}
