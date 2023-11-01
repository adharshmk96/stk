package ping

import (
	"github.com/adharshmk96/stk/gsk"
	"github.com/adharshmk96/stktemplate/internals/ping/api/handler"
	"github.com/adharshmk96/stktemplate/internals/ping/domain"
	"github.com/adharshmk96/stktemplate/internals/ping/service"
	"github.com/adharshmk96/stktemplate/internals/ping/storage"
	"github.com/adharshmk96/stktemplate/internals/ping/web"
	"github.com/adharshmk96/stktemplate/server/infra/db"
)

func initializePingHandler() domain.PingHandlers {
	conn := db.GetSqliteConnection()

	pingStorage := storage.NewSqliteRepo(conn)
	pingService := service.NewPingService(pingStorage)
	pingHandler := handler.NewPingHandler(pingService)

	return pingHandler
}

func SetupApiRoutes(rg *gsk.RouteGroup) {
	pingHandler := initializePingHandler()

	pingRoutes := rg.RouteGroup("/ping")

	pingRoutes.Get("/", pingHandler.PingHandler)
}

func SetupWebRoutes(rg *gsk.RouteGroup) {
	rg.Get("/ping", web.HomeHandler)
}
