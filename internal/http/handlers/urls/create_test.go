package urls_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/henrywhitaker3/boiler"
	"github.com/henrywhitaker3/go-template/internal/http/handlers/urls"
	"github.com/henrywhitaker3/go-template/internal/test"
	iurls "github.com/henrywhitaker3/go-template/internal/urls"
	"github.com/stretchr/testify/require"
)

func TestItCreatesAUrl(t *testing.T) {
	b := test.Boiler(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	rec := test.Post(
		t,
		b,
		"/urls",
		urls.CreateRequest{},
		"",
	)

	require.Equal(t, http.StatusOK, rec.Code)

	resp := urls.CreateResponse{}
	require.Nil(t, json.Unmarshal(rec.Body.Bytes(), &resp))

	svc, err := boiler.Resolve[iurls.Urls](b)
	require.Nil(t, err)

	_, err = svc.Get(ctx, resp.ID)
	require.ErrorIs(t, err, sql.ErrNoRows)
}
