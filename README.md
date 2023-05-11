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
    "github.com/adharshmk96/stk"
)

func main() {
    // create new server
    server := stk.NewServer()
    
    // add routes
    server.Get("/", func(ctx *stk.Context) {
        ctx.JSON(200, stk.H{
            "message": "hello world",
        })
    })
    
    // start server
    server.Start()
}
```
