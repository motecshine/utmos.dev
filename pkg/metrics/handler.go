package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Handler returns a Gin handler for the /metrics endpoint.
func Handler(collector *Collector) gin.HandlerFunc {
	h := promhttp.HandlerFor(collector.Registry(), promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	})

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
