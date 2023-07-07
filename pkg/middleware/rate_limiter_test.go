package middleware_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/adharshmk96/stk/gsk"
	"github.com/adharshmk96/stk/pkg/middleware"
	"github.com/stretchr/testify/assert"
)

func dummyHandler(c gsk.Context) {
	c.Status(http.StatusOK).JSONResponse("OK")
}

func TestRateLimiter(t *testing.T) {
	// Create a new server instance
	config := &gsk.ServerConfig{
		Port: "8888",
	}
	s := gsk.New(config)

	// rate limiter middleware
	requestsPerInterval := 5
	interval := 1 * time.Second

	rlConfig := middleware.RateLimiterConfig{
		RequestsPerInterval: requestsPerInterval,
		Interval:            interval,
	}

	rateLimiter := middleware.NewRateLimiter(rlConfig)

	s.Use(rateLimiter.Middleware)

	s.Get("/test", dummyHandler)

	for i := 0; i < requestsPerInterval; i++ {
		rr, _ := s.Test("GET", "/test", nil)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200 OK, got: %d for request %d", rr.Code, i+1)
		}
	}

	rr, _ := s.Test("GET", "/test", nil)

	assert.Equal(t, http.StatusTooManyRequests, rr.Code)
}
