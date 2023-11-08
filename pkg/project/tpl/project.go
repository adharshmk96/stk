package tpl

var GITIGNORE_TPL = Template{
	FilePath: ".gitignore",
	Render: true,
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

var GORELEASERYAML_TPL = Template{
	FilePath: ".goreleaser.yaml",
	Render: true,
	Content: `project_name: {{ .AppName }}

before:
  hooks:
    # You may remove this if you don't use go modules.
    - go mod tidy

builds:
  - main: ./main.go
    binary: {{ .AppName }}
    ldflags:
      - -s -w -X "{{ .PkgName }}/cmd.SemVer={{"{{ .Tag }}"}}"
    env:
      - CGO_ENABLED=0
    goos:
      - linux
      - windows
      - darwin
    goarch:
      - amd64
      - arm64

archives:
  - format: tar.gz
    format_overrides:
      - goos: windows
        format: zip
    
    # this name template makes the OS and Arch compatible with the results of uname.
    name_template: "{{"{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"}}"

changelog:
  sort: asc
  use: github
  filters:
    exclude:
      - '^docs:'
      - '^test:'
    include:
      - "^feat:"
      - "^fix:"
      - "^refactor:"
      - "^chore:"
      - "^perf:"


`,
}

var MAINGO_TPL = Template{
	FilePath: "main.go",
	Render: true,
	Content: `package main

import "{{ .PkgName }}/cmd"

func main() {
	cmd.Execute()
}
`,
}

var MAKEFILE_TPL = Template{
	FilePath: "makefile",
	Render: true,
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

var READMEMD_TPL = Template{
	FilePath: "README.md",
	Render: false,
	Content: `# {{ .AppName }}

Run

go run main.go serve -p 8080

- or -

make run


## Project Structure Documentation

This is the project structure generated by "stk init" command.

read more about project structure [here](https://stk-docs.netlify.app/getting-started/project-structure)`,
}

var GITHOOKS_PREPUSH_TPL = Template{
	FilePath: ".githooks/pre-push",
	Render: true,
	Content: `#!/bin/sh
make testci`,
}

var GITHUB_WORKFLOWS_GOBUILDTESTYML_TPL = Template{
	FilePath: ".github/workflows/go-build-test.yml",
	Render: false,
	Content: `# This workflow will build a golang project
# For more information see: https://docs.github.com/en/actions/automating-builds-and-tests/building-and-testing-go

name: Go Build and Test

on:
  push:
    branches: [ "main" ]
  pull_request:
    branches: [ "main" ]

jobs:

  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3

    - name: Build
      run: "make build"

    - name: Test
      run: "make testci"
`,
}

var GITHUB_WORKFLOWS_GORELEASEYML_TPL = Template{
	FilePath: ".github/workflows/go-release.yml",
	Render: false,
	Content: `name: Go Release Workflow

on:
  push:
    # Sequence of patterns matched against refs/tags
    tags:
      - "v*.*.*" # Push events to matching v*, i.e. v1.0, v20.15.10

jobs:
  release:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: 0
      - run: git fetch --force --tags # Ensures go releaser picks us previous tags.
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21' # The Go version to download (if necessary) and use.
     
      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v5
        with:
          distribution: goreleaser
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}`,
}

var VSCODE_LAUNCHJSON_TPL = Template{
	FilePath: ".vscode/launch.json",
	Render: true,
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

var CMD_ROOTGO_TPL = Template{
	FilePath: "cmd/root.go",
	Render: true,
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
var SemVer = "development"

func displaySemverInfo() {
	if SemVer != "development" {
		fmt.Printf("v%s", SemVer)
		return
	}
	version, ok := debug.ReadBuildInfo()
	if ok && version.Main.Version != "(devel)" && version.Main.Version != "" {
		SemVer = version.Main.Version
	}
	fmt.Printf("v%s", SemVer)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "{{ .AppName }}",
	Short: "{{ .AppName }} is an stk project.",
	Long:  "{{ .AppName }} is generated using stk cli.",
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("version").Value.String() == "true" {
			displaySemverInfo()
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
	Render: true,
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

var PUBLIC_ASSETS_SCRIPTJS_TPL = Template{
	FilePath: "public/assets/script.js",
	Render: true,
	Content: ``,
}

var PUBLIC_ASSETS_STYLESCSS_TPL = Template{
	FilePath: "public/assets/styles.css",
	Render: true,
	Content: ``,
}

var PUBLIC_TEMPLATES_INDEXHTML_TPL = Template{
	FilePath: "public/templates/index.html",
	Render: false,
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

var SERVER_SETUPGO_TPL = Template{
	FilePath: "server/setup.go",
	Render: true,
	Content: `package server

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/adharshmk96/stk/gsk"
	"github.com/adharshmk96/stk/pkg/middleware"
	"{{ .PkgName }}/server/infra"
	svrmw "{{ .PkgName }}/server/middleware"
	"{{ .PkgName }}/server/routing"
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
	// once it receives a signal, it will shut down close the done channel
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

var SERVER_INFRA_CONFIGGO_TPL = Template{
	FilePath: "server/infra/config.go",
	Render: true,
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
	Render: true,
	Content: `package infra

// Environment Variable Names
const (
	ENV_SQLITE_FILEPATH = "SQLITE_FILEPATH"
)
`,
}

var SERVER_INFRA_LOGGERGO_TPL = Template{
	FilePath: "server/infra/logger.go",
	Render: true,
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
	Render: true,
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

var SERVER_MIDDLEWARE_RATELIMITERGO_TPL = Template{
	FilePath: "server/middleware/rateLimiter.go",
	Render: true,
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

var SERVER_ROUTING_SETUPROUTESGO_TPL = Template{
	FilePath: "server/routing/setupRoutes.go",
	Render: true,
	Content: `package routing

import (
	"github.com/adharshmk96/stk/gsk"
)

var webRouteGroups = []func(*gsk.RouteGroup){}
var apiRouteGroups = []func(*gsk.RouteGroup){}

func RegisterApiRoutes(routeGroup func(*gsk.RouteGroup)) {
	apiRouteGroups = append(apiRouteGroups, routeGroup)
}

func RegisterWebRoutes(routeGroup func(*gsk.RouteGroup)) {
	webRouteGroups = append(webRouteGroups, routeGroup)
}

func SetupApiRoutes(server *gsk.Server) {
	apiRoutes := server.RouteGroup("/api")

	for _, routeGroup := range apiRouteGroups {
		routeGroup(apiRoutes)
	}
}

func SetupTemplateRoutes(server *gsk.Server) {
	templateRoutes := server.RouteGroup("/")

	for _, routeGroup := range webRouteGroups {
		routeGroup(templateRoutes)
	}

}
`,
}

var ProjectTemplates = []Template{
	GITIGNORE_TPL,
	GORELEASERYAML_TPL,
	MAINGO_TPL,
	MAKEFILE_TPL,
	READMEMD_TPL,
	GITHOOKS_PREPUSH_TPL,
	GITHUB_WORKFLOWS_GOBUILDTESTYML_TPL,
	GITHUB_WORKFLOWS_GORELEASEYML_TPL,
	VSCODE_LAUNCHJSON_TPL,
	CMD_ROOTGO_TPL,
	CMD_SERVEGO_TPL,
	PUBLIC_ASSETS_SCRIPTJS_TPL,
	PUBLIC_ASSETS_STYLESCSS_TPL,
	PUBLIC_TEMPLATES_INDEXHTML_TPL,
	SERVER_SETUPGO_TPL,
	SERVER_INFRA_CONFIGGO_TPL,
	SERVER_INFRA_CONSTANTSGO_TPL,
	SERVER_INFRA_LOGGERGO_TPL,
	SERVER_INFRA_DB_SQLITEGO_TPL,
	SERVER_MIDDLEWARE_RATELIMITERGO_TPL,
	SERVER_ROUTING_SETUPROUTESGO_TPL,
}
