# STK

Server toolkit - minimal and simple framework for developing server in golang

[![Build and Test](https://github.com/adharshmk96/stk/actions/workflows/go-build-test.yml/badge.svg)](https://github.com/adharshmk96/stk/actions/workflows/go-build-test.yml)
[![codecov](https://codecov.io/gh/adharshmk96/stk/graph/badge.svg?token=HMGG55CCLT)](https://codecov.io/gh/adharshmk96/stk)
[![Go Release Workflow](https://github.com/adharshmk96/stk/actions/workflows/go-release.yml/badge.svg)](https://github.com/adharshmk96/stk/actions/workflows/go-release.yml)

STK provides a suite of tools tailored for building and managing server applications.

## Features:

- [gsk (library)](docs/gsk.md): Ideal for constructing REST API servers.
- [STK CLI](#get-started): 
  - Quickly scaffold your project and add modules with ease. It uses gsk package to run the server.
  - [Migrator](#migrator): Generate migration files, perform migration on your sql database.

## Installation

with go install

```bash
go install github.com/adharshmk96/stk@latest
```

If go isn't configured properly run this
```bash
echo 'export PATH="$PATH:/snap/bin"' >> ~/.profile
echo 'export PATH="$PATH:~/go/bin"' >> ~/.profile
source ~/.profile
```

## Get started

1. Setup and initialize a project scaffolded using gsk and clean arch format. Read more about the project structure [here](docs/project.md)

```bash
stk init
```

STK init will generate a project in the current directory (default) or directory specified by `-w` flag. with the following structure.


```
│   .gitignore
│   go.mod
│   go.sum
│   main.go
│   makefile
│   README.md
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



2. Start the server

```bash
make run
```

it will run `go run . serve -p 8080` command

3. Test the server

```bash
curl http://localhost:8080/ping
```

### Add Modules to project

To add a module to the project run the following command


```bash
stk add module <module-name>
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



## Development

[refer development docs](docs/development.md)
