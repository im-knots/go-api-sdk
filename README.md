# Go API SDK  
  
[![Tests](https://github.com/im-knots/go-api-sdk/actions/workflows/tests.yml/badge.svg)](https://github.com/im-knots/go-api-sdk/actions/)

A Software Development Kit (SDK) that assists in the creation of REST APIs using the Go language and the Gin framework. This SDK is designed with Kubernetes and Prometheus in mind, providing out-of-the-box support for a /health and /metrics endpoint, which are useful in a Kubernetes environment for readiness and liveness probes, and for observability with Prometheus, respectively.  

## Features

* Basic setup for a web server using the Gin framework
* Flexible configuration loading from a YAML file, environment variables, or both.
* Health check endpoint at /health
* Prometheus metrics endpoint at /metrics
* Prometheus middleware for tracking request count, latency, and request/response sizes differentiated by HTTP response code
* Ability to add custom metrics to the Prometheus middleware for more granular monitoring and observability
* An example service with its route registered to demonstrate the usage


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

This SDK is not a standalone application, but a toolkit to facilitate the development of REST APIs. The provided main.go is an example of how to use the SDK and should be modified as per the needs of your project.

Here are the basic steps to use this SDK:

1. Create a new Go project or module and add this SDK as a dependency.
2. Implement the `server.Service` interface for each of your services. This involves creating a struct for your service and a method `RegisterRoutes(*gin.Engine)` for it.
3. Create your configuration struct with fields corresponding to the configuration keys in your YAML file.
4. Load your configuration file and unmarshal it into your configuration struct.
5. Create a new server by calling `server.NewServer(cfg)`, where cfg is your unmarshalled configuration struct.
6. Register each of your services to the server by calling server.`RegisterService(yourService)`.
7. Start the server by calling `server.Start()`.

### Implementing a Custom Service
To create a new service, first define a new struct type that will represent your service. The methods specific to this service, such as its API handlers, will be attached to this struct.

The service must satisfy the `server.Service` interface, which means it must implement a method `RegisterRoutes(*gin.Engine)`. Inside this method, you should register all the routes and handlers for your service.

Here's an example of how to create a new service:
```go
type MyService struct {
	// Fields for your service, if needed
}

func (s *MyService) RegisterRoutes(r *gin.Engine) {
	// Register routes for your service
	r.GET("/myroute", s.myHandler)
}
```

For each route, you should create a corresponding handler method on your service struct. The handler should have the signature `func(*gin.Context)`, and inside it, you can write the logic to process the request and send a response.

Here's an example of how to create a handler for a route:

```go
func (s *MyService) myHandler(c *gin.Context) {
	// Process the request and send a response
	c.JSON(http.StatusOK, gin.H{"message": "Hello from my service!"})
}
```

Then, to use your service with the server, you need to create an instance of your service and register it to the server:

```go
myService := &MyService{}
server.RegisterService(myService)
```

### Adding Custom Metrics

The SDK provides a Prometheus middleware that tracks request count and latency with success and error response differentiation. If you want to add your own custom metrics to this middleware, you can do so by following these steps:

1. Import the necessary Prometheus packages in your code:

    ```go
    import (
        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promhttp"
    )
    ```

2. Create your custom Prometheus metrics using the prometheus.NewXXX functions provided by the Prometheus client library. For example, to create a counter metric:

    ```go
    var myCustomCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "my_custom_counter",
            Help: "This is my custom counter",
        },
        []string{"label1", "label2"},
    )
    ```

3. Register your custom metrics using the prometheus.MustRegister function. You can register multiple metrics at once:

    ```go
    prometheus.MustRegister(myCustomCounter, myOtherMetric, ...)
    ```

4. In your custom service's handler method, increment the custom metric as needed:

    ```go
    myCustomCounter.WithLabelValues("label1_value", "label2_value").Inc()
    ```

    You can place this code inside the handler method where it makes sense in relation to your application logic.

That's it! Your custom metrics will now be tracked by the Prometheus middleware in addition to the default metrics provided by the SDK.


### Configuration
The SDK supports flexible configuration loading from a YAML file or environment variables, or a combination of both. 

#### YAML Configuration File

You can define the configuration in a YAML file. For instance, a default.yaml file may look like this:

```yaml
port: "8080"
```

You load this configuration file and unmarshal it into your configuration struct. Here is an example of how to define a configuration struct:

```go
type MyConfig struct {
	Port string `mapstructure:"port"`
	// Add additional fields as needed
}
```

And here is how to load the configuration file and unmarshal it into your struct:

```go
// initialize your config
cfg := config.NewConfig("default.yaml")

// define your config struct
var myConfig MyConfig
err := cfg.Unmarshal(&myConfig)
if err != nil {
	log.Fatalf("Unable to unmarshal config, %v", err)
}
```

Once the configuration is loaded and unmarshalled, you can access the configuration data via your struct. For example, you can access the port number as `myConfig.Port`.

#### Environment Variables
You can also use environment variables for configuration. If you have a nested structure in your YAML file like:

```yaml
db:
  host: localhost
  port: 5432
```
You can specify it in an environment variable by separating the levels with underscores:

```shell
export DB_HOST=localhost
export DB_PORT=5432
```

Note: The SDK automatically uses environment variables that match the keys in your configuration struct, overriding the values from the configuration file if both are present.

#### Combination of Both
You can even use a combination of both configuration file and environment variables. If a setting is specified in both places, the environment variable will take precedence. This is useful for handling sensitive data like passwords, which can be kept out of the configuration file and instead supplied as environment variables.

Remember to customize the configuration as per your service's requirements. Also note that the YAML configuration file is optional. If it doesn't exist, the application can rely entirely on environment variables.


# Contributors
im-knots

chazapp 
