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

func main() {
	cfg, err := config.Load("default.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %s", err)
	}

	srv := server.NewServer(cfg)
	srv.RegisterService(&ExampleService{})
	srv.Start()
}
