package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adharshmk96/stk"
	"github.com/adharshmk96/stk/middleware"
	"github.com/stretchr/testify/assert"
)

func TestSecurityHeaders(t *testing.T) {
	// Create a new server instance
	config := &stk.ServerConfig{
		Port:           "8080",
		RequestLogging: false,
	}
	s := stk.NewServer(config)

	s.Use(middleware.SecurityHeaders)

	// Register a test route and handler
	s.Get("/", func(c stk.Context) {
		c.Status(http.StatusOK).JSONResponse("OK")
	})

	// Run the test request
	req, _ := http.NewRequest("GET", "/", nil)
	respRec := httptest.NewRecorder()

	s.Router.ServeHTTP(respRec, req)

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
		value := respRec.Header().Get(header)
		assert.Equal(t, expectedValue, value)
	}
}
