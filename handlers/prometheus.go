package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	httpRequestSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "Size of HTTP requests",
			Buckets: prometheus.ExponentialBuckets(1024, 2, 10), // 1KB, 2KB, 4KB, ... 512KB
		},
		[]string{"method", "endpoint", "status"},
	)
	httpResponseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "Size of HTTP responses",
			Buckets: prometheus.ExponentialBuckets(1024, 2, 10), // 1KB, 2KB, 4KB, ... 512KB
		},
		[]string{"method", "endpoint", "status"},
	)
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration, httpRequestSize, httpResponseSize)
}

func RegisterCustomMetrics(metrics ...prometheus.Collector) {
	for _, metric := range metrics {
		prometheus.MustRegister(metric)
	}
}

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
			statusCode := strconv.Itoa(c.Writer.Status())
			httpRequestDuration.WithLabelValues(c.Request.Method, c.Request.URL.Path, statusCode).Observe(v)
		}))

		// Process request
		c.Next()

		// Stop timer and record duration of the request
		timer.ObserveDuration()

		statusCode := strconv.Itoa(c.Writer.Status())
		httpRequestsTotal.WithLabelValues(c.Request.Method, c.Request.URL.Path, statusCode).Inc()

		// Record request and response size
		requestSize := float64(c.Request.ContentLength)
		httpRequestSize.WithLabelValues(c.Request.Method, c.Request.URL.Path, statusCode).Observe(requestSize)

		responseSize := float64(c.Writer.Size())
		httpResponseSize.WithLabelValues(c.Request.Method, c.Request.URL.Path, statusCode).Observe(responseSize)
	}
}

func PrometheusHandler() http.Handler {
	return promhttp.Handler()
}
