package middleware_test

import (
	"net/http"
	"testing"

	"github.com/adharshmk96/stk/gsk"
	"github.com/adharshmk96/stk/pkg/middleware"
	"github.com/stretchr/testify/assert"
)

func TestCORSDefault(t *testing.T) {
	// Create a new server instance
	config := &gsk.ServerConfig{
		Port: "8888",
	}
	s := gsk.New(config)

	s.Use(middleware.CORS())

	// Register a test route and handler
	s.Get("/", func(c gsk.Context) {
		c.Status(http.StatusOK).JSONResponse("OK")
	})

	t.Run("Non-preflight request", func(t *testing.T) {
		// Run the test request
		testParams := gsk.TestParams{
			Headers: map[string]string{
				"Origin": "example.com",
			},
		}
		rr, _ := s.Test("GET", "/", nil, testParams)

		expectedHeaders := map[string]string{
			"Access-Control-Allow-Origin":  "example.com",
			"Access-Control-Allow-Methods": "POST, GET, OPTIONS, PUT, DELETE, PATCH",
			"Access-Control-Allow-Headers": "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
		}

		// expect http.StatusOK
		if rr.Code != http.StatusOK {
			t.Errorf("Expected response code %d, but got %d", http.StatusOK, rr.Code)
		}

		for header, expectedValue := range expectedHeaders {
			if value := rr.Header().Get(header); value != expectedValue {
				t.Errorf("Expected %s header to be %q, but got %q", header, expectedValue, value)
			}
		}
	})

}

func TestCORSAllowedOrigin(t *testing.T) {
	// Create a new server instance
	config := &gsk.ServerConfig{
		Port: "8888",
	}

	AllowedOrigins := []string{
		"example.com",
	}
	s := gsk.New(config)

	s.Use(middleware.CORS(middleware.CORSConfig{
		AllowedOrigins: AllowedOrigins,
	}))

	// Register a test route and handler
	s.Get("/", func(c gsk.Context) {
		c.Status(http.StatusOK).JSONResponse("OK")
	})

	t.Run("Non-preflight request from example.com", func(t *testing.T) {

		// Run the test request
		testParams := gsk.TestParams{
			Headers: map[string]string{
				"Origin": "example.com",
			},
		}
		rr, _ := s.Test("GET", "/", nil, testParams)

		expectedHeaders := map[string]string{
			"Access-Control-Allow-Origin":  "example.com",
			"Access-Control-Allow-Methods": "POST, GET, OPTIONS, PUT, DELETE, PATCH",
			"Access-Control-Allow-Headers": "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
		}

		assert.Equal(t, http.StatusOK, rr.Code)

		for header, expectedValue := range expectedHeaders {
			value := rr.Header().Get(header)
			assert.Equal(t, expectedValue, value)
		}
	})

	t.Run("Non-preflight request with invalid origin", func(t *testing.T) {

		// Run the test request
		testParams := gsk.TestParams{
			Headers: map[string]string{
				"Origin": "invalid.com",
			},
		}
		rr, _ := s.Test("GET", "/", nil, testParams)

		expectedHeaders := map[string]string{
			"Access-Control-Allow-Origin":  "",
			"Access-Control-Allow-Methods": "",
			"Access-Control-Allow-Headers": "",
		}

		assert.Equal(t, http.StatusForbidden, rr.Code)

		for header, expectedValue := range expectedHeaders {
			value := rr.Header().Get(header)
			assert.Equal(t, expectedValue, value)
		}
	})

	t.Run("Preflight request with example.com", func(t *testing.T) {

		// Run the test request
		testParams := gsk.TestParams{
			Headers: map[string]string{
				"Origin":                        "example.com",
				"Access-Control-Request-Method": "POST",
			},
		}
		rr, _ := s.Test("OPTIONS", "/", nil, testParams)

		// NOTE: thie is behaviour from the router package
		// change this if we are chaning the router
		expectedHeaders := map[string]string{
			"Access-Control-Allow-Origin":  "example.com",
			"Access-Control-Allow-Methods": "POST, GET, OPTIONS, PUT, DELETE, PATCH",
			"Access-Control-Allow-Headers": "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
		}

		assert.Equal(t, http.StatusNoContent, rr.Code)

		for header, expectedValue := range expectedHeaders {
			value := rr.Header().Get(header)
			assert.Equal(t, expectedValue, value)
		}
	})

	t.Run("Preflight request with invalid origin", func(t *testing.T) {

		// Run the test request
		testParams := gsk.TestParams{
			Headers: map[string]string{
				"Origin":                        "invalid.com",
				"Access-Control-Request-Method": "POST",
			},
		}
		rr, _ := s.Test("OPTIONS", "/", nil, testParams)

		expectedHeaders := map[string]string{
			"Access-Control-Allow-Origin":  "",
			"Access-Control-Allow-Methods": "",
			"Access-Control-Allow-Headers": "",
		}

		// TODO - this should be checked later on
		assert.Equal(t, http.StatusForbidden, rr.Code)

		for header, expectedValue := range expectedHeaders {
			value := rr.Header().Get(header)
			assert.Equal(t, expectedValue, value)
		}
	})

}
