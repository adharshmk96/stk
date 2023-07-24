# STK

Server toolkit - minimal and simple framework for developing server in golang

## Library

GSK - Web server framework [here](#gsk---web-server-framework--library-)

## CLI Tools

There are few cli tools that comes with stk
- Migrator - Database migration tool [here](#migrator)
- Project generator - Generates a new project with gsk following clean architecture (WIP)
- Verify - Verify the project structure for arch rules (WIP)

### Install
```bash
go install github.com/adharshmk96/stk
```

## GSK - Web server framework ( library )

[docs](docs/gsk.md)

- A web server framework with go's native http server wrapper and httprouter for routing
- Middleware support
- Logrus Logger
- DB Connection helper functions
- Utilities

### Get started

```go
package main

import (
	"net/http"

	"github.com/adharshmk96/stk/gsk"
)

func main() {
	// create new server
	server := gsk.New(&config)

	// add routes
	server.Get("/", func(gc *gsk.Context) {
		gc.Status(http.StatusOK).JSONResponse(gsk.Map{"message": "Hello World"})
	})

	// start server
	server.Start()
}
```

### Middleware

you can add any middleware by simply creating a function like this and adding it to server.Use()

```go
middleware := func(next stk.HandlerFunc) stk.HandlerFunc {
	return func(gc stk.Context) {
		if gc.Request.URL.Path == "/blocked" {
  			gc.Status(http.StatusForbidden).JSONResponse("blocked")
			return
  		}
		next(c)
	}
}

server.Use(middleware)
```

# CLI Tools

## Migrator
- CLI tool for generating migration files and running migrations
- Supports sqlite3 (default)

### Get started

Generate migration files ( optinally name it and fill )

```bash
stk migrator generate -n "initial migration" --fill
```

migrate up ( applies all migrations, or specified number of steps )

```bash
stk migrator up
```

migrate down ( applies all down migrations, or specified number of steps )

```bash
stk migrator down
```

History - Shows history of applied migrations

```bash
stk migrator history
```

```
Number  Name               Type  Created     
000001  initial_migration  up    2023-07-01  
000002  initial_migration  up    2023-07-01  
000003  initial_migration  up    2023-07-01  
000004  initial_migration  up    2023-07-01  
000005  initial_migration  up    2023-07-01  
000005  initial_migration  down  2023-07-01  
000004  initial_migration  down  2023-07-01  
000003  initial_migration  down  2023-07-01  
000002  initial_migration  down  2023-07-01  
000001  initial_migration  down  2023-07-01
```


## Development

[refer development docs](docs/development.md)