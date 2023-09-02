# STK

Server toolkit - minimal and simple framework for developing server in golang

## Library

GSK - Web server framework [here](#gsk---web-server-framework--library-)

## CLI Tools

There are few cli tools that comes with stk
- Migrator - Database migration tool [here](#migrator)
- Project generator - Generates a new project with gsk following clean architecture
- Verify - Verify the project structure for arch rules (WIP)

### Install
```bash
go install github.com/adharshmk96/stk
```

## GSK - Web server framework ( library )

[docs](docs/gsk.md)

- A web server framework with go's native http server wrapper and httprouter for routing
- Middleware support
- slog Logger
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
	server := gsk.New()

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

NOTE: Middleware functions only wraps registered routes.

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

## Project Generator

- Generates a new project with gsk following clean architecture

### Get started

1. goto working directory `cd <target directory>`
2. run the following command

```bash
stk project generate
```

The command will generate a project with the following structure

```
│   .gitignore
│   go.mod
│   go.sum
│   main.go
│   makefile
│   README.md
│   request.http
│
├───.github
│   └───workflows
│           go-build-test.yml
│           go-release.yml
│
├───cmd
│       root.go
│       serve.go
│       version.go
│
├───internals
│   ├───core
│   │   │   handler.go
│   │   │   service.go
│   │   │   storage.go
│   │   │
│   │   ├───ds
│   │   │       ping.go
│   │   │
│   │   └───serr
│   │           pingerr.go
│   │
│   ├───http
│   │   └───handler
│   │           handler.go
│   │           ping.go
│   │
│   ├───service
│   │       ping.go
│   │       service.go
│   │
│   └───storage
│       └───sqlite
│               ping.go
│               sqlite.go
│
└───server
    │   setup.go
    │
    ├───infra
    │       config.go
    │       constants.go
    │       logger.go
    │
    ├───middleware
    │       middleware.go
    │
    └───routing
            routing.go

```

find more about the project structure [here](docs/project.md)



## Development

[refer development docs](docs/development.md)
