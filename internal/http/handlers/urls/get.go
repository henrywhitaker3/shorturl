package urls

import (
	"net/http"

	"github.com/henrywhitaker3/boiler"
	"github.com/henrywhitaker3/shorturl/internal/http/common"
	"github.com/henrywhitaker3/shorturl/internal/http/middleware"
	"github.com/henrywhitaker3/shorturl/internal/tracing"
	"github.com/henrywhitaker3/shorturl/internal/urls"
	"github.com/henrywhitaker3/shorturl/internal/uuid"
	"github.com/labstack/echo/v4"
)

type GetHandler struct {
	urls   urls.Urls
	clicks *urls.Clicks
}

func NewGetHandler(b *boiler.Boiler) *GetHandler {
	return &GetHandler{
		urls:   boiler.MustResolve[urls.Urls](b),
		clicks: boiler.MustResolve[*urls.Clicks](b),
	}
}

type GetRequest struct {
	ID uuid.UUID `param:"id"`
}

func (g GetRequest) Validate() error {
	return nil
}

type GetResponse struct {
	Url   *urls.Url   `json:"url"`
	Stats *urls.Stats `json:"stats"`
}

func (g *GetHandler) Handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracing.NewSpan(c.Request().Context(), "GetUrl")
		defer span.End()

		req, ok := common.GetRequest[GetRequest](ctx)
		if !ok {
			return common.ErrBadRequest
		}

		url, err := g.urls.Get(ctx, req.ID)
		if err != nil {
			return common.Stack(err)
		}

		stats, err := g.clicks.Stats(ctx, url.ID)
		if err != nil {
			return common.Stack(err)
		}

		return c.JSON(http.StatusOK, GetResponse{
			Url:   url,
			Stats: stats,
		})
	}
}

func (g *GetHandler) Method() string {
	return http.MethodGet
}

func (g *GetHandler) Path() string {
	return "/urls/:id"
}

func (g *GetHandler) Middleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		middleware.Bind[GetRequest](),
	}
}
