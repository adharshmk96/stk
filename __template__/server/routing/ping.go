package routing

import "github.com/adharshmk96/stktemplate/internals/ping"

func init() {
	RegisterApiRoutes(ping.SetupApiRoutes)
	RegisterWebRoutes(ping.SetupWebRoutes)
}
