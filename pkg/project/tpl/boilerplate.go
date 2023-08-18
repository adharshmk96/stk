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

var SERVER_SETUPGO_TPL = Template{
	FilePath: "server/setup.go",
	Content: `package server

import (
	"os"
	"os/signal"
	"syscall"

	"{{ .PkgName }}/pkg/http/handler"
	"{{ .PkgName }}/pkg/service"
	"{{ .PkgName }}/pkg/storage/sqlite"
	"{{ .PkgName }}/server/infra"
	svrmw "{{ .PkgName }}/server/middleware"
	"{{ .PkgName }}/server/routing"
	"github.com/adharshmk96/stk/gsk"
	"github.com/adharshmk96/stk/pkg/db"
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

	intializeServer(server)

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
			logger.Error(err)
		}

		close(done)
	}()

	return server, done
}

func intializeServer(server *gsk.Server) {
	conn := db.GetSqliteConnection("sqlite.db")

	{{ .AppName }}Storage := sqlite.NewSqliteRepo(conn)
	{{ .AppName }}Service := service.NewPingService({{ .AppName }}Storage)
	{{ .AppName }}Handler := handler.NewPingHandler({{ .AppName }}Service)

	routing.SetupPingRoutes(server, {{ .AppName }}Handler)
}
`,
}

var SERVER_MIDDLEWARE_MIDDLEWAREGO_TPL = Template{
	FilePath: "server/middleware/middleware.go",
	Content: `package server

import (
	"time"

	"github.com/adharshmk96/stk/gsk"
	"github.com/adharshmk96/stk/pkg/middleware"
)

func RateLimiter() gsk.Middleware {
	rlConfig := middleware.RateLimiterConfig{
		RequestsPerInterval: 10,
		Interval:            60 * time.Second,
	}
	rateLimiter := middleware.NewRateLimiter(rlConfig)
	return rateLimiter.Middleware
}
`,
}

var SERVER_ROUTING_ROUTINGGO_TPL = Template{
	FilePath: "server/routing/routing.go",
	Content: `package routing

import (
	"{{ .PkgName }}/pkg/core"
	"github.com/adharshmk96/stk/gsk"
)

func SetupPingRoutes(server *gsk.Server, {{ .AppName }}Handler core.PingHandlers) {
	server.Get("/ping", {{ .AppName }}Handler.PingHandler)
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

var PKG_SERVICE_SERVICEGO_TPL = Template{
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

var PKG_SERVICE_PINGGO_TPL = Template{
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

var PKG_STORAGE_SQLITE_SQLITEGO_TPL = Template{
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

var PKG_STORAGE_SQLITE_PINGGO_TPL = Template{
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

var PKG_HTTP_HANDLER_HANDLERGO_TPL = Template{
	FilePath: "pkg/http/handler/handler.go",
	Content: `package handler

import (
	"{{ .PkgName }}/pkg/core"
)

type pingHandler struct {
	pingService core.PingService
}

func NewPingHandler(pingService core.PingService) core.PingHandlers {
	return &pingHandler{
		pingService: pingService,
	}
}
`,
}

var PKG_HTTP_HANDLER_PINGGO_TPL = Template{
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

var PKG_CORE_SERVICEGO_TPL = Template{
	FilePath: "pkg/core/service.go",
	Content: `package core

type PingService interface {
	PingService() string
}
`,
}

var PKG_CORE_HANDLERGO_TPL = Template{
	FilePath: "pkg/core/handler.go",
	Content: `package core

import "github.com/adharshmk96/stk/gsk"

type PingHandlers interface {
	PingHandler(gc *gsk.Context)
}
`,
}

var PKG_CORE_STORAGEGO_TPL = Template{
	FilePath: "pkg/core/storage.go",
	Content: `package core

type PingStorage interface {
	Ping() error
}
`,
}

var PKG_CORE_SERR_PINGERRGO_TPL = Template{
	FilePath: "pkg/core/serr/pingerr.go",
	Content: `package serr

import "errors"

var (
	ErrPingFailed = errors.New("ping failed")
)
`,
}

var PKG_CORE_DS_PINGGO_TPL = Template{
	FilePath: "pkg/core/ds/ping.go",
	Content: `package ds

type User struct {
	pong string
}
`,
}

var CMD_VERSIONGO_TPL = Template{
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

var CMD_ROOTGO_TPL = Template{
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

var BoilerPlateTemplates = []Template{
	GITIGNORE_TPL,
	MAINGO_TPL,
	READMEMD_TPL,
	REQUESTHTTP_TPL,
	MAKEFILE_TPL,
	SERVER_SETUPGO_TPL,
	SERVER_MIDDLEWARE_MIDDLEWAREGO_TPL,
	SERVER_ROUTING_ROUTINGGO_TPL,
	SERVER_INFRA_CONFIGGO_TPL,
	SERVER_INFRA_CONSTANTSGO_TPL,
	SERVER_INFRA_LOGGERGO_TPL,
	PKG_SERVICE_SERVICEGO_TPL,
	PKG_SERVICE_PINGGO_TPL,
	PKG_STORAGE_SQLITE_SQLITEGO_TPL,
	PKG_STORAGE_SQLITE_PINGGO_TPL,
	PKG_HTTP_HANDLER_HANDLERGO_TPL,
	PKG_HTTP_HANDLER_PINGGO_TPL,
	PKG_CORE_SERVICEGO_TPL,
	PKG_CORE_HANDLERGO_TPL,
	PKG_CORE_STORAGEGO_TPL,
	PKG_CORE_SERR_PINGERRGO_TPL,
	PKG_CORE_DS_PINGGO_TPL,
	CMD_VERSIONGO_TPL,
	CMD_SERVEGO_TPL,
	CMD_ROOTGO_TPL,
}
