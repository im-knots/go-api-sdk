package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/im-knots/go-api-sdk/config"
	"github.com/im-knots/go-api-sdk/server"
	"github.com/im-knots/go-api-sdk/handlers"
)

var myCustomCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "my_custom_counter",
		Help: "This is my custom counter",
	},
	[]string{"label1", "label2"},
)

type ExampleService struct{}

func (m *ExampleService) RegisterRoutes(r *gin.Engine) {
	// Register Prometheus middleware
	r.Use(handlers.PrometheusMiddleware())

	r.GET("/api/v1/", m.RootHandler)
	r.GET("/api/v1/report", m.ReportHandler)
	r.GET("/api/v1/crud", m.CrudHandler)
	r.POST("/api/v1/input", m.InputHandler)
}

func (m *ExampleService) RootHandler(c *gin.Context) {
	log.Println("this is a test")
	c.JSON(http.StatusOK, gin.H{"message": "Hello World!"})
}

func (m *ExampleService) ReportHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (m *ExampleService) CrudHandler(c *gin.Context) {
	// Increment the custom metric
	myCustomCounter.WithLabelValues("label1_value", "label2_value").Inc()

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (m *ExampleService) InputHandler(c *gin.Context) {
	var jsonInput map[string]interface{}

	if err := c.ShouldBindJSON(&jsonInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, jsonInput)
}

type MyConfig struct {
	Port string `mapstructure:"port"`
	// Add additional fields as needed
}

func main() {
	// Initialize your config
	cfg := config.NewConfig("default.yaml")

	// Define your config struct
	var myConfig MyConfig
	err := cfg.Unmarshal(&myConfig)
	if err != nil {
		log.Fatalf("Unable to unmarshal config, %v", err)
	}

	s := server.NewServer(myConfig.Port)

	// Register custom metrics
	handlers.RegisterCustomMetrics(myCustomCounter)

	exampleService := &ExampleService{}
	s.RegisterService(exampleService)

	s.Start()
}
