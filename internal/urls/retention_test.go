package urls_test

import (
	"context"
	"testing"
	"time"

	"github.com/henrywhitaker3/boiler"
	"github.com/henrywhitaker3/shorturl/internal/config"
	"github.com/henrywhitaker3/shorturl/internal/test"
	"github.com/henrywhitaker3/shorturl/internal/urls"
	"github.com/stretchr/testify/require"
)

func TestItDeletesClick(t *testing.T) {
	b := test.Boiler(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	url := test.Url(t, b, test.UrlOpts{})

	clicks := boiler.MustResolve[*urls.Clicks](b)

	stats, err := clicks.Stats(ctx, url.ID)
	require.Nil(t, err)
	require.Zero(t, stats.Clicks)

	require.Nil(t, clicks.Click(ctx, urls.StoreClick{
		ID:   url.ID,
		IP:   "127.0.0.1",
		Time: time.Now().Add(-time.Hour),
	}))
	require.Nil(t, clicks.Click(ctx, urls.StoreClick{
		ID:   url.ID,
		IP:   "127.0.0.1",
		Time: time.Now(),
	}))

	stats, err = clicks.Stats(ctx, url.ID)
	require.Nil(t, err)
	require.Equal(t, 2, stats.Clicks)

	retention := urls.NewRetention(urls.RetentionOpts{
		Clicks: clicks,
		Config: config.Retention{
			Enabled: true,
			Period:  time.Minute,
		},
	})
	require.Nil(t, retention.Run(ctx))

	stats, err = clicks.Stats(ctx, url.ID)
	require.Nil(t, err)
	require.Equal(t, 1, stats.Clicks)
}
