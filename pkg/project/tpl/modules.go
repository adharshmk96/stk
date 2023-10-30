package tpl

var INTERNALS_CORE_ENTITY_PINGGO_MOD_TPL = Template{
	FilePath: "internals/core/entity/ping.go",
	Render: true,
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

var INTERNALS_CORE_SERR_PINGGO_MOD_TPL = Template{
	FilePath: "internals/core/serr/ping.go",
	Render: true,
	Content: `package serr

import "errors"

var (
	Err{{ .ExportedName }}Failed = errors.New("{{ .ModName }} failed")
)
`,
}

var INTERNALS_HTTP_HANDLER_PINGGO_MOD_TPL = Template{
	FilePath: "internals/http/handler/ping.go",
	Render: true,
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

var INTERNALS_HTTP_HANDLER_TEST_PING_TESTGO_MOD_TPL = Template{
	FilePath: "internals/http/handler_test/ping_test.go",
	Render: true,
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

*/
`,
}

var INTERNALS_HTTP_HELPERS_PINGGO_MOD_TPL = Template{
	FilePath: "internals/http/helpers/ping.go",
	Render: true,
	Content: `package helpers
`,
}

var INTERNALS_HTTP_TRANSPORT_PINGGO_MOD_TPL = Template{
	FilePath: "internals/http/transport/ping.go",
	Render: true,
	Content: `package transport
`,
}

var INTERNALS_SERVICE_PINGGO_MOD_TPL = Template{
	FilePath: "internals/service/ping.go",
	Render: true,
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

var INTERNALS_SERVICE_TEST_PING_TESTGO_MOD_TPL = Template{
	FilePath: "internals/service_test/ping_test.go",
	Render: true,
	Content: `package service_test

// run the following command to generate mocks for {{ .ExportedName }}Storage and {{ .ExportedName }} interfaces
//
// mockery --dir=internals/core/entity --name=^{{ .ExportedName }}.*
//
// and uncomment the following code

/*

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

*/
`,
}

var INTERNALS_STORAGE_PINGSTORAGE_PINGGO_MOD_TPL = Template{
	FilePath: "internals/storage/pingStorage/ping.go",
	Render: true,
	Content: `package {{ .ModName }}Storage

import (
	"database/sql"
	"fmt"

	"{{ .PkgName }}/internals/core/entity"
	"{{ .PkgName }}/internals/core/serr"
	"{{ .PkgName }}/server/infra"
)

type sqliteRepo struct {
	conn *sql.DB
}

func NewSqliteRepo(conn *sql.DB) entity.{{ .ExportedName }}Storage {
	return &sqliteRepo{
		conn: conn,
	}
}

// Repository Methods
func (s *sqliteRepo) {{ .ExportedName }}() error {
	rows, err := s.conn.Query("SELECT 1")
	if err != nil {
		return serr.Err{{ .ExportedName }}Failed
	}
	defer rows.Close()

	var result int

	if rows.Next() {
		err = rows.Scan(&result)
		if err != nil {
			return serr.Err{{ .ExportedName }}Failed
		}
	} else {
		return serr.Err{{ .ExportedName }}Failed
	}

	logger := infra.GetLogger()
	logger.Info(fmt.Sprintf("connection result: %d", result))
	return nil
}
`,
}

var INTERNALS_STORAGE_PINGSTORAGE_PINGQUERIESGO_MOD_TPL = Template{
	FilePath: "internals/storage/pingStorage/pingQueries.go",
	Render: true,
	Content: `package {{ .ModName }}Storage

const (
	SELECT_ONE_TEST = "SELECT 1"
)
`,
}

var SERVER_ROUTING_PINGGO_MOD_TPL = Template{
	FilePath: "server/routing/ping.go",
	Render: true,
	Content: `package routing

import (
	"{{ .PkgName }}/internals/core/entity"
	"{{ .PkgName }}/internals/http/handler"
	"{{ .PkgName }}/internals/service"
	"{{ .PkgName }}/internals/storage/{{ .ModName }}Storage"
	"{{ .PkgName }}/server/infra/db"
	"github.com/adharshmk96/stk/gsk"
)

func initialize{{ .ExportedName }}() entity.{{ .ExportedName }}Handlers {
	conn := db.GetSqliteConnection()

	{{ .ModName }}Storage := {{ .ModName }}Storage.NewSqliteRepo(conn)
	{{ .ModName }}Service := service.New{{ .ExportedName }}Service({{ .ModName }}Storage)
	{{ .ModName }}Handler := handler.New{{ .ExportedName }}Handler({{ .ModName }}Service)

	return {{ .ModName }}Handler
}

func setup{{ .ExportedName }}Routes(rg *gsk.RouteGroup) {
	{{ .ModName }}Handler := initialize{{ .ExportedName }}()

	{{ .ModName }}Routes := rg.RouteGroup("/{{ .ModName }}")

	{{ .ModName }}Routes.Get("/", {{ .ModName }}Handler.{{ .ExportedName }}Handler)
}

func init() {
	routeGroups = append(routeGroups, setup{{ .ExportedName }}Routes)
}
`,
}

var ModuleTemplates = []Template{
	INTERNALS_CORE_ENTITY_PINGGO_MOD_TPL,
	INTERNALS_CORE_SERR_PINGGO_MOD_TPL,
	INTERNALS_HTTP_HANDLER_PINGGO_MOD_TPL,
	INTERNALS_HTTP_HANDLER_TEST_PING_TESTGO_MOD_TPL,
	INTERNALS_HTTP_HELPERS_PINGGO_MOD_TPL,
	INTERNALS_HTTP_TRANSPORT_PINGGO_MOD_TPL,
	INTERNALS_SERVICE_PINGGO_MOD_TPL,
	INTERNALS_SERVICE_TEST_PING_TESTGO_MOD_TPL,
	INTERNALS_STORAGE_PINGSTORAGE_PINGGO_MOD_TPL,
	INTERNALS_STORAGE_PINGSTORAGE_PINGQUERIESGO_MOD_TPL,
	SERVER_ROUTING_PINGGO_MOD_TPL,
}
