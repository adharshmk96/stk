package tpl

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

var MAKEFILE_TPL = Template{
	FilePath: "makefile",
	Content: `##########################
### Version Commands
##########################

patch:
	$(eval NEW_TAG := $(shell git semver patch --dryrun))
	$(call update_file)
	@git semver patch

minor:
	$(eval NEW_TAG := $(shell git semver minor --dryrun))
	$(call update_file)
	@git semver minor

major:
	$(eval NEW_TAG := $(shell git semver major --dryrun))
	$(call update_file)
	@git semver major

publish:
	@git push origin $(shell git semver get)


##########################
### Build Commands
##########################

BINARY_NAME=app

build:
	@go build -o ./out/$(BINARY_NAME) -v

run: build
	@go run . serve -p 8080

test:
	@go test ./... -coverprofile=coverage.out

coverage:
	@go test -v ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

testci:
	@go test ./... -coverprofile=coverage.out

clean:
	@go clean
	@rm -f ./out/$(BINARY_NAME)
	@rm -f coverage.out
	@rm -rf .keys
	@rm -f auth_database.db

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

	
##########################
### Helpers
##########################

define update_file
    @echo "updating files to version $(NEW_TAG)"
    @sed -i.bak "s/var version = \"[^\"]*\"/var version = \"$(NEW_TAG)\"/g" ./cmd/root.go
    @rm cmd/root.go.bak
    @git add cmd/root.go
    @git commit -m "bump version to $(NEW_TAG)" > /dev/null
endef

##########################
### Setup Commands
##########################

init: deps keygen initgithooks mockgen
	@echo "Project initialized."

initci: deps keygen
	@echo "Project initialized for CI."

initgithooks:
	@git config core.hooksPath .githooks

mockgen:
	@rm -rf ./mocks
	@mockery --all	

`,
}

var READMEMD_TPL = Template{
	FilePath: "README.md",
	Content: `# templates-for-go
repo with some files for go 
`,
}

var REQUESTHTTP_TPL = Template{
	FilePath: "request.http",
	Content:  `GET http://localhost:8080/ping`,
}

var CMDROOTGO_TPL = Template{
	FilePath: "cmd/root.go",
	Content: `package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version = "v0.0.0"
var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "{{ .AppName }}",
	Short: "{{ .AppName }} is a template for creating api servers.",
	Long:  "{{ .AppName }} is a template for creating api servers using stk.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "{{ .AppName }}.yaml", "config file.")
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

var CMDSERVEGO_TPL = Template{
	FilePath: "cmd/serve.go",
	Content: `package cmd

import (
	"github.com/spf13/cobra"
	"{{ .PkgName }}/server"
)

var startingPort string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		startAddr := "0.0.0.0:"
		server.StartServer(startAddr + startingPort)
	},
}

