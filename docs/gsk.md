# GSK Server Package Documentation

[back to main](../README.md)

The GSK server package is a lightweight and flexible HTTP server for Golang applications. It provides an interface for managing server lifecycle, adding middleware, routing HTTP methods, and handling static files. Moreover, the package includes features for automated testing of server responses.

## Basic Usage

Here's a basic "hello world" server example:

```go
package main

import (
	"net/http"

	"github.com/adharshmk96/stk/gsk"
)

func main() {
	// create a new server
	server := gsk.New()

	// add routes
	server.Get("/", func(gc gsk.Context) {
		gc.Status(http.StatusOK).JSONResponse(gsk.Map{"message": "Hello World"})
	})

	// start the server
	server.Start()
}
```

1. **Initialization:** Initialize a new server instance by invoking the `New` function.

```go
server := gsk.New()
```

Optionally, you can pass a `ServerConfig` object to `New` to set the server's configurations, such as the port, logger, and body size limit.

```go
config := &gsk.ServerConfig{
	Port:          ":8081",
	Logger:        logrus.New(),
	BodySizeLimit: 1<<20, // 1 MB
}
server := gsk.New(config)
```

2. **Starting and Stopping:** Use `Start` to start the server and `Shutdown` to stop it.

```go
server.Start()

// and later...
err := server.Shutdown()
if err != nil {
    // handle error
}
```

Here is an example function to start and gracefully shutdown the server:

```go
func setupRoutes(server gsk.Server) {
    server.Get("/", func(c gsk.Context) {
        c.Status(http.StatusOK).JSONResponse(gsk.Map{"message": "Hello World"})
    })
}

func StartHttpServer(port string) (gsk.Server, chan bool) {
    serverConfig := &gsk.ServerConfig{
		Port:   port,
		Logger: logger,
	}

	server := gsk.New(serverConfig)

    // add middlewares
	rateLimiter := rateLimiter()
	server.Use(rateLimiter)
	server.Use(middleware.RequestLogger)

    // setup routes after adding middleware
	setupRoutes(server)

	server.Start()

	// prepare for graceful shutdown
	done := make(chan bool)

	// A goroutine that listens for OS signals
	// it will block until it receives a signal
	// once it receives a signal, it will shutdown and close the done channel
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		if err := server.Shutdown(); err != nil {
			logger.Error(err)
		}

		close(done)
	}()

	return server, done
}


func main() {
    server, done := StartHttpServer("8080")
    <-done
    logger.Info("Server Stopped")
}
```

3. **Routing HTTP methods:** Define routes for each HTTP method (Get, Post, Put, Delete, Patch) by calling the appropriate function.

```go
server.Get("/path", func(c gsk.Context) {
    // handle the request
})
```

4. **Serving Static Files:** Use the `Static` function to serve static files from a specific directory.

```go
server.Static("/assets/*filepath", "/path/to/your/static/files")
```

## Usage with Middleware

Middleware executes code before the request is handled by the route handler. Middleware functions are defined separately and then added to the server using the `Use` function.

```go
// Define your middleware
var MyMiddleware gsk.Middleware = func(next gsk.HandlerFunc) gsk.HandlerFunc {
    return func(c gsk.Context) {
        // Middleware code here

        next(c)
    }
}

// Add middleware to the server
server.Use(MyMiddleware)
```

Note that middleware is applied when the route is registered. Therefore, make sure to register routes after adding the middleware.

### Middleware Ordering

Middlewares apply only to the routes registered after the middleware is added. You can add some routes which do not require a specific middleware.

```go
server.Use(MyMiddleware)
// applies my middleware to all routes registered after this
server.Get("/path", func(c gsk.Context) {
    // handle the request
})

server.Use(AnotherMiddleware)
// applies another middleware to all routes registered after this
// but not to the route registered before this
// *my middleware will still be applied to this route*
server.Get("/path2", func(c gsk.Context) {
    // handle the request
})
```

### Route Grouping

Group routes and apply middleware to the group:

```go
server.Use(MyMiddleware)
// applies MyMiddleware to all routes registered after this
authGroup := server.RouteGroup("/auth")
authGroup.Use(AuthMiddleware)

// applies auth middleware to all routes registered after this
authGroup.Get("/login", func(c gsk.Context) {
	// handle the /auth/login request
})

publicGroup := server.RouteGroup("/public")
// applies only MyMiddleware to all routes registered after this
publicGroup.Get("/home", func(c gsk.Context) {
	// handle the /public/home request
})
```

### Middleware Ordering

Middlewares will be applied only to the routes registered after the middleware is added. You can add some routes which doesn't require a specific middleware

```go
server.Use(MyMiddleware)
// applies my middleware to all routes registered after this
server.Get("/path", func(c gsk.Context) {
    // handle the request
})

server.Use(AnotherMiddleware)
// applies another middleware to all routes registered after this
// but not to the route registered before this
// *my middleware will still be applied to this route*
server.Get("/path2", func(c gsk.Context) {
    // handle the request
})

```

## Testing Usage

The server package provides a `Test` function to simulate HTTP requests and test server responses. This function takes the HTTP method, path, body, and optional parameters (cookies and headers), and returns a `httptest.ResponseRecorder` and an error.

```go
w, err := server.Test("GET", "/path", nil)

// w is a *httptest.ResponseRecorder, and can be used to check the response
// err should be checked to ensure the request was processed successfully
```

If you need to send headers or cookies with your test request, use the `TestParams` struct.

```go
params := gsk.TestParams{
	Cookies: []*http.Cookie{{Name: "name", Value: "value"}},
	Headers: map[string]string{"Content-Type": "application/json"},
}

w, err := server.Test("GET", "/path", nil, params)
```

Please note that this documentation is a simplified guide to using the GSK package. To use it fully, ensure to handle error cases properly and use the `Context` object correctly in your route handlers and middleware.

