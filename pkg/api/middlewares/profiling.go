package middlewares

import (
	"strconv"

	"cake/pkg/profiling"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func PrometheusMiddleware(c *gin.Context) {
	path := c.Request.URL.Path
	timer := prometheus.NewTimer(profiling.HttpDuration.WithLabelValues(path, "handler"))
	profiling.TotalRequests.WithLabelValues(path, "RECEIVED").Inc()
	c.Next()
	statusCode := c.Writer.Status()
	profiling.TotalRequests.WithLabelValues(path, "DONE").Inc()
	profiling.ResponseStatus.WithLabelValues(path, strconv.Itoa(statusCode)).Inc()
	timer.ObserveDuration()
}
