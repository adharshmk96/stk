package routing

import (
	"github.com/adharshmk96/stk/gsk"
)

func SetupApiRoutes(server *gsk.Server) {
	apiRoutes := server.RouteGroup("/api")

	setupPingRoutes(apiRoutes)
}

func SetupTemplateRoutes(server *gsk.Server) {
	templateRoutes := server.RouteGroup("/")

	templateRoutes.Get("/", func(gc *gsk.Context) {
		gc.TemplateResponse(&gsk.Tpl{
			TemplatePath: "public/templates/index.html",
			Variables: gsk.Map{
				"Title":   "STK",
				"Content": "Hello, World!",
			},
		})
	})
}
