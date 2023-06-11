package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/im-knots/go-api-sdk/config"
	"github.com/im-knots/go-api-sdk/handlers"
	"github.com/im-knots/go-api-sdk/server"
	"github.com/stretchr/testify/assert"
)

var s *server.Server
var r *gin.Engine

func init() {
	cfg := config.NewConfig("default.yaml")
	var myConfig MyConfig
	err := cfg.Unmarshal(&myConfig)
	if err != nil {
		log.Fatalf("Unable to unmarshal config, %v", err)
	}
	s = server.NewServer(myConfig.Port)

	exampleService := &ExampleService{}
	s.RegisterService(exampleService)
	handlers.RegisterCustomMetrics(myCustomCounter)

	r = gin.New()
	r.Use(gin.Logger())
	r.GET("/health", handlers.HealthCheckHandler)
	r.GET("/metrics", gin.WrapH(handlers.PrometheusHandler()))
	for _, service := range s.Services {
		service.RegisterRoutes(r)
	}
}

func TestRootHandler(t *testing.T) {
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/v1/", nil)
	r.ServeHTTP(w, req)
	// Convert the JSON response to a map
	var response map[string]string
	err := json.Unmarshal([]byte(w.Body.String()), &response) // Grab the value & whether or not it exists
	value, exists := response["message"]                      // Make some assertions on the correctness of the response.
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Hello World!", value)
}

func TestHealthRoute(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, w.Body.String(), "OK")
}

func TestCustomMetrics(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/crud", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	req, _ = http.NewRequest("GET", "/metrics", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "my_custom_counter")
}
