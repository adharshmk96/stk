# GSK Server Package Documentation

[back to main](../README.md)

The GSK server package is a lightweight and flexible HTTP server for Golang applications. It provides an interface for managing server lifecycle, adding middleware, routing HTTP methods, and handling static files. Moreover, the package includes features for automated testing of server responses.

- A web server framework with go's native http server wrapper and httprouter for routing
- Middleware support
- slog Logger
- DB Connection helper functions
- Utilities


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
	server.Get("/", func(gc *gsk.Context) {
		gc.Status(http.StatusOK).JSONResponse(gsk.Map{"message": "Hello World"})
	})

	// start the server
	server.Start()
}
```

### Initialization:

Initialize a new server instance by invoking the `New` function.

```go
server := gsk.New()
```

Optionally, you can pass a `ServerConfig` object to `New` to set the server's configurations, such as the port, logger, and body size limit.

```go
config := &gsk.ServerConfig{
	Port:          ":8081",
	Logger:        slog.New(),
	BodySizeLimit: 1<<20, // 1 MB
}
server := gsk.New(config)
```

### Starting and Stopping:
 
Use `Start` to start the server and `Shutdown` to stop it.

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
package main

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/adharshmk96/stk/gsk"
	"github.com/adharshmk96/stk/pkg/middleware"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func setupRoutes(server *gsk.Server) {
	server.Get("/", func(c *gsk.Context) {
		c.Status(http.StatusOK).JSONResponse(gsk.Map{"message": "Hello World"})
	})
}

func StartHttpServer(port string) (*gsk.Server, chan bool) {

	serverConfig := &gsk.ServerConfig{
		Port:   port,
		Logger: logger,
	}

	server := gsk.New(serverConfig)

	rateLimiter := middleware.NewRateLimiter()
	server.Use(rateLimiter.Middleware)

	server.Use(middleware.RequestLogger)
	server.Use(middleware.CORS(middleware.CORSConfig{
		AllowAll: true,
	}))

	setupRoutes(server)

	server.Start()

	// graceful shutdown
	done := make(chan bool)

	// A go routine that listens for os signals
	// it will block until it receives a signal
	// once it receives a signal, it will shutdown close the done channel
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		if err := server.Shutdown(); err != nil {
			logger.Error("Server Shutdown Failed", err)
		}

		close(done)
	}()

	return server, done
}

func main() {

	startAddr := "0.0.0.0:"
	startingPort := "8080"

	_, done := StartHttpServer(startAddr + startingPort)
	// blocks the routine until done is closed
	<-done

}
```

### Routing HTTP methods: 
Define routes for each HTTP method (Get, Post, Put, Delete, Patch) by calling the appropriate function.

```go
server.Get("/path", func(c *gsk.Context) {
    // handle the request
})
```

#### Route Parameters:

Use the `Param` function to get the value of a route parameter.

```go
server.Get("/path/:id", func(c *gsk.Context) {
	id := c.Param("id")
	// handle the request
})
```

#### Query Parameters:

Use the `Query` function to get the value of a query parameter.

```go
server.Get("/path", func(c *gsk.Context) {
	// example: /path?id=123, id = 123
	id := c.Query("id")
	// handle the request
})
```

#### Request Body:

Use the `Body` function to get the request body as a string.

```go
server.Post("/path", func(c *gsk.Context) {
	body := c.Body()
	// handle the request
})
```

### Route Groups:

Use the `RouteGroup` function to group routes under a common prefix.


```go
server := gsk.New()

apiRoutes := server.RouteGroup("/api")

// all routes registered under /api will be prefixed with /api
apiRoutes.Get("/path", func(c *gsk.Context) {
	// handle the request for /api/path
})

otherRoutes := server.RouteGroup("/other")

// all routes registered under /other will be prefixed with /other
otherRoutes.Get("/path", func(c *gsk.Context) {
	// handle the request for /other/path
})

```

### Serving Static Files: 

Use the `Static` function to serve static files from a specific directory.

```go
server.Static("/assets/*filepath", "/path/to/your/static/files")
```

## Middlewares

Middleware executes code before the request is handled by the route handler. Middleware functions are defined separately and then added to the server using the `Use` function.

```go
// Define your middleware
var MyMiddleware gsk.Middleware = func(next gsk.HandlerFunc) gsk.HandlerFunc {
    return func(c *gsk.Context) {
        // Execute code Before handler 

        next(c)

		// Execute code after handler
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
server.Get("/path", func(c *gsk.Context) {
    // handle the request
})

server.Use(AnotherMiddleware)
// applies another middleware to all routes registered after this
// but not to the route registered before this
// *my middleware will still be applied to this route*
server.Get("/path2", func(c *gsk.Context) {
    // handle the request
})
```

### With Route Groups

Group routes and apply middleware to the group:

```go
server.Use(MyMiddleware)
// applies MyMiddleware to all routes registered after this
authGroup := server.RouteGroup("/auth")
authGroup.Use(AuthMiddleware)

// applies auth middleware to all routes registered after this
authGroup.Get("/login", func(c *gsk.Context) {
	// handle the /auth/login request
})

publicGroup := server.RouteGroup("/public")
// applies only MyMiddleware to all routes registered after this
publicGroup.Get("/home", func(c *gsk.Context) {
	// handle the /public/home request
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

