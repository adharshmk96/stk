package tpl

var INTERNALS_PING_ROUTESGO_MOD_TPL = Template{
	FilePath: "internals/ping/routes.go",
	Render: true,
	Content: `package {{ .ModName }}

import (
	"github.com/adharshmk96/stk/gsk"
	"{{ .PkgName }}/internals/{{ .ModName }}/api/handler"
	"{{ .PkgName }}/internals/{{ .ModName }}/domain"
	"{{ .PkgName }}/internals/{{ .ModName }}/service"
	"{{ .PkgName }}/internals/{{ .ModName }}/storage"
	"{{ .PkgName }}/internals/{{ .ModName }}/web"
	"{{ .PkgName }}/server/infra/db"
)

func initialize{{ .ExportedName }}Handler() domain.{{ .ExportedName }}Handlers {
	conn := db.GetSqliteConnection()

	{{ .ModName }}Storage := storage.NewSqliteRepo(conn)
	{{ .ModName }}Service := service.New{{ .ExportedName }}Service({{ .ModName }}Storage)
	{{ .ModName }}Handler := handler.New{{ .ExportedName }}Handler({{ .ModName }}Service)

	return {{ .ModName }}Handler
}

func SetupApiRoutes(rg *gsk.RouteGroup) {
	{{ .ModName }}Handler := initialize{{ .ExportedName }}Handler()

	{{ .ModName }}Routes := rg.RouteGroup("/{{ .ModName }}")

	{{ .ModName }}Routes.Get("/", {{ .ModName }}Handler.{{ .ExportedName }}Handler)
}

func SetupWebRoutes(rg *gsk.RouteGroup) {
	rg.Get("/{{ .ModName }}", web.HomeHandler)
}
`,
}

