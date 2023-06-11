package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/im-knots/go-api-sdk/config"
	"github.com/im-knots/go-api-sdk/server"
)

type ExampleService struct{}

func (m *ExampleService) RegisterRoutes(r *gin.Engine) {
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
	// initialize your config
	cfg := config.NewConfig("default.yaml")

	// define your config struct
	var myConfig MyConfig
	err := cfg.Unmarshal(&myConfig)
	if err != nil {
		log.Fatalf("Unable to unmarshal config, %v", err)
	}

	s := server.NewServer(myConfig.Port)

	// registering service
	exampleService := &ExampleService{}
	s.RegisterService(exampleService)

	s.Start()
}