func init() {
	serveCmd.Flags().StringVarP(&startingPort, "port", "p", "8080", "Port to start the server on")

	rootCmd.AddCommand(serveCmd)
}
`,
}

var CMDVERSIONGO_TPL = Template{
	FilePath: "cmd/version.go",
	Content: `package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "display the version of {{ .AppName }}",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("{{ .AppName }} version: %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
`,
}

var PKGCOREHANDLERGO_TPL = Template{
	FilePath: "pkg/core/handler.go",
	Content: `package core

import "github.com/adharshmk96/stk/gsk"

type PingHandler interface {
	PingHandler(gc *gsk.Context)
}
`,
}

var PKGCORESERVICEGO_TPL = Template{
	FilePath: "pkg/core/service.go",
	Content: `package core

type PingService interface {
	PingService() string
}
`,
}

var PKGCORESTORAGEGO_TPL = Template{
	FilePath: "pkg/core/storage.go",
	Content: `package core

type PingStorage interface {
	Ping() error
}
`,
}

var PKGCOREDSPINGGO_TPL = Template{
	FilePath: "pkg/core/ds/ping.go",
	Content: `package ds

type User struct {
	pong string
}
`,
}

var PKGCORESERRPINGERRGO_TPL = Template{
	FilePath: "pkg/core/serr/pingerr.go",
	Content: `package serr

import "errors"

var (
	ErrPingFailed = errors.New("ping failed")
)
`,
}

var PKGHTTPHANDLERHANDLERGO_TPL = Template{
	FilePath: "pkg/http/handler/handler.go",
	Content: `package handler

import (
	"{{ .PkgName }}/pkg/core"
)

type pingHandler struct {
	pingService core.PingService
}

func NewPingHandler(pingService core.PingService) core.PingHandler {
	return &pingHandler{
		pingService: pingService,
	}
}
`,
}

var PKGHTTPHANDLERPINGGO_TPL = Template{
	FilePath: "pkg/http/handler/ping.go",
	Content: `package handler

import (
	"net/http"

	"github.com/adharshmk96/stk/gsk"
)

/*
PingHandler returns ping 200 response
Response:
- 200: OK
- 500: Internal Server Error
*/
func (h *pingHandler) PingHandler(gc *gsk.Context) {
	
	ping := h.pingService.PingService()

	gc.Status(http.StatusOK).JSONResponse(gsk.Map{
		"message": ping,
	})
}	
`,
}

var PKGSERVICEPINGGO_TPL = Template{
	FilePath: "pkg/service/ping.go",
	Content: `package service

func (s *pingService) PingService() string {
	err := s.pingStorage.Ping()
	if err != nil {
		return "error"
	}
	return "pong"
}
`,
}

var PKGSERVICESERVICEGO_TPL = Template{
	FilePath: "pkg/service/service.go",
	Content: `package service

import (
	"{{ .PkgName }}/pkg/core"
)

type pingService struct {
	pingStorage core.PingStorage
}

func NewPingService(storage core.PingStorage) core.PingService {
	return &pingService{
		pingStorage: storage,
	}
}
`,
}

var PKGSTORAGESQLITEPINGGO_TPL = Template{
	FilePath: "pkg/storage/sqlite/ping.go",
	Content: `package sqlite

import "{{ .PkgName }}/pkg/core/serr"

func (s *sqliteRepo) Ping() error {
	err := s.conn.Ping()
	if err != nil {
		return serr.ErrPingFailed
	}
	return nil
}
`,
}

var PKGSTORAGESQLITESQLITEGO_TPL = Template{
	FilePath: "pkg/storage/sqlite/sqlite.go",
	Content: `package sqlite

import (
	"database/sql"

	"{{ .PkgName }}/pkg/core"
)

type sqliteRepo struct {
	conn *sql.DB
}

func NewSqliteRepo(conn *sql.DB) core.PingStorage {
	return &sqliteRepo{
		conn: conn,
	}
}
`,
}

var SERVERMIDDLEWAREGO_TPL = Template{
	FilePath: "server/middleware.go",
	Content: `package server

import (
	"time"

	"github.com/adharshmk96/stk/gsk"
	"github.com/adharshmk96/stk/pkg/middleware"
)

func rateLimiter() gsk.Middleware {
	rlConfig := middleware.RateLimiterConfig{
		RequestsPerInterval: 10,
		Interval:            60 * time.Second,
	}
	rateLimiter := middleware.NewRateLimiter(rlConfig)
	return rateLimiter.Middleware
}
`,
}

var SERVERROUTINGGO_TPL = Template{
	FilePath: "server/routing.go",
	Content: `package server

import (
	"{{ .PkgName }}/pkg/http/handler"
	"{{ .PkgName }}/pkg/service"
	"{{ .PkgName }}/pkg/storage/sqlite"
	"github.com/adharshmk96/stk/gsk"
	"github.com/adharshmk96/stk/pkg/db"
)

func setupRoutes(server *gsk.Server) {

	conn := db.GetSqliteConnection("sqlite.db")

	{{ .AppName }}Storage := sqlite.NewSqliteRepo(conn)
	{{ .AppName }}Service := service.NewPingService({{ .AppName }}Storage)
	{{ .AppName }}Handler := handler.NewPingHandler({{ .AppName }}Service)

	server.Get("/ping", {{ .AppName }}Handler.PingHandler)
}
`,
}

var SERVERSETUPGO_TPL = Template{
	FilePath: "server/setup.go",
	Content: `package server

import "github.com/adharshmk96/stk/gsk"

func StartServer(port string) *gsk.Server {

	serverConfig := &gsk.ServerConfig{
		Port: port,
	}

	server := gsk.New(serverConfig)

	setupRoutes(server)

	rateLimiter := rateLimiter()
	server.Use(rateLimiter)

	server.Start()

	return server
}
`,
}

var BoilerPlateTemplates = []Template{
	GITIGNORE_TPL,
	MAINGO_TPL,
	MAKEFILE_TPL,
	READMEMD_TPL,
	REQUESTHTTP_TPL,
	CMDROOTGO_TPL,
	CMDSERVEGO_TPL,
	CMDVERSIONGO_TPL,
	PKGCOREHANDLERGO_TPL,
	PKGCORESERVICEGO_TPL,
	PKGCORESTORAGEGO_TPL,
	PKGCOREDSPINGGO_TPL,
	PKGCORESERRPINGERRGO_TPL,
	PKGHTTPHANDLERHANDLERGO_TPL,
	PKGHTTPHANDLERPINGGO_TPL,
	PKGSERVICEPINGGO_TPL,
	PKGSERVICESERVICEGO_TPL,
	PKGSTORAGESQLITEPINGGO_TPL,
	PKGSTORAGESQLITESQLITEGO_TPL,
	SERVERMIDDLEWAREGO_TPL,
	SERVERROUTINGGO_TPL,
	SERVERSETUPGO_TPL,
}