var INTERNALS_PING_API_HANDLER_PINGGO_MOD_TPL = Template{
	FilePath: "internals/ping/api/handler/ping.go",
	Render: true,
	Content: `package handler

import (
	"net/http"

	"{{ .PkgName }}/internals/{{ .ModName }}/domain"

	"github.com/adharshmk96/stk/gsk"
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

	"{{ .PkgName }}/internals/{{ .ModName }}/api/handler"

	"github.com/adharshmk96/stk/gsk"
	"{{ .PkgName }}/mocks"
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

var INTERNALS_PING_DOMAIN_HANDLERGO_MOD_TPL = Template{
	FilePath: "internals/ping/domain/handler.go",
	Render: true,
	Content: `package domain

import "github.com/adharshmk96/stk/gsk"

// Handler
type {{ .ExportedName }}Handlers interface {
	{{ .ExportedName }}Handler(gc *gsk.Context)
}
`,
}

var INTERNALS_PING_DOMAIN_PINGGO_MOD_TPL = Template{
	FilePath: "internals/ping/domain/ping.go",
	Render: true,
	Content: `package domain

// Domain
type {{ .ExportedName }}Data struct {
	{{ .ModName }} string
}
`,
}

var INTERNALS_PING_DOMAIN_SERVICEGO_MOD_TPL = Template{
	FilePath: "internals/ping/domain/service.go",
	Render: true,
	Content: `package domain

// Service
type {{ .ExportedName }}Service interface {
	{{ .ExportedName }}Service() (string, error)
}
`,
}

var INTERNALS_PING_DOMAIN_STORAGEGO_MOD_TPL = Template{
	FilePath: "internals/ping/domain/storage.go",
	Render: true,
	Content: `package domain

// Storage
type {{ .ExportedName }}Storage interface {
	{{ .ExportedName }}() error
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
	"{{ .PkgName }}/internals/{{ .ModName }}/domain"
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

	"{{ .PkgName }}/internals/{{ .ModName }}/service"
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

var INTERNALS_PING_STORAGE_PINGGO_MOD_TPL = Template{
	FilePath: "internals/ping/storage/ping.go",
	Render: true,
	Content: `package storage

import (
	"database/sql"
	"fmt"

	"{{ .PkgName }}/internals/{{ .ModName }}/domain"
	"{{ .PkgName }}/internals/{{ .ModName }}/serr"
	"{{ .PkgName }}/server/infra"
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
	"github.com/adharshmk96/stk/gsk"
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

var MOCKS_PINGHANDLERSGO_MOD_TPL = Template{
	FilePath: "mocks/PingHandlers.go",
	Render: true,
	Content: `// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	gsk "github.com/adharshmk96/stk/gsk"
	mock "github.com/stretchr/testify/mock"
)

// {{ .ExportedName }}Handlers is an autogenerated mock type for the {{ .ExportedName }}Handlers type
type {{ .ExportedName }}Handlers struct {
	mock.Mock
}

// {{ .ExportedName }}Handler provides a mock function with given fields: gc
func (_m *{{ .ExportedName }}Handlers) {{ .ExportedName }}Handler(gc *gsk.Context) {
	_m.Called(gc)
}

type mockConstructorTestingTNew{{ .ExportedName }}Handlers interface {
	mock.TestingT
	Cleanup(func())
}

// New{{ .ExportedName }}Handlers creates a new instance of {{ .ExportedName }}Handlers. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func New{{ .ExportedName }}Handlers(t mockConstructorTestingTNew{{ .ExportedName }}Handlers) *{{ .ExportedName }}Handlers {
	mock := &{{ .ExportedName }}Handlers{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
`,
}

var MOCKS_PINGSERVICEGO_MOD_TPL = Template{
	FilePath: "mocks/PingService.go",
	Render: true,
	Content: `// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// {{ .ExportedName }}Service is an autogenerated mock type for the {{ .ExportedName }}Service type
type {{ .ExportedName }}Service struct {
	mock.Mock
}

// {{ .ExportedName }}Service provides a mock function with given fields:
func (_m *{{ .ExportedName }}Service) {{ .ExportedName }}Service() (string, error) {
	ret := _m.Called()

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNew{{ .ExportedName }}Service interface {
	mock.TestingT
	Cleanup(func())
}

// New{{ .ExportedName }}Service creates a new instance of {{ .ExportedName }}Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func New{{ .ExportedName }}Service(t mockConstructorTestingTNew{{ .ExportedName }}Service) *{{ .ExportedName }}Service {
	mock := &{{ .ExportedName }}Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
`,
}

var MOCKS_PINGSTORAGEGO_MOD_TPL = Template{
	FilePath: "mocks/PingStorage.go",
	Render: true,
	Content: `// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// {{ .ExportedName }}Storage is an autogenerated mock type for the {{ .ExportedName }}Storage type
type {{ .ExportedName }}Storage struct {
	mock.Mock
}

// {{ .ExportedName }} provides a mock function with given fields:
func (_m *{{ .ExportedName }}Storage) {{ .ExportedName }}() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNew{{ .ExportedName }}Storage interface {
	mock.TestingT
	Cleanup(func())
}

// New{{ .ExportedName }}Storage creates a new instance of {{ .ExportedName }}Storage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func New{{ .ExportedName }}Storage(t mockConstructorTestingTNew{{ .ExportedName }}Storage) *{{ .ExportedName }}Storage {
	mock := &{{ .ExportedName }}Storage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
`,
}

var SERVER_ROUTING_PINGGO_MOD_TPL = Template{
	FilePath: "server/routing/ping.go",
	Render: true,
	Content: `package routing

import "{{ .PkgName }}/internals/{{ .ModName }}"

func init() {
	RegisterApiRoutes({{ .ModName }}.SetupApiRoutes)
	RegisterWebRoutes({{ .ModName }}.SetupWebRoutes)
}
`,
}

var ModuleTemplates = []Template{
	INTERNALS_PING_ROUTESGO_MOD_TPL,
	INTERNALS_PING_API_HANDLER_PINGGO_MOD_TPL,
	INTERNALS_PING_API_HANDLER_PING_TESTGO_MOD_TPL,
	INTERNALS_PING_API_TRANSPORT_PINGGO_MOD_TPL,
	INTERNALS_PING_DOMAIN_HANDLERGO_MOD_TPL,
	INTERNALS_PING_DOMAIN_PINGGO_MOD_TPL,
	INTERNALS_PING_DOMAIN_SERVICEGO_MOD_TPL,
	INTERNALS_PING_DOMAIN_STORAGEGO_MOD_TPL,
	INTERNALS_PING_SERR_PINGGO_MOD_TPL,
	INTERNALS_PING_SERVICE_PINGGO_MOD_TPL,
	INTERNALS_PING_SERVICE_PING_TESTGO_MOD_TPL,
	INTERNALS_PING_STORAGE_PINGGO_MOD_TPL,
	INTERNALS_PING_STORAGE_PINGQUERIESGO_MOD_TPL,
	INTERNALS_PING_WEB_PINGGO_MOD_TPL,
	MOCKS_PINGHANDLERSGO_MOD_TPL,
	MOCKS_PINGSERVICEGO_MOD_TPL,
	MOCKS_PINGSTORAGEGO_MOD_TPL,
	SERVER_ROUTING_PINGGO_MOD_TPL,
}
