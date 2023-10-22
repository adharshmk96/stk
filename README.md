# STK

Server toolkit - minimal and simple framework for developing server in golang

[![Build and Test](https://github.com/adharshmk96/stk/actions/workflows/go-build-test.yml/badge.svg)](https://github.com/adharshmk96/stk/actions/workflows/go-build-test.yml)
[![codecov](https://codecov.io/gh/adharshmk96/stk/graph/badge.svg?token=HMGG55CCLT)](https://codecov.io/gh/adharshmk96/stk)
[![Go Release Workflow](https://github.com/adharshmk96/stk/actions/workflows/go-release.yml/badge.svg)](https://github.com/adharshmk96/stk/actions/workflows/go-release.yml)

STK provides a suite of tools tailored for building and managing server applications.

## Features:

- [GSK (library)](docs/gsk.md): Ideal for constructing servers.
- [STK CLI](#get-started): 
  - Quickly scaffold your project and add modules with ease. It uses gsk package to run the server.
  - [Migrator](#migrator): Generate migration files, perform migration on your sql database.

## Installation

with go install

```bash
go install github.com/adharshmk96/stk@latest
```

If installation fails, check the GOPATH and GOBIN environment variables. Make sure that GOBIN is added to your PATH.
```bash

echo export PATH=$PATH:$GOBIN >> ~/.bashrc
source ~/.bashrc

```

## Get started ( with CLI )

1. Setup and initialize a project scaffolded using gsk and clean arch format. Read more about the project structure [here](docs/project.md)

```bash
stk init
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

Checkout the library documentation [here](docs/gsk.md)

### Add Modules to project

To add a new module to the project run the following command

```bash
stk add module <module-name>
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
