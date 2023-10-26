package tpl

var MAKEFILE_TPL = Template{
	FilePath: "makefile",
	Content: `publish:
	@git push && semver push


##########################
### Build Commands
##########################

BINARY_NAME={{ .AppName }}

build:
	@go build -o ./out/$(BINARY_NAME) -v

run: 
	@go run . serve -p 8080

test:
	@go test ./...

coverage:
	@go test -v ./... -coverprofile=coverage.out 
	@go tool cover -html=coverage.out

testci:
	@go test ./... -coverprofile=coverage.out

clean:
	@go clean
	@rm -f ./out/$(BINARY_NAME)
	@rm -f coverage.out

deps:
	@go mod download

tidy:
	@go mod tidy

lint:
	@golangci-lint run --enable-all

vet:
	@go vet

clean-branch:
	@git branch | egrep -v "(^\*|main|master)" | xargs git branch -D
	@git tag -d $(shell git tag -l)

	
##########################
### Setup Commands
##########################

init: 
	@go mod tidy
# Install tools
	@go install github.com/adharshmk96/semver@latest
	@go install github.com/vektra/mockery/v2@latest
# Setup Git hooks
	@git config core.hooksPath .githooks

# mockgen:
	@rm -rf ./mocks
	@mockery --all	

	@echo "Project initialized."
`,
}

var GITIGNORE_TPL = Template{
	FilePath: ".gitignore",
	Content: `# If you prefer the allow list template instead of the deny list, see community template:
# https://github.com/github/gitignore/blob/main/community/Golang/Go.AllowList.gitignore
#
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out
dist

*.db`,
}

var MAINGO_TPL = Template{
	FilePath: "main.go",
	Content: `package main

import "{{ .PkgName }}/cmd"

func main() {
	cmd.Execute()
}
`,
}

var READMEMD_TPL = Template{
	FilePath: "README.md",
	Content: `# {{ .AppName }}

Run

go run main.go serve -p 8080

- or -

make run


## Project Structure Documentation

This is the project structure generated by "stk init" command.

## Components

### Internals

The internals directory contains the application logic and http , organized into various segments:
- Entity & Errors
- Http Handlers & Helpers
- Service
- Storage

Structure:

- internals/core
  - /entity - Primary data structures and interfaces for a module.
  - /serr  - Server specific errors for each modules.
- internals/http
  - /handler - Http handler functions for each modules which process the request and response. Implements the interfaces defined in internals/core/entity for each modules.
  - /helpers - Helper functions for http handlers.
  - /transport - Data structure and Functions related to http response.
- internals/service - Business logic for each modules. Implements the interfaces defined in internals/core/entity for each modules.
- internals/storage
  - /<module-name>Storage - Storage implementation for each modules, Initialize and implement the storage interface defined in internals/core/entity for each modules.
---

### Server

The server directory contains logic related to platform level concerns, such as:
- Server configurations
- Middleware
- Routing
- Infrastructure

Structure:

- server/infra - Server configurations and constants.
- server/middleware - Middleware functions for the server.
- server/routing - Routing logic for the server, Binds the http handlers to the routes for each Modules.
- server/server.go - Server initialization and start logic.

---

### Cmd

- cmd - Entry point for the application or any related command-line interfaces (CLI). These scripts initialize and run the application, utilizing the Cobra CLI library. Serve and version commands are generated by default.

---

### Platform Configurations

- .github/workflows 
  
  For github workflows, The default workflow for build, test, release is generated by default. Go releaser is used to execute the release job. Release workflow is triggered when a new tag is pushed to the repository.

- .vscode

  For vscode configurations, such as launch.json. The default launch configurations for debugging is generated by default.

---

For managing versions, [semver](https://github.com/adharshmk96/semver) is reccomended.
For testing, [mockery](https://github.com/vektra/mockery) is reccomended.
`,
}

