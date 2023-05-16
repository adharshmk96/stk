# stk

Server tool kit - framework for developing server in golang

## what is included

- [x] go's native http server wrapper
  - use `net/http` package
  - Get, Post, Put, Delete, Patch methods
- [x] Middleware support
  - support middleware for all routes
- [x] logger
  - zap logger by uber go zap package
- [x] utils
  - password hashing using argon2
  - loading env variables
    
## usage

```go
package main

import (
	"net/http"

	"github.com/adharshmk96/stk"
)

func main() {
	config := stk.ServerConfig{
		Port:           "0.0.0.0:8080",
		RequestLogging: true,
	}
	// create new server
	server := stk.NewServer(&config)

	// add routes
	server.Get("/", func(c *stk.Context) {
		c.Status(http.StatusOK).JSONResponse(stk.Map{"message": "Hello World"})
	})

	// start server
	server.Start()
}
```
