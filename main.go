package main

import (
	"github.com/im-knots/go-api/config"
	"github.com/im-knots/go-api/server"
)

func main() {
	cfg, err := config.Load("default.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %s", err)
	}

	srv := server.NewServer(cfg)
	srv.Start()
}