var CMD_ROOTGO_TPL = Template{
	FilePath: "cmd/root.go",
	Content: `package cmd

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var SemVer = "v0.0.0"

func GetSemverInfo() string {
	if SemVer != "v0.0.0" {
		return SemVer
	}
	version, ok := debug.ReadBuildInfo()
	if ok && version.Main.Version != "(devel)" && version.Main.Version != "" {
		return version.Main.Version
	}
	return SemVer
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "{{ .AppName }}",
	Short: "{{ .AppName }} is an stk project.",
	Long:  "{{ .AppName }} is generated using stk cli.",
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("version").Value.String() == "true" {
			fmt.Println(GetSemverInfo())
		} else {
			cmd.Help()
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", ".stk.yaml", "config file.")
	rootCmd.Flags().BoolP("version", "v", false, "display {{ .AppName }} version")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv()

	// Set the key replacer for env variables.
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

}
`,
}

var CMD_SERVEGO_TPL = Template{
	FilePath: "cmd/serve.go",
	Content: `package cmd

import (
	"sync"

	"{{ .PkgName }}/server"
	"github.com/spf13/cobra"
)

var startingPort string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup

		wg.Add(1)

		startAddr := "0.0.0.0:"

		go func() {
			defer wg.Done()
			_, done := server.StartHttpServer(startAddr + startingPort)
			// blocks the routine until done is closed
			<-done
		}()

		wg.Wait()
	},
}

func init() {
	serveCmd.Flags().StringVarP(&startingPort, "port", "p", "8080", "Port to start the server on")

	rootCmd.AddCommand(serveCmd)
}
`,
}

var VSCODE_LAUNCHJSON_TPL = Template{
	FilePath: ".vscode/launch.json",
	Content: `{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "serve in port 8080",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/main.go",
            "args": ["serve", "-p", "8080"]
          }
    ]
}`,
}

var SERVER_SETUPGO_TPL = Template{
	FilePath: "server/setup.go",
	Content: `package server

import (
	"os"
	"os/signal"
	"syscall"

	"{{ .PkgName }}/server/infra"
	svrmw "{{ .PkgName }}/server/middleware"
	"{{ .PkgName }}/server/routing"
	"github.com/adharshmk96/stk/gsk"
	"github.com/adharshmk96/stk/pkg/middleware"
)

func StartHttpServer(port string) (*gsk.Server, chan bool) {

	logger := infra.GetLogger()

	serverConfig := &gsk.ServerConfig{
		Port:   port,
		Logger: logger,
	}

	server := gsk.New(serverConfig)

	rateLimiter := svrmw.RateLimiter()
	server.Use(rateLimiter)
	server.Use(middleware.RequestLogger)
	server.Use(middleware.CORS(middleware.CORSConfig{
		AllowAll: true,
	}))

	infra.LoadDefaultConfig()

	routing.SetupTemplateRoutes(server)
	routing.SetupApiRoutes(server)

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
			logger.Error(err.Error())
		}

		close(done)
	}()

	return server, done
}
`,
}

var SERVER_ROUTING_INITROUTESGO_TPL = Template{
	FilePath: "server/routing/initRoutes.go",
	Content: `package routing

import (
	"github.com/adharshmk96/stk/gsk"
)

func SetupApiRoutes(server *gsk.Server) {
	apiRoutes := server.RouteGroup("/api")

	setup{{ .ExportedName }}Routes(apiRoutes)
}

func SetupTemplateRoutes(server *gsk.Server) {
	templateRoutes := server.RouteGroup("/")

	templateRoutes.Get("/", func(gc *gsk.Context) {
		gc.TemplateResponse(&gsk.Tpl{
			TemplatePath: "public/templates/index.html",
			Variables: gsk.Map{
				"Title":   "STK",
				"Content": "Hello, World!",
			},
		})
	})
}
`,
}

var SERVER_MIDDLEWARE_MIDDLEWAREGO_TPL = Template{
	FilePath: "server/middleware/middleware.go",
	Content: `package middleware

import (
	"time"

	"github.com/adharshmk96/stk/gsk"
	gskmw "github.com/adharshmk96/stk/pkg/middleware"
)

func RateLimiter() gsk.Middleware {
	rlConfig := gskmw.RateLimiterConfig{
		RequestsPerInterval: 10,
		Interval:            60 * time.Second,
	}
	rateLimiter := gskmw.NewRateLimiter(rlConfig)
	return rateLimiter.Middleware
}
`,
}

