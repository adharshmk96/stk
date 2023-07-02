package tpl

func ServerSetupTemplate() []byte {
	return []byte(`package {{.DirNames.ServerDir}}

import "github.com/adharshmk96/stk/gsk"

func StartServer(port string) *gsk.Server {

	serverConfig := &gsk.ServerConfig{
		Port:           port,
		RequestLogging: true,
	}

	server := gsk.NewServer(serverConfig)

	setupRoutes(server)

	rateLimiter := rateLimiter()
	server.Use(rateLimiter)

	server.Start()

	return server
}
`)
}

func ServerMiddlewareTemplate() []byte {
	return []byte(`package {{.DirNames.ServerDir}}

import (
	"time"

	"github.com/adharshmk96/stk/gsk"
	"github.com/adharshmk96/stk/pkg/middleware"
)

func rateLimiter() gsk.Middleware {
	rateLimiter := middleware.NewRateLimiter(60, 10*time.Second)
	return rateLimiter.Middleware
}
`)
}

// Server router with handler, service and storage
func ServerRouterTemplate() []byte {
	return []byte(`package {{.DirNames.ServerDir}}

import (
	"github.com/adharshmk96/stk/gsk"
	"{{.PkgName}}/{{.DirTree.HandlerPath}}"
	"{{.PkgName}}/{{.DirTree.ServicePath}}"
	"{{.PkgName}}/{{.DirTree.StoragePath}}/{{.DirNames.SqliteRepoDir}}"
	"github.com/adharshmk96/stk/pkg/db"
)

func setupRoutes(server *gsk.Server) {

	conn := db.GetSqliteConnection("sqlite.db")

	{{.AppName}}Storage := {{.DirNames.SqliteRepoDir}}.NewSqliteRepo(conn)
	{{.AppName}}Service := {{.DirNames.ServiceDir}}.NewPingService({{.AppName}}Storage)
	{{.AppName}}Handler := {{.DirNames.HandlerDir}}.NewPingHandler({{.AppName}}Service)

	server.Get("/ping", {{.AppName}}Handler.PingHandler)
}
`)
}

// Server router without storage
func ServerRouterNoStorageTemplate() []byte {
	return []byte(`package {{.DirNames.ServerDir}}

import (
	"github.com/adharshmk96/stk/gsk"
	"{{.PkgName}}/{{.DirTree.HandlerPath}}"
	"{{.PkgName}}/pkg/services"
)

func setupRoutes(server *gsk.Server) {
	{{.AppName}}Service := services.NewPingService({{.AppName}}Storage)
	{{.AppName}}Handler := handlers.NewPingHandler({{.AppName}}Service)

	server.GET("/ping", handlers.PingHandler)
}
`)
}
