package internal

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	totalRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests handled by the service.",
	})
)

func init() {
	prometheus.MustRegister(totalRequests)
}

func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
