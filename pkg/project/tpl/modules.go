package tpl

var PINGGO_TPL_MOD = Template{
	FilePath: "{{ .ModName }}.go",
	Content:  ``,
}

var INTERNALSCOREENTITYPINGGO_TPL_MOD = Template{
	FilePath: "internals/core/entity/{{ .ModName }}.go",
	Content: `package entity

import "github.com/adharshmk96/stk/gsk"

// Domain
type User struct {
	pong string
}

// Storage
type {{ .ExportedName }}Storage interface {
	{{ .ExportedName }}() error
}

// Service
type {{ .ExportedName }}Service interface {
	{{ .ExportedName }}Service() string
}

// Handler
type {{ .ExportedName }}Handlers interface {
	{{ .ExportedName }}Handler(gc *gsk.Context)
}
`,
}

var INTERNALSCORESERRPINGGO_TPL_MOD = Template{
	FilePath: "internals/core/serr/{{ .ModName }}.go",
	Content: `package serr

import "errors"

var (
	Err{{ .ExportedName }}Failed = errors.New("{{ .ModName }} failed")
)
`,
}

var INTERNALSHTTPHANDLERPINGGO_TPL_MOD = Template{
	FilePath: "internals/http/handler/{{ .ModName }}.go",
	Content: `package handler

import (
	"net/http"

	"github.com/adharshmk96/stk/gsk"
)

/*
{{ .ExportedName }}Handler returns {{ .ModName }} 200 response
Response:
- 200: OK
- 500: Internal Server Error
*/
func (h *{{ .ModName }}Handler) {{ .ExportedName }}Handler(gc *gsk.Context) {

	{{ .ModName }} := h.service.{{ .ExportedName }}Service()

	gc.Status(http.StatusOK).JSONResponse(gsk.Map{
		"message": {{ .ModName }},
	})
}
`,
}

var INTERNALSHTTPHELPERSPINGGO_TPL_MOD = Template{
	FilePath: "internals/http/helpers/{{ .ModName }}.go",
	Content: `package helpers
`,
}

var INTERNALSHTTPTRANSPORTPINGGO_TPL_MOD = Template{
	FilePath: "internals/http/transport/{{ .ModName }}.go",
	Content: `package transport
`,
}

var INTERNALSSERVICEPINGGO_TPL_MOD = Template{
	FilePath: "internals/service/{{ .ModName }}.go",
	Content: `package service

func (s *{{ .ModName }}Service) {{ .ExportedName }}Service() string {
	err := s.storage.{{ .ExportedName }}()
	if err != nil {
		return "error"
	}
	return "pong"
}
`,
}

var INTERNALSSTORAGEPINGSTORAGEPINGGO_TPL_MOD = Template{
	FilePath: "internals/storage/{{ .ModName }}Storage/{{ .ModName }}.go",
	Content: `package {{ .ModName }}Storage

import "{{ .PkgName }}/internals/core/serr"

func (s *sqliteRepo) {{ .ExportedName }}() error {
	err := s.conn.{{ .ExportedName }}()
	if err != nil {
		return serr.Err{{ .ExportedName }}Failed
	}
	return nil
}
`,
}

var SERVERROUTINGPINGGO_TPL_MOD = Template{
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
	PINGGO_TPL_MOD,
	INTERNALSCOREENTITYPINGGO_TPL_MOD,
	INTERNALSCORESERRPINGGO_TPL_MOD,
	INTERNALSHTTPHANDLERPINGGO_TPL_MOD,
	INTERNALSHTTPHELPERSPINGGO_TPL_MOD,
	INTERNALSHTTPTRANSPORTPINGGO_TPL_MOD,
	INTERNALSSERVICEPINGGO_TPL_MOD,
	INTERNALSSTORAGEPINGSTORAGEPINGGO_TPL_MOD,
	SERVERROUTINGPINGGO_TPL_MOD,
}
