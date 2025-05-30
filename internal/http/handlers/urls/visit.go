package urls

import (
	"net/http"
	"time"

	"github.com/henrywhitaker3/boiler"
	"github.com/henrywhitaker3/shorturl/internal/config"
	"github.com/henrywhitaker3/shorturl/internal/http/common"
	"github.com/henrywhitaker3/shorturl/internal/http/middleware"
	"github.com/henrywhitaker3/shorturl/internal/logger"
	"github.com/henrywhitaker3/shorturl/internal/queue"
	"github.com/henrywhitaker3/shorturl/internal/tracing"
	"github.com/henrywhitaker3/shorturl/internal/urls"
	"github.com/labstack/echo/v4"
)

type VisitHandler struct {
	urls  urls.Urls
	queue *queue.Publisher
	track bool
}

func NewVisitHandler(b *boiler.Boiler) *VisitHandler {
	return &VisitHandler{
		urls:  boiler.MustResolve[urls.Urls](b),
		queue: boiler.MustResolve[*queue.Publisher](b),
		track: boiler.MustResolve[*config.Config](b).Tracking.Enabled,
	}
}

type VisitRequest struct {
	Alias string `param:"alias"`
}

func (v VisitRequest) Validate() error {
	return nil
}

func (v *VisitHandler) Handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracing.NewSpan(c.Request().Context(), "VisitUrl")
		defer span.End()

		req, ok := common.GetRequest[VisitRequest](ctx)
		if !ok {
			return common.ErrBadRequest
		}

		url, err := v.urls.GetAlias(ctx, req.Alias)
		if err != nil {
			return common.Stack(err)
		}

		if v.track {
			if err := v.queue.Push(ctx, queue.ClickTask, queue.ClickJob{
				ID:   url.ID,
				IP:   c.RealIP(),
				Time: time.Now(),
			}); err != nil {
				logger.Logger(ctx).Error("failed to queue click", "error", err)
			}
		}

		c.Response().
			Header().
			Set(echo.HeaderCacheControl, "no-cache, no-store, max-age=0, must-revalidate")
		c.Response().Header().Set("Pragma", "no-cache")

		return c.Redirect(http.StatusPermanentRedirect, url.Url)
	}
}

func (v *VisitHandler) Method() string {
	return http.MethodGet
}

func (v *VisitHandler) Path() string {
	return "/:alias"
}

func (v *VisitHandler) Middleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		middleware.Bind[VisitRequest](),
	}
}
