package server

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"

	"github.com/im-knots/go-api/config"
	"github.com/im-knots/go-api/handlers"
)

type Server struct {
	Config *config.Config
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		Config: cfg,
	}
}

func (s *Server) Start() {
	router := mux.NewRouter()

	router.HandleFunc("/health", handlers.HealthCheckHandler).Methods("GET")
	router.Handle("/metrics", handlers.PrometheusHandler())

	log.Printf("Starting server on port %s", s.Config.Server.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", s.Config.Server.Port), router))
}
