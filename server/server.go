package server

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/im-knots/go-api-sdk/handlers"
	"github.com/im-knots/go-api-sdk/instrumentation"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// Define the Service interface globally.
type Service interface {
	RegisterRoutes(*gin.Engine)
}

type Server struct {
	Port     string
	Engine   *gin.Engine
	Services []Service
	Exporter string
	Name     string
}

func NewServer(port string) *Server {
	s := &Server{
		Port:     port,
		Engine:   gin.New(),
		Services: make([]Service, 0),
		Exporter: "",
		Name:     "",
	}

	s.Engine.Use(gin.Logger())
	s.Engine.Use(otelgin.Middleware(s.Name))
	s.Engine.GET("/health", handlers.HealthCheckHandler)
	s.Engine.GET("/metrics", gin.WrapH(handlers.PrometheusHandler()))

	return s
}

func (s *Server) RegisterService(service Service) {
	s.Services = append(s.Services, service)
}

func (s *Server) Start() {
	cleanup := instrumentation.InitTracer(s.Exporter, s.Name)
	defer cleanup(context.Background())

	for _, service := range s.Services {
		service.RegisterRoutes(s.Engine)
	}

	log.Printf("Starting server on port %s", s.Port)
	log.Fatal(s.Engine.Run(fmt.Sprintf(":%s", s.Port)))
}
