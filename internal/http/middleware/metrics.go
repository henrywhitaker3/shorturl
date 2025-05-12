package middleware

import (
	"net/http"
	"strings"

	"github.com/henrywhitaker3/shorturl/internal/config"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
)

func Metrics(conf config.Telemetry, reg prometheus.Registerer) echo.MiddlewareFunc {
	return echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
		Subsystem: strings.ReplaceAll(conf.Tracing.ServiceName, "-", "_"),
		Skipper: func(c echo.Context) bool {
			return c.Request().Method == http.MethodOptions
		},
		Registerer: reg,
		HistogramOptsFunc: func(opts prometheus.HistogramOpts) prometheus.HistogramOpts {
			opts.Buckets = []float64{
				.0001,
				.0002,
				.0003,
				.0004,
				.0005,
				.0006,
				.0007,
				.0008,
				.0009,
				.001,
				.002,
				.003,
				.004,
				.005,
				.006,
				.007,
				.008,
				.009,
				.01,
				.025,
				.05,
				.1,
				.25,
				.5,
				1,
				2.5,
				5,
				10,
			}
			return opts
		},
	})
}
