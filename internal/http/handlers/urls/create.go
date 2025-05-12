package urls

import (
	"fmt"
	"net/http"

	"github.com/henrywhitaker3/boiler"
	"github.com/henrywhitaker3/go-template/internal/http/common"
	"github.com/henrywhitaker3/go-template/internal/http/middleware"
	"github.com/henrywhitaker3/go-template/internal/queue"
	"github.com/henrywhitaker3/go-template/internal/tracing"
	"github.com/henrywhitaker3/go-template/internal/uuid"
	"github.com/labstack/echo/v4"
)

type CreateHandler struct {
	queue *queue.Publisher
}

func NewCreateHandler(b *boiler.Boiler) *CreateHandler {
	return &CreateHandler{
		queue: boiler.MustResolve[*queue.Publisher](b),
	}
}

type CreateRequest struct {
	Url string `json:"url"`
}

func (c CreateRequest) Validate() error {
	if c.Url == "" {
		return fmt.Errorf("%w url", common.ErrRequiredField)
	}
	return nil
}

type CreateResponse struct {
	ID uuid.UUID `json:"id"`
}

func (h *CreateHandler) Handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracing.NewSpan(c.Request().Context(), "CreateUrl")
		defer span.End()

		req, ok := common.GetRequest[CreateRequest](ctx)
		if !ok {
			return common.ErrBadRequest
		}

		id, err := uuid.Ordered()
		if err != nil {
			return common.Stack(err)
		}

		if err := h.queue.Push(ctx, queue.CreateTask, queue.CreateJob{
			ID:     id,
			Url:    req.Url,
			Domain: c.Request().Host,
		}); err != nil {
			return common.Stack(err)
		}

		return c.JSON(http.StatusAccepted, CreateResponse{
			ID: id,
		})
	}
}

func (h *CreateHandler) Method() string {
	return http.MethodPost
}

func (h *CreateHandler) Path() string {
	return "/urls"
}

func (h *CreateHandler) Middleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		middleware.Bind[CreateRequest](),
	}
}
