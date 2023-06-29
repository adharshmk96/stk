package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/adharshmk96/stk/gsk"
	"github.com/adharshmk96/stk/middleware"
	"github.com/stretchr/testify/assert"
)

func dummyHandler(c gsk.Context) {
	c.Status(http.StatusOK).JSONResponse("OK")
}

func TestRateLimiter(t *testing.T) {
	// Create a new server instance
	config := &gsk.ServerConfig{
		Port:           "8080",
		RequestLogging: false,
	}
	s := gsk.NewServer(config)

	// rate limiter middleware
	requestsPerInterval := 5
	interval := 1 * time.Second
	rateLimiter := middleware.NewRateLimiter(requestsPerInterval, interval)

	s.Use(rateLimiter.Middleware)

	s.Get("/test", dummyHandler)

	for i := 0; i < requestsPerInterval; i++ {
		req, _ := http.NewRequest("GET", "/test", nil)
		respRec := httptest.NewRecorder()
		s.Router.ServeHTTP(respRec, req)

		if respRec.Code != http.StatusOK {
			t.Errorf("Expected 200 OK, got: %d for request %d", respRec.Code, i+1)
		}
	}

	req, _ := http.NewRequest("GET", "/test", nil)
	respRec := httptest.NewRecorder()
	s.Router.ServeHTTP(respRec, req)

	assert.Equal(t, http.StatusTooManyRequests, respRec.Code)
}
