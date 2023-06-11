# Go API SDK

A Software Development Kit (SDK) that assists in the creation of REST APIs using the Go language and the Gin framework. This SDK is designed with Kubernetes and Prometheus in mind, providing out-of-the-box support for a /health and /metrics endpoint, which are useful in a Kubernetes environment for readiness and liveness probes, and for observability with Prometheus, respectively.

Note the main.go included in this repository is an example implementation of the SDK

## Features

* Basic setup for a web server using the Gin framework
* Flexible configuration loading from a YAML file
* Health check endpoint at /health
* Prometheus metrics endpoint at /metrics
* Prometheus middleware for tracking request count and latency with success and error response differentiation
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


### Configuration
The SDK loads the service configuration from a YAML file, the name and location of which are supplied when creating the config object in your code. The configuration loading mechanism is highly flexible, and can accommodate arbitrary fields as per the requirements of your service.

The loaded configuration is unmarshalled into a structure that you define in your service. This structure should include fields that correspond to the keys in your configuration file. For each field, you can add a struct tag `mapstructure:"key"` where `"key"` is the corresponding key in your configuration file. This tag is used to map the configuration data to the correct fields in your struct.

Here is an example of how to define a configuration struct:

```go
type MyConfig struct {
	Port string `mapstructure:"port"`
	// Add additional fields as needed
}
```

And here is how to load the configuration file and unmarshal it into your struct:

```
// initialize your config
cfg := config.NewConfig("default.yaml")

// define your config struct
var myConfig MyConfig
err := cfg.Unmarshal(&myConfig)
if err != nil {
	log.Fatalf("Unable to unmarshal config, %v", err)
}
```

In this example, the port field in `MyConfig` corresponds to the port key in the default.yaml file, which might look like this:

```yaml
port: "8080"
```

Once the configuration is loaded and unmarshalled, you can access the configuration data via your struct. For example, you can access the port number as `myConfig.Port`.

Remember to customize the default.yaml file as per your service's requirements and locate it in the same directory where the application is started.




# Contributors
im-knots