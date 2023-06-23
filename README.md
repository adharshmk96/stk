# STK

Server toolkit - minimal and simple framework for developing server in golang

## what is included

- [x] go's native http server wrapper
  - uses `net/http` package
  - Get, Post, Put, Delete, Patch methods
- [x] httprouter for routing
- [x] middleware support for all routes
- [x] logger by uber go logrus package
- [x] utils
  - db connection helpers
  - password hashing using argon2
  - loading env variables
    
## get started

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
	server.Get("/", func(c stk.Context) {
		c.Status(http.StatusOK).JSONResponse(stk.Map{"message": "Hello World"})
	})

	// start server
	server.Start()
}
```

## middleware

you can add any middleware by simply creating a function like this and adding it to server.Use()

```go
middleware := func(next stk.HandlerFunc) stk.HandlerFunc {
	return func(c stk.Context) {
		if ctx.Request.URL.Path == "/blocked" {
  			ctx.Status(http.StatusForbidden).JSONResponse("blocked")
			return
  		}
		next(c)
	}
}

server.Use(middleware)
```
