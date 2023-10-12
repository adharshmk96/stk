package handler

import (
	"net/http"

	"github.com/adharshmk96/stk-template/singlemod/internals/core/entity"
	"github.com/adharshmk96/stk/gsk"
)

type pingHandler struct {
	service entity.PingService
}

func NewPingHandler(service entity.PingService) entity.PingHandlers {
	return &pingHandler{
		service: service,
	}
}

/*
PingHandler returns ping 200 response
Response:
- 200: OK
- 500: Internal Server Error
*/
func (h *pingHandler) PingHandler(gc *gsk.Context) {

	ping, err := h.service.PingService()
	if err != nil {
		gc.Status(http.StatusInternalServerError).JSONResponse(gsk.Map{
			"error": err.Error(),
		})
		return
	}

	gc.Status(http.StatusOK).JSONResponse(gsk.Map{
		"message": ping,
	})
}
