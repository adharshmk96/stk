package routing

import (
	"github.com/adharshmk96/stk-template/singlemod/internals/core/entity"
	"github.com/adharshmk96/stk-template/singlemod/internals/http/handler"
	"github.com/adharshmk96/stk-template/singlemod/internals/service"
	"github.com/adharshmk96/stk-template/singlemod/internals/storage/pingStorage"
	"github.com/adharshmk96/stk-template/singlemod/server/infra/db"
	"github.com/adharshmk96/stk/gsk"
)

func initializePing() entity.PingHandlers {
	conn := db.GetSqliteConnection()

	pingStorage := pingStorage.NewSqliteRepo(conn)
	pingService := service.NewPingService(pingStorage)
	pingHandler := handler.NewPingHandler(pingService)

	return pingHandler
}

func setupPingRoutes(rg *gsk.RouteGroup) {
	pingHandler := initializePing()

	pingRoutes := rg.RouteGroup("/ping")

	pingRoutes.Get("/", pingHandler.PingHandler)
}
