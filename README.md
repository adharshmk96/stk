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

If go isn't configured properly run this
```bash
echo 'export PATH="$PATH:/snap/bin"' >> ~/.profile
echo 'export PATH="$PATH:~/go/bin"' >> ~/.profile
source ~/.profile
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
|
├───.github
│   └───workflows
├───.vscode
├───cmd
├───internals
│   ├───core
│   │   ├───entity
│   │   └───serr
│   ├───http
│   │   ├───handler
│   │   ├───helpers
│   │   └───transport
│   ├───service
│   └───storage
│       └───pingStorage
├───mocks
└───server
    ├───infra
    ├───middleware
    └───routing

```

find more about the project structure [here](docs/project.md)

### Add Modules to project

To add a module to the project run the following command


```bash
stk project module <module-name>
```

It will generate the module in the project structure

```
├───internals
│   |
│   ├───core
│   │   ├───entity
│   │   │       <module-name>.go
│   │   │
│   │   └───serr
│   │           <module-name>.go
│   │
│   ├───http
│   │   ├───handler
│   │   │       <module-name>.go
│   │   │       <module-name>_test.go
│   │   │
│   │   ├───helpers
│   │   │       <module-name>.go
│   │   │
│   │   └───transport
│   │           <module-name>.go
│   │
│   ├───service
│   │       <module-name>.go
│   │       <module-name>_test.go
│   │
│   └───storage
│        └───<module-name>Storage
│               <module-name>.go
│               <module-name>Connection.go
│               <module-name>Queries.go
└───server
    |
    └───routing
        <module-name>.go
```

you can use it by adding `setup<module-name>Routes` to the `routing/initRoutes.go` file

example:

```go
package routing

import (
	"github.com/adharshmk96/stk/gsk"
)

func SetupRoutes(server *gsk.Server) {
	setupPingRoutes(server)
	setupModuleRoutes(server)
}
```

## Development

[refer development docs](docs/development.md)
