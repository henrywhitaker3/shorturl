package middleware

import (
	"github.com/henrywhitaker3/shorturl/internal/http/common"
	"github.com/henrywhitaker3/shorturl/internal/tracing"
	"github.com/labstack/echo/v4"
)

type request interface {
	Validate() error
}

func Bind[T request]() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			_, span := tracing.NewSpan(c.Request().Context(), "BindRequest")
			defer span.End()

			var req T
			if err := c.Bind(&req); err != nil {
				return common.ErrBadRequest
			}
			if err := req.Validate(); err != nil {
				return err
			}
			ctx := common.SetRequest(c.Request().Context(), req)
			c.SetRequest(c.Request().WithContext(ctx))
			span.End()
			return next(c)
		}
	}
}
