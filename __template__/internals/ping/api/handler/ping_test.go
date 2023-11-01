package handler_test

// run the following command to generate mocks for Ping interfaces
//
// mockery --dir=internals/ping/ping --name=^Ping.*

import (
	"net/http"
	"testing"

	"github.com/adharshmk96/stktemplate/internals/ping/api/handler"

	"github.com/adharshmk96/stk/gsk"
	"github.com/adharshmk96/stktemplate/mocks"
	"github.com/stretchr/testify/assert"
)

func TestPingHandler(t *testing.T) {
	t.Run("Ping Handler returns 200", func(t *testing.T) {

		// Arrange
		s := gsk.New()
		service := mocks.NewPingService(t)
		service.On("PingService").Return("pong", nil)

		pingHandler := handler.NewPingHandler(service)

		s.Get("/ping", pingHandler.PingHandler)

		// Act
		w, _ := s.Test("GET", "/ping", nil)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
