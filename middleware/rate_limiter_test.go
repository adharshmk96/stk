package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/adharshmk96/stk"
	"github.com/adharshmk96/stk/middleware"
)

func dummyHandler(c *stk.Context) {
	c.Status(http.StatusOK).JSONResponse("OK")
}

func TestRateLimiter(t *testing.T) {
	// Create a new server instance
	config := &stk.ServerConfig{
		Port:           "8080",
		RequestLogging: true,
	}
	s := stk.NewServer(config)

	// rate limiter middleware
	requestsPerInterval := 5
	interval := 1 * time.Second
	rlMiddleware := middleware.NewRateLimiter(requestsPerInterval, interval)

	s.Use(rlMiddleware)

	s.Get("/test", dummyHandler)

	respRec := httptest.NewRecorder()

	for i := 0; i < requestsPerInterval+1; i++ {
		req, _ := http.NewRequest("GET", "/test", nil)
		s.Router.ServeHTTP(respRec, req)
	}

	if respRec.Code != http.StatusTooManyRequests {
		t.Errorf("Expected 429 Too Many Requests, got: %d", respRec.Code)
	}
}
