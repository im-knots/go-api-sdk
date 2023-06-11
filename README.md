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

## Versioning
This go module uses github actions to bump the module version on PR merges into main

### Bumping

Manual Bumping: Any commit message that includes #major, #minor, #patch, or #none will trigger the respective version bump. If two or more are present, the highest-ranking one will take precedence. If #none is contained in the merge commit message, it will skip bumping regardless DEFAULT_BUMP.

Automatic Bumping: If no #major, #minor or #patch tag is contained in the merge commit message, it will bump whichever DEFAULT_BUMP is set to (which is minor by default). Disable this by setting DEFAULT_BUMP to none.

# Contributors
im-knots

this is a test of gh actions