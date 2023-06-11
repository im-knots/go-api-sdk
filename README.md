# Go API SDK

This is a Software Development Kit (SDK) for building APIs in Go. This project uses the Gin web framework, and it is designed to follow best practices for structuring API applications, making it particularly suitable for building APIs intended to work with Kubernetes.

## Getting Started

These instructions will help you understand how to use this SDK for creating your own API applications.

### Prerequisites

- Go 1.20 or later
- Git

### Installation

1. This is currently a private repo. You will need to tell go to not check checksums with GOPRIVATE
    ```
    go env -w GOPRIVATE=github.com/im-knots/*
    ```

2. To use this SDK in your application, import it with:
    ```
    import "github.com/im-knots/go-api-sdk
    ```

3. Use go mod to download the Go dependencies
    ```
    go mod tidy
    ```

## Usage

The main.go file included in the repository is an example of how you can use this SDK to build your own API.

To create your own application, create a new Service and register it with the server. This allows you to define your own routes and handlers.

The server started with this SDK will have two endpoints available by default:

    /health: Returns OK if the server is running.
    /metrics: Returns metrics for Prometheus.

