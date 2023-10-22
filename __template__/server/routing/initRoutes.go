package routing

import (
	"github.com/adharshmk96/stk/gsk"
)

func SetupRoutes(server *gsk.Server) {
	apiRoutes := server.RouteGroup("/api")

	setupPingRoutes(apiRoutes)
}
