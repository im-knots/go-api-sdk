package handlers

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func PrometheusHandler() http.Handler {
	return promhttp.Handler()
}
