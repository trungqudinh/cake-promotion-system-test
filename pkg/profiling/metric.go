package profiling

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	TotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of get requests.",
		},
		[]string{"path", "state"},
	)

	ResponseStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_response_status",
			Help: "Status of HTTP response",
		},
		[]string{"path", "status"},
	)

	HttpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_response_time_seconds",
		Help:    "Duration of HTTP requests.",
		Buckets: []float64{.001, .003, .005, .01, .025, .05, .1, .25, .5, .7, .8, .9, 1, 2.5, 5, 10},
	}, []string{"path", "phase"})
)

func init() {
	prometheus.Register(TotalRequests)
	prometheus.Register(ResponseStatus)
	prometheus.Register(HttpDuration)
}