var SERVER_INFRA_CONFIGGO_TPL = Template{
	FilePath: "server/infra/config.go",
	Content: `package infra

import "github.com/spf13/viper"

// Configurations are loaded from the environment variables using viper.
// callin this function will reLoad the config. (useful for testing)
// WARN: this will reload all the config.
func LoadDefaultConfig() {
	viper.SetDefault(ENV_SQLITE_FILEPATH, "database.db")

	viper.AutomaticEnv()
}
`,
}

var SERVER_INFRA_CONSTANTSGO_TPL = Template{
	FilePath: "server/infra/constants.go",
	Content: `package infra

// Environment Variable Names
const (
	ENV_SQLITE_FILEPATH = "SQLITE_FILEPATH"
)
`,
}

var SERVER_INFRA_LOGGERGO_TPL = Template{
	FilePath: "server/infra/logger.go",
	Content: `package infra

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func GetLogger() *slog.Logger {
	return logger
}
`,
}

var SERVER_INFRA_DB_SQLITEGO_TPL = Template{
	FilePath: "server/infra/db/sqlite.go",
	Content: `package db

import (
	"database/sql"
	"sync"

	"{{ .PkgName }}/server/infra"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

var (
	sqliteInstance *sql.DB
	sqliteOnce     sync.Once
)

// GetSqliteConnection returns a singleton database connection
func GetSqliteConnection() *sql.DB {
	filepath := viper.GetString(infra.ENV_SQLITE_FILEPATH)
	sqliteOnce.Do(func() {
		db, err := sql.Open("sqlite3", filepath)
		if err != nil {
			panic(err)
		}
		sqliteInstance = db
	})
	return sqliteInstance
}

// ResetSqliteConnection resets the singleton database connection
func ResetSqliteConnection() {
	sqliteInstance = nil
	sqliteOnce = sync.Once{}
}
`,
}

var PUBLIC_ASSETS_SCRIPTJS_TPL = Template{
	FilePath: "public/assets/script.js",
	Content: ``,
}

var PUBLIC_ASSETS_STYLESCSS_TPL = Template{
	FilePath: "public/assets/styles.css",
	Content: ``,
}

var PUBLIC_TEMPLATES_INDEXHTML_TPL = Template{
	FilePath: "public/templates/index.html",
	Content: `<!DOCTYPE html>
<html lang="en">
    <head>
    <meta charset="UTF-8">
    <title>{{ .Var.Title }}</title>
    <link rel="stylesheet" href="{{ .Config.Static }}/style.css">
</head>
<body>
    <h1>{{ .Var.Title }}</h1>
    <p>{{ .Var.Content }}</p>

    <script src="{{ .Config.Static }}/script.js"></script>
</body>
</html>
`,
}

var ProjectTemplates = []Template{
	MAKEFILE_TPL,
	GITIGNORE_TPL,
	MAINGO_TPL,
	READMEMD_TPL,
	CMD_ROOTGO_TPL,
	CMD_SERVEGO_TPL,
	VSCODE_LAUNCHJSON_TPL,
	SERVER_SETUPGO_TPL,
	SERVER_ROUTING_INITROUTESGO_TPL,
	SERVER_MIDDLEWARE_MIDDLEWAREGO_TPL,
	SERVER_INFRA_CONFIGGO_TPL,
	SERVER_INFRA_CONSTANTSGO_TPL,
	SERVER_INFRA_LOGGERGO_TPL,
	SERVER_INFRA_DB_SQLITEGO_TPL,
	PUBLIC_ASSETS_SCRIPTJS_TPL,
	PUBLIC_ASSETS_STYLESCSS_TPL,
	PUBLIC_TEMPLATES_INDEXHTML_TPL,
}
