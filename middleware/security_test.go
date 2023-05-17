package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adharshmk96/stk"
)

func TestSecurityHeaders(t *testing.T) {
	// Create a new server instance
	config := &stk.ServerConfig{
		Port:           "8080",
		RequestLogging: true,
	}
	s := stk.NewServer(config)

	// Register a test route and handler
	s.Get("/", func(c *stk.Context) {
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
		if value := respRec.Header().Get(header); value != expectedValue {
			t.Errorf("Expected %s header to be %q, but got %q", header, expectedValue, value)
		}
	}
}
