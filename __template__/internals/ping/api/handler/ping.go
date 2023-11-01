package handler

import (
	"net/http"

	"github.com/adharshmk96/stktemplate/internals/ping/domain"

	"github.com/adharshmk96/stk/gsk"
)

type pingHandler struct {
	service domain.PingService
}

func NewPingHandler(service domain.PingService) domain.PingHandlers {
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

	message, err := h.service.PingService()
	if err != nil {
		gc.Status(http.StatusInternalServerError).JSONResponse(gsk.Map{
			"error": err.Error(),
		})
		return
	}

	gc.Status(http.StatusOK).JSONResponse(gsk.Map{
		"message": message,
	})
}
