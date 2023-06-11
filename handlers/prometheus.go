package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of HTTP operations",
		},
		[]string{"method", "endpoint", "status"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.LinearBuckets(0.01, 0.01, 10), // 10 buckets, each 10ms wide.
		},
		[]string{"method", "endpoint", "status"},
	)
)


func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration)
}

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
			status := "success"
			if c.Writer.Status() >= 400 {
				status = "error"
			}
			httpRequestDuration.WithLabelValues(c.Request.Method, c.Request.URL.Path, status).Observe(v)
		}))

		// Process request
		c.Next()

		// Stop timer and record duration of the request
		timer.ObserveDuration()

		status := "success"
		if c.Writer.Status() >= 400 {
			status = "error"
		}
		httpRequestsTotal.WithLabelValues(c.Request.Method, c.Request.URL.Path, status).Inc()
	}
}


func PrometheusHandler() http.Handler {
	return promhttp.Handler()
}
