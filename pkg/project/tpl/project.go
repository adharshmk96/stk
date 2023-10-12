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

var DATABASEDB_TPL = Template{
	FilePath: "database.db",
	Content: ``,
}

var MAINGO_TPL = Template{
	FilePath: "main.go",
	Content: `package main

import "github.com/adharshmk96/stk-template/{{ .AppName }}/cmd"

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

run: 
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
	Content: `GET http://localhost:8080/{{ .ModName }}`,
}

var SQLITEDB_TPL = Template{
	FilePath: "sqlite.db",
	Content: ``,
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

var CMD_SERVEGO_TPL = Template{
	FilePath: "cmd/serve.go",
	Content: `package cmd

import (
	"sync"

	"github.com/adharshmk96/stk-template/{{ .AppName }}/server"
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

var INTERNALS_CORE_ENTITY_PINGGO_TPL = Template{
	FilePath: "internals/core/entity/ping.go",
	Content: `package entity

import "github.com/adharshmk96/stk/gsk"

// Domain
type {{ .ExportedName }}Data struct {
	{{ .ModName }} string
}

// Storage
type {{ .ExportedName }}Storage interface {
	{{ .ExportedName }}() error
}

// Service
type {{ .ExportedName }}Service interface {
	{{ .ExportedName }}Service() (string, error)
}

// Handler
type {{ .ExportedName }}Handlers interface {
	{{ .ExportedName }}Handler(gc *gsk.Context)
}
`,
}

var INTERNALS_CORE_SERR_PINGGO_TPL = Template{
	FilePath: "internals/core/serr/ping.go",
	Content: `package serr

import "errors"

var (
	Err{{ .ExportedName }}Failed = errors.New("{{ .ModName }} failed")
)
`,
}

var INTERNALS_HTTP_HANDLER_PINGGO_TPL = Template{
	FilePath: "internals/http/handler/ping.go",
	Content: `package handler

import (
	"net/http"

	"github.com/adharshmk96/stk-template/{{ .AppName }}/internals/core/entity"
	"github.com/adharshmk96/stk/gsk"
)

type {{ .ModName }}Handler struct {
	service entity.{{ .ExportedName }}Service
}

func New{{ .ExportedName }}Handler(service entity.{{ .ExportedName }}Service) entity.{{ .ExportedName }}Handlers {
	return &{{ .ModName }}Handler{
		service: service,
	}
}

/*
{{ .ExportedName }}Handler returns {{ .ModName }} 200 response
Response:
- 200: OK
- 500: Internal Server Error
*/
func (h *{{ .ModName }}Handler) {{ .ExportedName }}Handler(gc *gsk.Context) {

	{{ .ModName }}, err := h.service.{{ .ExportedName }}Service()
	if err != nil {
		gc.Status(http.StatusInternalServerError).JSONResponse(gsk.Map{
			"error": err.Error(),
		})
		return
	}

	gc.Status(http.StatusOK).JSONResponse(gsk.Map{
		"message": {{ .ModName }},
	})
}
`,
}

var INTERNALS_HTTP_HANDLER_PING_TESTGO_TPL = Template{
	FilePath: "internals/http/handler/ping_test.go",
	Content: `package handler_test

// run the following command to generate mocks for {{ .ExportedName }} interfaces
//
// mockery --dir=internals/core/entity --name=^{{ .ExportedName }}.*
//
// and uncomment the following code

/*

import (
	"net/http"
	"testing"

	"github.com/adharshmk96/stk-template/{{ .AppName }}/internals/http/handler"
	"github.com/adharshmk96/stk-template/{{ .AppName }}/mocks"
	"github.com/adharshmk96/stk/gsk"
	"github.com/stretchr/testify/assert"
)

func Test{{ .ExportedName }}Handler(t *testing.T) {
	t.Run("{{ .ExportedName }} Handler returns 200", func(t *testing.T) {

		// Arrange
		s := gsk.New()
		service := mocks.New{{ .ExportedName }}Service(t)
		service.On("{{ .ExportedName }}Service").Return("pong", nil)

		{{ .ModName }}Handler := handler.New{{ .ExportedName }}Handler(service)

		s.Get("/{{ .ModName }}", {{ .ModName }}Handler.{{ .ExportedName }}Handler)

		// Act
		w, _ := s.Test("GET", "/{{ .ModName }}", nil)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

*/
`,
}

var INTERNALS_HTTP_HELPERS_PINGGO_TPL = Template{
	FilePath: "internals/http/helpers/ping.go",
	Content: `package helpers
`,
}

var INTERNALS_HTTP_TRANSPORT_PINGGO_TPL = Template{
	FilePath: "internals/http/transport/ping.go",
	Content: `package transport
`,
}

var INTERNALS_SERVICE_PINGGO_TPL = Template{
	FilePath: "internals/service/ping.go",
	Content: `package service

import "github.com/adharshmk96/stk-template/{{ .AppName }}/internals/core/entity"

type {{ .ModName }}Service struct {
	storage entity.{{ .ExportedName }}Storage
}

func New{{ .ExportedName }}Service(storage entity.{{ .ExportedName }}Storage) entity.{{ .ExportedName }}Service {
	return &{{ .ModName }}Service{
		storage: storage,
	}
}

func (s *{{ .ModName }}Service) {{ .ExportedName }}Service() (string, error) {
	err := s.storage.{{ .ExportedName }}()
	if err != nil {
		return "", err
	}
	return "pong", nil
}
`,
}

var INTERNALS_SERVICE_PING_TESTGO_TPL = Template{
	FilePath: "internals/service/ping_test.go",
	Content: `package service_test

// run the following command to generate mocks for {{ .ExportedName }}Storage and {{ .ExportedName }} interfaces
//
// mockery --dir=internals/core/entity --name=^{{ .ExportedName }}.*
//
// and uncomment the following code

/*

import (
	"testing"

	"github.com/adharshmk96/stk-template/{{ .AppName }}/internals/service"
	"github.com/adharshmk96/stk-template/{{ .AppName }}/mocks"
	"github.com/stretchr/testify/assert"
)

func Test{{ .ExportedName }}Service(t *testing.T) {
	t.Run("{{ .ExportedName }}Service returns pong", func(t *testing.T) {

		// Arrange
		storage := mocks.New{{ .ExportedName }}Storage(t)
		storage.On("{{ .ExportedName }}").Return(nil)

		svc := service.New{{ .ExportedName }}Service(storage)

		// Act
		msg, err := svc.{{ .ExportedName }}Service()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "pong", msg)
	})
}

*/
`,
}

var INTERNALS_STORAGE_PINGSTORAGE_PINGGO_TPL = Template{
	FilePath: "internals/storage/pingStorage/ping.go",
	Content: `package {{ .ModName }}Storage

import (
	"fmt"

	"github.com/adharshmk96/stk-template/multimod/server/infra"
	"github.com/adharshmk96/stk-template/{{ .AppName }}/internals/core/serr"
)

// Repository Methods
func (s *sqliteRepo) {{ .ExportedName }}() error {
	res, err := s.conn.Exec("SELECT 1")
	if err != nil {
		return serr.Err{{ .ExportedName }}Failed
	}
	num, err := res.RowsAffected()
	if err != nil {
		return serr.Err{{ .ExportedName }}Failed
	}

	logger := infra.GetLogger()
	logger.Info(fmt.Sprintf("{{ .ExportedName }} Success: %d", num))
	return nil
}
`,
}

var INTERNALS_STORAGE_PINGSTORAGE_PINGCONNECTIONGO_TPL = Template{
	FilePath: "internals/storage/pingStorage/pingConnection.go",
	Content: `package {{ .ModName }}Storage

import (
	"database/sql"

	"github.com/adharshmk96/stk-template/{{ .AppName }}/internals/core/entity"
)

type sqliteRepo struct {
	conn *sql.DB
}

func NewSqliteRepo(conn *sql.DB) entity.{{ .ExportedName }}Storage {
	return &sqliteRepo{
		conn: conn,
	}
}
`,
}

var INTERNALS_STORAGE_PINGSTORAGE_PINGQUERIESGO_TPL = Template{
	FilePath: "internals/storage/pingStorage/pingQueries.go",
	Content: `package {{ .ModName }}Storage

const (
	SELECT_ONE_TEST = "SELECT 1"
)
`,
}

var SERVER_SETUPGO_TPL = Template{
	FilePath: "server/setup.go",
	Content: `package server

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/adharshmk96/stk-template/{{ .AppName }}/server/infra"
	svrmw "github.com/adharshmk96/stk-template/{{ .AppName }}/server/middleware"
	"github.com/adharshmk96/stk-template/{{ .AppName }}/server/routing"
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

	routing.SetupRoutes(server)

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

var SERVER_ROUTING_INITROUTESGO_TPL = Template{
	FilePath: "server/routing/initRoutes.go",
	Content: `package routing

import (
	"github.com/adharshmk96/stk/gsk"
)

func SetupRoutes(server *gsk.Server) {
	setup{{ .ExportedName }}Routes(server)
}
`,
}

var SERVER_ROUTING_PINGGO_TPL = Template{
	FilePath: "server/routing/ping.go",
	Content: `package routing

import (
	"github.com/adharshmk96/stk-template/{{ .AppName }}/internals/http/handler"
	"github.com/adharshmk96/stk-template/{{ .AppName }}/internals/service"
	"github.com/adharshmk96/stk-template/{{ .AppName }}/internals/storage/{{ .ModName }}Storage"
	"github.com/adharshmk96/stk-template/{{ .AppName }}/server/infra"
	"github.com/adharshmk96/stk/gsk"
	"github.com/adharshmk96/stk/pkg/db"
	"github.com/spf13/viper"
)

func setup{{ .ExportedName }}Routes(server *gsk.Server) {
	dbConfig := viper.GetString(infra.ENV_SQLITE_FILEPATH)
	conn := db.GetSqliteConnection(dbConfig)

	{{ .ModName }}Storage := {{ .ModName }}Storage.NewSqliteRepo(conn)
	{{ .ModName }}Service := service.New{{ .ExportedName }}Service({{ .ModName }}Storage)
	{{ .ModName }}Handler := handler.New{{ .ExportedName }}Handler({{ .ModName }}Service)

	server.Get("/{{ .ModName }}", {{ .ModName }}Handler.{{ .ExportedName }}Handler)
}
`,
}

var ProjectTemplates = []Template{
	GITIGNORE_TPL,
	DATABASEDB_TPL,
	MAINGO_TPL,
	MAKEFILE_TPL,
	READMEMD_TPL,
	REQUESTHTTP_TPL,
	SQLITEDB_TPL,
	VSCODE_LAUNCHJSON_TPL,
	CMD_ROOTGO_TPL,
	CMD_SERVEGO_TPL,
	CMD_VERSIONGO_TPL,
	INTERNALS_CORE_ENTITY_PINGGO_TPL,
	INTERNALS_CORE_SERR_PINGGO_TPL,
	INTERNALS_HTTP_HANDLER_PINGGO_TPL,
	INTERNALS_HTTP_HANDLER_PING_TESTGO_TPL,
	INTERNALS_HTTP_HELPERS_PINGGO_TPL,
	INTERNALS_HTTP_TRANSPORT_PINGGO_TPL,
	INTERNALS_SERVICE_PINGGO_TPL,
	INTERNALS_SERVICE_PING_TESTGO_TPL,
	INTERNALS_STORAGE_PINGSTORAGE_PINGGO_TPL,
	INTERNALS_STORAGE_PINGSTORAGE_PINGCONNECTIONGO_TPL,
	INTERNALS_STORAGE_PINGSTORAGE_PINGQUERIESGO_TPL,
	SERVER_SETUPGO_TPL,
	SERVER_INFRA_CONFIGGO_TPL,
	SERVER_INFRA_CONSTANTSGO_TPL,
	SERVER_INFRA_LOGGERGO_TPL,
	SERVER_MIDDLEWARE_MIDDLEWAREGO_TPL,
	SERVER_ROUTING_INITROUTESGO_TPL,
	SERVER_ROUTING_PINGGO_TPL,
}
