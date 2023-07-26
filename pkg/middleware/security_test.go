package middleware_test

import (
	"net/http"
	"testing"

	"github.com/adharshmk96/stk/gsk"
	"github.com/adharshmk96/stk/pkg/middleware"
	"github.com/stretchr/testify/assert"
)

func TestSecurityHeaders(t *testing.T) {
	// Create a new server instance
	config := &gsk.ServerConfig{
		Port: "8888",
	}
	s := gsk.New(config)

	s.Use(middleware.SecurityHeaders)

	// Register a test route and handler
	s.Get("/", func(c *gsk.Context) {
		c.Status(http.StatusOK).JSONResponse("OK")
	})

	// Run the test request
	rr, _ := s.Test("GET", "/", nil)

	expectedHeaders := map[string]string{
		"X-Content-Type-Options":            "nosniff",
		"X-Frame-Options":                   "SAMEORIGIN",
		"X-XSS-Protection":                  "1; mode=block",
		"Content-Security-Policy":           "default-src 'self';",
		"X-Permitted-Cross-Domain-Policies": "master-only",
		"Strict-Transport-Security":         "max-age=31536000; includeSubDomains",
		"Referrer-Policy":                   "strict-origin-when-cross-origin",
	}

	for header, expectedValue := range expectedHeaders {
		value := rr.Header().Get(header)
		assert.Equal(t, expectedValue, value)
	}
}
