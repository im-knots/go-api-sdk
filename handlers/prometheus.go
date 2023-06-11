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
	httpRequestSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "Size of HTTP requests",
			Buckets: prometheus.ExponentialBuckets(1024, 2, 10), // 1KB, 2KB, 4KB, ... 512KB
		},
		[]string{"method", "endpoint"},
	)
)

type CustomMetrics interface {
	Describe(ch chan<- *prometheus.Desc)
	Collect(ch chan<- prometheus.Metric)
}

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration, httpRequestSize)
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

		// Record request size
		requestSize := float64(c.Request.ContentLength)
		httpRequestSize.WithLabelValues(c.Request.Method, c.Request.URL.Path).Observe(requestSize)
	}
}

func PrometheusHandler() http.Handler {
	return promhttp.Handler()
}
