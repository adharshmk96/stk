package tpl

var INTERNALSCOREENTITYPINGGO_MOD = Template{
	FilePath: "internals/core/entity/{{ .ModName }}.go",
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

var INTERNALSCORESERRPINGGO_MOD = Template{
	FilePath: "internals/core/serr/{{ .ModName }}.go",
	Content: `package serr

import "errors"

var (
	Err{{ .ExportedName }}Failed = errors.New("{{ .ModName }} failed")
)
`,
}

var INTERNALSHTTPHANDLERPINGGO_MOD = Template{
	FilePath: "internals/http/handler/{{ .ModName }}.go",
	Content: `package handler

import (
	"net/http"

	"{{ .PkgName }}/internals/core/entity"
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

var INTERNALSHTTPHANDLERPING_TESTGO_MOD = Template{
	FilePath: "internals/http/handler/{{ .ModName }}_test.go",
	Content: `package handler_test

import (
	"net/http"
	"testing"

	"{{ .PkgName }}/internals/http/handler"
	"{{ .PkgName }}/mocks"
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
`,
}

var INTERNALSHTTPHELPERSPINGGO_MOD = Template{
	FilePath: "internals/http/helpers/{{ .ModName }}.go",
	Content: `package helpers
`,
}

var INTERNALSHTTPTRANSPORTPINGGO_MOD = Template{
	FilePath: "internals/http/transport/{{ .ModName }}.go",
	Content: `package transport
`,
}

var INTERNALSSERVICEPINGGO_MOD = Template{
	FilePath: "internals/service/{{ .ModName }}.go",
	Content: `package service

import "{{ .PkgName }}/internals/core/entity"

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

var INTERNALSSERVICEPING_TESTGO_MOD = Template{
	FilePath: "internals/service/{{ .ModName }}_test.go",
	Content: `package service_test

import (
	"testing"

	"{{ .PkgName }}/internals/service"
	"{{ .PkgName }}/mocks"
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
`,
}

var INTERNALSSTORAGEPINGSTORAGEPINGGO_MOD = Template{
	FilePath: "internals/storage/{{ .ModName }}Storage/{{ .ModName }}.go",
	Content: `package {{ .ModName }}Storage

import "{{ .PkgName }}/internals/core/serr"

// Repository Methods
func (s *sqliteRepo) {{ .ExportedName }}() error {
	err := s.conn.{{ .ExportedName }}()
	if err != nil {
		return serr.Err{{ .ExportedName }}Failed
	}
	return nil
}
`,
}

var INTERNALSSTORAGEPINGSTORAGEPINGCONNECTIONGO_MOD = Template{
	FilePath: "internals/storage/{{ .ModName }}Storage/{{ .ModName }}Connection.go",
	Content: `package {{ .ModName }}Storage

import (
	"database/sql"

	"{{ .PkgName }}/internals/core/entity"
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

var INTERNALSSTORAGEPINGSTORAGEPINGQUERIESGO_MOD = Template{
	FilePath: "internals/storage/{{ .ModName }}Storage/{{ .ModName }}Queries.go",
	Content: `package {{ .ModName }}Storage

const (
	SELECT_ONE_TEST = "SELECT 1"
)
`,
}

var SERVERROUTINGPINGGO_MOD = Template{
	FilePath: "server/routing/{{ .ModName }}.go",
	Content: `package routing

import (
	"{{ .PkgName }}/internals/http/handler"
	"{{ .PkgName }}/internals/service"
	"{{ .PkgName }}/internals/storage/{{ .ModName }}Storage"
	"{{ .PkgName }}/server/infra"
	"github.com/adharshmk96/stk/gsk"
	"github.com/adharshmk96/stk/pkg/db"
	"github.com/spf13/viper"
)

func setup{{ .ExportedName }}Routes(server *gsk.Server) {
	dbConfig := viper.GetString(infra.ENV_SQLITE_FILEPATH)
	conn := db.GetSqliteConnection(dbConfig)

	{{ .AppName }}Storage := {{ .ModName }}Storage.NewSqliteRepo(conn)
	{{ .AppName }}Service := service.New{{ .ExportedName }}Service({{ .AppName }}Storage)
	{{ .AppName }}Handler := handler.New{{ .ExportedName }}Handler({{ .AppName }}Service)

	server.Get("/{{ .ModName }}", {{ .AppName }}Handler.{{ .ExportedName }}Handler)
}
`,
}

var ModuleTemplates = []Template{
	INTERNALSCOREENTITYPINGGO_MOD,
	INTERNALSCORESERRPINGGO_MOD,
	INTERNALSHTTPHANDLERPINGGO_MOD,
	INTERNALSHTTPHANDLERPING_TESTGO_MOD,
	INTERNALSHTTPHELPERSPINGGO_MOD,
	INTERNALSHTTPTRANSPORTPINGGO_MOD,
	INTERNALSSERVICEPINGGO_MOD,
	INTERNALSSERVICEPING_TESTGO_MOD,
	INTERNALSSTORAGEPINGSTORAGEPINGGO_MOD,
	INTERNALSSTORAGEPINGSTORAGEPINGCONNECTIONGO_MOD,
	INTERNALSSTORAGEPINGSTORAGEPINGQUERIESGO_MOD,
	SERVERROUTINGPINGGO_MOD,
}
