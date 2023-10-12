package routing

import (
	"github.com/adharshmk96/stk-template/singlemod/internals/http/handler"
	"github.com/adharshmk96/stk-template/singlemod/internals/service"
	"github.com/adharshmk96/stk-template/singlemod/internals/storage/pingStorage"
	"github.com/adharshmk96/stk-template/singlemod/server/infra"
	"github.com/adharshmk96/stk/gsk"
	"github.com/adharshmk96/stk/pkg/db"
	"github.com/spf13/viper"
)

func setupPingRoutes(server *gsk.Server) {
	dbConfig := viper.GetString(infra.ENV_SQLITE_FILEPATH)
	conn := db.GetSqliteConnection(dbConfig)

	pingStorage := pingStorage.NewSqliteRepo(conn)
	pingService := service.NewPingService(pingStorage)
	pingHandler := handler.NewPingHandler(pingService)

	server.Get("/ping", pingHandler.PingHandler)
}
