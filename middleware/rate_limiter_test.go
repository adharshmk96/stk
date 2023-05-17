package middleware_test

import (
	"fmt"
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
	rateLimiter := middleware.NewRateLimiter(requestsPerInterval, interval)

	s.Use(rateLimiter.Middleware)

	s.Get("/test", dummyHandler)

	req, _ := http.NewRequest("GET", "/test", nil)
	respRec := httptest.NewRecorder()

	for i := 0; i < requestsPerInterval; i++ {
		s.Router.ServeHTTP(respRec, req)

		if respRec.Code != http.StatusOK {
			t.Errorf("Expected 200 OK, got: %d for request %d", respRec.Code, i+1)
		}
	}

	fmt.Println(rateLimiter.AccessCounter)
	s.Router.ServeHTTP(respRec, req)

	if respRec.Code != http.StatusTooManyRequests {
		t.Errorf("Expected 429 Too Many Requests, got: %d", respRec.Code)
	}
}
