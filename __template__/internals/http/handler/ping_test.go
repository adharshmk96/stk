package handler_test

// run the following command to generate mocks for Ping interfaces
//
// mockery --dir=internals/core/entity --name=^Ping.*
//
// and uncomment the following code

/*

import (
	"net/http"
	"testing"

	"github.com/adharshmk96/stk-template/singlemod/internals/http/handler"
	"github.com/adharshmk96/stk-template/singlemod/mocks"
	"github.com/adharshmk96/stk/gsk"
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

*/
