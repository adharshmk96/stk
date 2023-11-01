package service_test

// run the following command to generate mocks for PingStorage and Ping interfaces
//
// mockery --dir=internals/ping/ping --name=^Ping.*

import (
	"testing"

	"github.com/adharshmk96/stktemplate/internals/ping/service"
	"github.com/adharshmk96/stktemplate/mocks"
	"github.com/stretchr/testify/assert"
)

func TestPingService(t *testing.T) {
	t.Run("PingService returns pong", func(t *testing.T) {

		// Arrange
		storage := mocks.NewPingStorage(t)
		storage.On("Ping").Return(nil)

		svc := service.NewPingService(storage)

		// Act
		msg, err := svc.PingService()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "pong", msg)
	})
}
