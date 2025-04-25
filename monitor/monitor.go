package monitor

import (
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
)

const (
	promSubsystemName = "vivasoft"
)

// NewEchoPrometheusClient adds metricsPath to echo server as middleware
func NewEchoPrometheusClient(e *echo.Echo, metricsPath *string) {
	prom := prometheus.NewPrometheus(promSubsystemName, nil)

	// Scrape metrics from Main Server
	e.Use(prom.HandlerFunc)

	if metricsPath != nil {
		prom.MetricsPath = *metricsPath
	}
	// Setup metrics endpoint at application server
	prom.SetMetricsPath(e)
}
