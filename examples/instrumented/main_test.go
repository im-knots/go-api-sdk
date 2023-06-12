package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/im-knots/go-api-sdk/server"
	"github.com/stretchr/testify/assert"
)

func TestSomeJob(t *testing.T) {
	s := server.NewServer("8080")

	instrumentedService := &InstrumentedService{
		Name: "Instrumented Service",
	}
	s.Name = "Instrumented Server"
	s.RegisterService(instrumentedService)
	for _, service := range s.Services {
		service.RegisterRoutes(s.Engine)
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/", nil)
	s.Engine.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "Worked")
}
