package server

import (
	"fmt"
	"log"
	
	"github.com/gin-gonic/gin"

	"github.com/im-knots/go-api-sdk/config"
	"github.com/im-knots/go-api-sdk/handlers"
)

// Define the Service interface globally.
type Service interface {
	RegisterRoutes(*gin.Engine)
}

type Server struct {
	Config  *config.Config
	Services []Service
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		Config:  cfg,
		Services: make([]Service, 0),
	}
}

func (s *Server) RegisterService(service Service) {
	s.Services = append(s.Services, service)
}

func (s *Server) Start() {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(handlers.PrometheusMiddleware()) // Apply the middleware

	r.GET("/health", handlers.HealthCheckHandler)
	r.GET("/metrics", gin.WrapH(handlers.PrometheusHandler()))

	for _, service := range s.Services {
		service.RegisterRoutes(r)
	}

	log.Printf("Starting server on port %s", s.Config.Server.Port)
	log.Fatal(r.Run(fmt.Sprintf(":%s", s.Config.Server.Port)))
}
