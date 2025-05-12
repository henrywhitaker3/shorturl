package urls_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/henrywhitaker3/boiler"
	"github.com/henrywhitaker3/shorturl/internal/http/handlers/urls"
	"github.com/henrywhitaker3/shorturl/internal/test"
	iurls "github.com/henrywhitaker3/shorturl/internal/urls"
	"github.com/stretchr/testify/require"
)

func TestItCreatesAUrl(t *testing.T) {
	b := test.Boiler(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	rec := test.Post(
		t,
		b,
		"/urls",
		urls.CreateRequest{
			Url: "https://synthetigo.com",
		},
		"",
	)

	require.Equal(t, http.StatusAccepted, rec.Code)

	resp := urls.CreateResponse{}
	require.Nil(t, json.Unmarshal(rec.Body.Bytes(), &resp))

	// Fill up the buffer
	require.Nil(t, boiler.MustResolve[*iurls.AliasGenerator](b).Run(ctx))

	svc, err := boiler.Resolve[iurls.Urls](b)
	require.Nil(t, err)

	_, err = svc.Get(ctx, resp.ID)
	require.ErrorIs(t, err, sql.ErrNoRows)

	t.Log("running queue")
	test.RunQueues(t, b, ctx)

	time.Sleep(time.Second * 2)

	url, err := svc.Get(ctx, resp.ID)
	require.Nil(t, err)
	t.Log(url)
}
