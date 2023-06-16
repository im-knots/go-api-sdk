package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/im-knots/go-api-sdk/handlers"
	"github.com/im-knots/go-api-sdk/server"
	"github.com/urfave/cli/v2"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
)

type InstrumentedService struct {
	Name string
}

func (m *InstrumentedService) RegisterRoutes(r *gin.Engine) {
	// Register Prometheus middleware
	r.Use(handlers.PrometheusMiddleware())

	r.GET("/api/v1/", m.RootHandler)
	r.Use(otelgin.Middleware(m.Name))
}

func someJob(ctx context.Context, workTime int) {
	tp := otel.Tracer("some-job")
	_, span := tp.Start(ctx, fmt.Sprintf("some-job-%d", workTime))
	defer span.End()
	time.Sleep(time.Duration(workTime) * time.Millisecond)
}

func (m *InstrumentedService) RootHandler(c *gin.Context) {
	workTime := rand.Intn(100)
	someJob(c.Request.Context(), workTime)
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Worked %d ms", workTime)})
}

func main() {
	app := &cli.App{
		Name:  "Instrumented",
		Usage: "Run an API and instrument its code",
		Commands: []*cli.Command{
			{
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "port",
						Value: "8080",
						Usage: "Port to bind the API server",
					},
					&cli.StringFlag{
						Name:  "exporter",
						Value: "localhost:4317",
						Usage: "Address of OTLP Exporter",
					},
				},
				Name:  "start",
				Usage: "Starts the API server",
				Action: func(c *cli.Context) error {
					s := server.NewServer(c.String("port"))
					instrumentedService := &InstrumentedService{
						Name: "Instrumented Service",
					}
					s.Name = "Instrumented Server"
					s.Exporter = c.String("exporter")
					s.RegisterService(instrumentedService)
					s.Start()
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
