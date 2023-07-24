package tpl

func HandlerHandlersTemplate() []byte {
	return []byte(`package {{.DirNames.HandlerDir}}

import (
	"{{.PkgName}}/{{.DirTree.EntityPath}}"
)

type pingHandler struct {
	pingService {{.DirNames.EntityDir}}.PingService
}

func NewPingHandler(pingService {{.DirNames.EntityDir}}.PingService) {{.DirNames.EntityDir}}.PingHandler {
	return &pingHandler{
		pingService: pingService,
	}
}
`)
}

func HandlerPingTemplate() []byte {
	return []byte(`package {{.DirNames.HandlerDir}}

import (
	"net/http"

	"github.com/adharshmk96/stk/gsk"
)

/*
PingHandler returns ping 200 response
Response:
- 200: OK
- 500: Internal Server Error
*/
func (h *pingHandler) PingHandler(gc *gsk.Context) {
	
	ping := h.pingService.PingService()

	gc.Status(http.StatusOK).JSONResponse(gsk.Map{
		"message": ping,
	})
}	
`)
}
