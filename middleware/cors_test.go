package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adharshmk96/stk"
	"github.com/adharshmk96/stk/middleware"
)

func TestCORSDefault(t *testing.T) {
	// Create a new server instance
	config := &stk.ServerConfig{
		Port:           "8080",
		RequestLogging: true,
	}
	s := stk.NewServer(config)

	s.Use(middleware.CORS)

	// Register a test route and handler
	s.Get("/", func(c *stk.Context) {
		c.Status(http.StatusOK).JSONResponse("OK")
	})

	t.Run("Non-preflight request", func(t *testing.T) {
		// Run the test request
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Host", "example.com")
		respRec := httptest.NewRecorder()

		s.Router.ServeHTTP(respRec, req)

		expectedHeaders := map[string]string{
			"Access-Control-Allow-Origin":  "example.com",
			"Access-Control-Allow-Methods": "POST, GET, OPTIONS, PUT, DELETE, PATCH",
			"Access-Control-Allow-Headers": "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
		}

		// expect http.StatusOK
		if respRec.Code != http.StatusOK {
			t.Errorf("Expected response code %d, but got %d", http.StatusOK, respRec.Code)
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
			"example.com",
		},
	}
	s := stk.NewServer(config)

	s.Use(middleware.CORS)

	// Register a test route and handler
	s.Get("/", func(c *stk.Context) {
		c.Status(http.StatusOK).JSONResponse("OK")
	})

	t.Run("Non-preflight request from example.com", func(t *testing.T) {
		// Run the test request
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Host", "example.com")
		respRec := httptest.NewRecorder()

		s.Router.ServeHTTP(respRec, req)

		expectedHeaders := map[string]string{
			"Access-Control-Allow-Origin":  "example.com",
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
		req.Header.Set("Host", "invalid.com")
		respRec := httptest.NewRecorder()

		s.Router.ServeHTTP(respRec, req)

		expectedHeaders := map[string]string{
			"Access-Control-Allow-Origin":  "",
			"Access-Control-Allow-Methods": "",
			"Access-Control-Allow-Headers": "",
		}

		for header, expectedValue := range expectedHeaders {
			if value := respRec.Header().Get(header); value != expectedValue {
				t.Errorf("Expected %s header to be %q, but got %q", header, expectedValue, value)
			}
		}
	})

	t.Run("Preflight request with example.com", func(t *testing.T) {
		// Run the test request
		req, _ := http.NewRequest("OPTIONS", "/", nil)
		req.Header.Set("Host", "example.com")
		req.Header.Set("Access-Control-Request-Method", "POST")
		respRec := httptest.NewRecorder()

		s.Router.ServeHTTP(respRec, req)

		// NOTE: thie is behaviour from the router package
		// change this if we are chaning the router
		expectedHeaders := map[string]string{
			"Access-Control-Allow-Origin":  "",
			"Access-Control-Allow-Methods": "",
			"Access-Control-Allow-Headers": "",
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
		req.Header.Set("Host", "invalid.com")
		req.Header.Set("Access-Control-Request-Method", "POST")
		respRec := httptest.NewRecorder()

		s.Router.ServeHTTP(respRec, req)

		expectedHeaders := map[string]string{
			"Access-Control-Allow-Origin":  "",
			"Access-Control-Allow-Methods": "",
			"Access-Control-Allow-Headers": "",
		}

		for header, expectedValue := range expectedHeaders {
			if value := respRec.Header().Get(header); value != expectedValue {
				t.Errorf("Expected %s header to be %q, but got %q", header, expectedValue, value)
			}
		}
	})

}
