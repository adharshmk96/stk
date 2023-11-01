package tpl

var INTERNALS_PING_API_HANDLER_PINGGO_MOD_TPL = Template{
	FilePath: "internals/ping/api/handler/ping.go",
	Render: true,
	Content: `package handler

import (
	"net/http"

	"{{ .PkgName }}template/internals/{{ .ModName }}/domain"

	"{{ .PkgName }}/gsk"
)

type {{ .ModName }}Handler struct {
	service domain.{{ .ExportedName }}Service
}

func New{{ .ExportedName }}Handler(service domain.{{ .ExportedName }}Service) domain.{{ .ExportedName }}Handlers {
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

	message, err := h.service.{{ .ExportedName }}Service()
	if err != nil {
		gc.Status(http.StatusInternalServerError).JSONResponse(gsk.Map{
			"error": err.Error(),
		})
		return
	}

	gc.Status(http.StatusOK).JSONResponse(gsk.Map{
		"message": message,
	})
}
`,
}

var INTERNALS_PING_API_HANDLER_PING_TESTGO_MOD_TPL = Template{
	FilePath: "internals/ping/api/handler/ping_test.go",
	Render: true,
	Content: `package handler_test

// run the following command to generate mocks for {{ .ExportedName }} interfaces
//
// mockery --dir=internals/{{ .ModName }}/{{ .ModName }} --name=^{{ .ExportedName }}.*

import (
	"net/http"
	"testing"

	"{{ .PkgName }}template/internals/{{ .ModName }}/api/handler"

	"{{ .PkgName }}/gsk"
	"{{ .PkgName }}template/mocks"
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

var INTERNALS_PING_API_TRANSPORT_PINGGO_MOD_TPL = Template{
	FilePath: "internals/ping/api/transport/ping.go",
	Render: true,
	Content: `package transport
`,
}

var INTERNALS_PING_DOMAIN_PINGGO_MOD_TPL = Template{
	FilePath: "internals/ping/domain/ping.go",
	Render: true,
	Content: `package domain

import "{{ .PkgName }}/gsk"

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

var INTERNALS_PING_SERR_PINGGO_MOD_TPL = Template{
	FilePath: "internals/ping/serr/ping.go",
	Render: true,
	Content: `package serr

import "errors"

var (
	Err{{ .ExportedName }}Failed = errors.New("{{ .ModName }} failed")
)
`,
}

var INTERNALS_PING_SERVICE_PINGGO_MOD_TPL = Template{
	FilePath: "internals/ping/service/ping.go",
	Render: true,
	Content: `package service

import (
	"{{ .PkgName }}template/internals/{{ .ModName }}/domain"
)

type {{ .ModName }}Service struct {
	storage domain.{{ .ExportedName }}Storage
}

func New{{ .ExportedName }}Service(storage domain.{{ .ExportedName }}Storage) domain.{{ .ExportedName }}Service {
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

var INTERNALS_PING_SERVICE_PING_TESTGO_MOD_TPL = Template{
	FilePath: "internals/ping/service/ping_test.go",
	Render: true,
	Content: `package service_test

// run the following command to generate mocks for {{ .ExportedName }}Storage and {{ .ExportedName }} interfaces
//
// mockery --dir=internals/{{ .ModName }}/{{ .ModName }} --name=^{{ .ExportedName }}.*

import (
	"testing"

	"{{ .PkgName }}template/internals/{{ .ModName }}/service"
	"{{ .PkgName }}template/mocks"
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

var INTERNALS_PING_STORAGE_PINGGO_MOD_TPL = Template{
	FilePath: "internals/ping/storage/ping.go",
	Render: true,
	Content: `package storage

import (
	"database/sql"
	"fmt"

	"{{ .PkgName }}template/internals/{{ .ModName }}/domain"
	"{{ .PkgName }}template/internals/{{ .ModName }}/serr"
	"{{ .PkgName }}template/server/infra"
)

type sqliteRepo struct {
	conn *sql.DB
}

func NewSqliteRepo(conn *sql.DB) domain.{{ .ExportedName }}Storage {
	return &sqliteRepo{
		conn: conn,
	}
}

func (s *sqliteRepo) {{ .ExportedName }}() error {
	logger := infra.GetLogger()
	rows, err := s.conn.Query(SELECT_ONE_TEST)
	if err != nil {
		return serr.Err{{ .ExportedName }}Failed
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logger.Error("connection close failed.")
		}
	}(rows)

	var result int

	if rows.Next() {
		err = rows.Scan(&result)
		if err != nil {
			return serr.Err{{ .ExportedName }}Failed
		}
	} else {
		return serr.Err{{ .ExportedName }}Failed
	}

	logger.Info(fmt.Sprintf("connection result: %d", result))
	return nil
}
`,
}

var INTERNALS_PING_STORAGE_PINGQUERIESGO_MOD_TPL = Template{
	FilePath: "internals/ping/storage/pingQueries.go",
	Render: true,
	Content: `package storage

const (
	SELECT_ONE_TEST = "SELECT 1"
)
`,
}

var INTERNALS_PING_WEB_PINGGO_MOD_TPL = Template{
	FilePath: "internals/ping/web/ping.go",
	Render: true,
	Content: `package web

import (
	"{{ .PkgName }}/gsk"
)

func HomeHandler(gc *gsk.Context) {

	gc.TemplateResponse(&gsk.Tpl{
		TemplatePath: "public/templates/index.html",
		Variables: gsk.Map{
			"Title":   "{{ .ExportedName }}",
			"Content": "Welcome to the {{ .ModName }} page!",
		},
	})

}
`,
}

var SERVER_ROUTING_PINGGO_MOD_TPL = Template{
	FilePath: "server/routing/ping.go",
	Render: true,
	Content: `package routing

import "{{ .PkgName }}template/internals/{{ .ModName }}"

func init() {
	RegisterApiRoutes({{ .ModName }}.SetupApiRoutes)
	RegisterWebRoutes({{ .ModName }}.SetupWebRoutes)
}
`,
}

var ModuleTemplates = []Template{
	INTERNALS_PING_API_HANDLER_PINGGO_MOD_TPL,
	INTERNALS_PING_API_HANDLER_PING_TESTGO_MOD_TPL,
	INTERNALS_PING_API_TRANSPORT_PINGGO_MOD_TPL,
	INTERNALS_PING_DOMAIN_PINGGO_MOD_TPL,
	INTERNALS_PING_SERR_PINGGO_MOD_TPL,
	INTERNALS_PING_SERVICE_PINGGO_MOD_TPL,
	INTERNALS_PING_SERVICE_PING_TESTGO_MOD_TPL,
	INTERNALS_PING_STORAGE_PINGGO_MOD_TPL,
	INTERNALS_PING_STORAGE_PINGQUERIESGO_MOD_TPL,
	INTERNALS_PING_WEB_PINGGO_MOD_TPL,
	SERVER_ROUTING_PINGGO_MOD_TPL,
}
