package stk_test

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

func TestCORSDefault(t *testing.T) {
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

	t.Run("Non-preflight request", func(t *testing.T) {
		// Run the test request
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Origin", "https://example.com")
		respRec := httptest.NewRecorder()

		s.Router.ServeHTTP(respRec, req)

		expectedHeaders := map[string]string{
			"Access-Control-Allow-Origin":  "https://example.com",
			"Access-Control-Allow-Methods": "POST, GET, OPTIONS, PUT, DELETE, PATCH",
			"Access-Control-Allow-Headers": "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
		}

		for header, expectedValue := range expectedHeaders {
			if value := respRec.Header().Get(header); value != expectedValue {
				t.Errorf("Expected %s header to be %q, but got %q", header, expectedValue, value)
			}
		}
	})

}

func TestCORSAllowedOrigin(t *testing.T) {
	// Create a new server instance
	config := &stk.ServerConfig{
		Port:           "8080",
		RequestLogging: true,
		AllowedOrigins: []string{
			"https://example.com",
		},
	}
	s := stk.NewServer(config)

	// Register a test route and handler
	s.Get("/", func(c *stk.Context) {
		c.Status(http.StatusOK).JSONResponse("OK")
	})

	t.Run("Non-preflight request", func(t *testing.T) {
		// Run the test request
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Origin", "https://example.com")
		respRec := httptest.NewRecorder()

		s.Router.ServeHTTP(respRec, req)

		expectedHeaders := map[string]string{
			"Access-Control-Allow-Origin":  "https://example.com",
			"Access-Control-Allow-Methods": "POST, GET, OPTIONS, PUT, DELETE, PATCH",
			"Access-Control-Allow-Headers": "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
		}

		for header, expectedValue := range expectedHeaders {
			if value := respRec.Header().Get(header); value != expectedValue {
				t.Errorf("Expected %s header to be %q, but got %q", header, expectedValue, value)
			}
		}
	})

	t.Run("Non-preflight request with invalid origin", func(t *testing.T) {
		// Run the test request
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Origin", "https://invalid.com")
		respRec := httptest.NewRecorder()

		s.Router.ServeHTTP(respRec, req)

		expectedHeaders := map[string]string{
			"Access-Control-Allow-Origin":  "",
			"Access-Control-Allow-Methods": "POST, GET, OPTIONS, PUT, DELETE, PATCH",
			"Access-Control-Allow-Headers": "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
		}

		for header, expectedValue := range expectedHeaders {
			if value := respRec.Header().Get(header); value != expectedValue {
				t.Errorf("Expected %s header to be %q, but got %q", header, expectedValue, value)
			}
		}
	})

	t.Run("Preflight request", func(t *testing.T) {
		// Run the test request
		req, _ := http.NewRequest("OPTIONS", "/", nil)
		req.Header.Set("Origin", "https://example.com")
		req.Header.Set("Access-Control-Request-Method", "POST")
		respRec := httptest.NewRecorder()

		s.Router.ServeHTTP(respRec, req)

		expectedHeaders := map[string]string{
			"Access-Control-Allow-Origin":  "https://example.com",
			"Access-Control-Allow-Methods": "POST, GET, OPTIONS, PUT, DELETE, PATCH",
			"Access-Control-Allow-Headers": "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
		}

		for header, expectedValue := range expectedHeaders {
			if value := respRec.Header().Get(header); value != expectedValue {
				t.Errorf("Expected %s header to be %q, but got %q", header, expectedValue, value)
			}
		}
	})

	t.Run("Preflight request with invalid origin", func(t *testing.T) {
		// Run the test request
		req, _ := http.NewRequest("OPTIONS", "/", nil)
		req.Header.Set("Origin", "https://invalid.com")
		req.Header.Set("Access-Control-Request-Method", "POST")
		respRec := httptest.NewRecorder()

		s.Router.ServeHTTP(respRec, req)

		expectedHeaders := map[string]string{
			"Access-Control-Allow-Origin":  "",
			"Access-Control-Allow-Methods": "POST, GET, OPTIONS, PUT, DELETE, PATCH",
			"Access-Control-Allow-Headers": "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
		}

		for header, expectedValue := range expectedHeaders {
			if value := respRec.Header().Get(header); value != expectedValue {
				t.Errorf("Expected %s header to be %q, but got %q", header, expectedValue, value)
			}
		}
	})

}
