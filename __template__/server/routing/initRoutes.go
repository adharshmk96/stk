package routing

import (
	"github.com/adharshmk96/stk/gsk"
)

var routeGroups = []func(*gsk.RouteGroup){}

func SetupApiRoutes(server *gsk.Server) {
	apiRoutes := server.RouteGroup("/api")

	for _, routeGroup := range routeGroups {
		routeGroup(apiRoutes)
	}
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
