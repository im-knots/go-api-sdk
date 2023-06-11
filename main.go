package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/im-knots/go-api-sdk/config"
	"github.com/im-knots/go-api-sdk/server"
)

type ExampleService struct{}

func (m *ExampleService) RegisterRoutes(r *gin.Engine) {
	r.GET("/example", m.ExampleHandler)
}

func (m *ExampleService) ExampleHandler(c *gin.Context) {
	log.Println("this is a test")
	c.JSON(200, gin.H{"message": "Hello World!"})
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
