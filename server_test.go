package stk_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/adharshmk96/stk"
)

// Test server routes

func TestServerRoutes(t *testing.T) {
	config := &stk.ServerConfig{
		Port:           "8080",
		RequestLogging: true,
	}
	s := stk.NewServer(config)

	test_status := http.StatusNoContent

	sampleHandler := func(ctx *stk.Context) {
		method := ctx.Request.Method
		ctx.Writer.WriteHeader(test_status)
		ctx.Writer.Write([]byte(method))
	}

	s.Get("/test-get", sampleHandler)
	s.Post("/test-post", sampleHandler)
	s.Put("/test-put", sampleHandler)
	s.Delete("/test-delete", sampleHandler)
	s.Patch("/test-patch", sampleHandler)

	queryParamHandler := func(ctx *stk.Context) {
		ctx.Writer.WriteHeader(test_status)
		ctx.Writer.Write([]byte(ctx.GetQueryParam("name")))
	}

	s.Get("/test/p", queryParamHandler)
	s.Post("/test/p", queryParamHandler)
	s.Put("/test/p", queryParamHandler)
	s.Delete("/test/p", queryParamHandler)
	s.Patch("/test/p", queryParamHandler)

	paramsHandler := func(ctx *stk.Context) {
		ctx.Writer.WriteHeader(test_status)
		ctx.Writer.Write([]byte(ctx.GetParam("id")))
	}

	s.Get("/test/d/:id", paramsHandler)
	s.Post("/test/d/:id", paramsHandler)
	s.Put("/test/d/:id", paramsHandler)
	s.Delete("/test/d/:id", paramsHandler)
	s.Patch("/test/d/:id", paramsHandler)

	serverHandler := http.HandlerFunc(s.Router.ServeHTTP)

	testCases := []struct {
		name       string
		method     string
		path       string
		statusCode int
		expected   string
	}{
		{name: "testing get for 200", method: http.MethodGet, path: "/test-get", statusCode: test_status, expected: "GET"},
		{name: "testing post for 200", method: http.MethodPost, path: "/test-post", statusCode: test_status, expected: "POST"},
		{name: "testing put for 200", method: http.MethodPut, path: "/test-put", statusCode: test_status, expected: "PUT"},
		{name: "testing delete for 200", method: http.MethodDelete, path: "/test-delete", statusCode: test_status, expected: "DELETE"},
		{name: "testing patch for 200", method: http.MethodPatch, path: "/test-patch", statusCode: test_status, expected: "PATCH"},

		{name: "testing get with dynamic route", method: http.MethodGet, path: "/test/d/123", statusCode: test_status, expected: "123"},
		{name: "testing post with dynamic route", method: http.MethodPost, path: "/test/d/123", statusCode: test_status, expected: "123"},
		{name: "testing put with dynamic route", method: http.MethodPut, path: "/test/d/123", statusCode: test_status, expected: "123"},
		{name: "testing delete with dynamic route", method: http.MethodDelete, path: "/test/d/123", statusCode: test_status, expected: "123"},
		{name: "testing patch with dynamic route", method: http.MethodPatch, path: "/test/d/123", statusCode: test_status, expected: "123"},

		{name: "testing get with param name=adha", method: http.MethodGet, path: "/test/p?name=adha", statusCode: test_status, expected: "adha"},
		{name: "testing post with param name=adha", method: http.MethodPost, path: "/test/p?name=adha", statusCode: test_status, expected: "adha"},
		{name: "testing put with param name=adha", method: http.MethodPut, path: "/test/p?name=adha", statusCode: test_status, expected: "adha"},
		{name: "testing delete with param name=adha", method: http.MethodDelete, path: "/test/p?name=adha", statusCode: test_status, expected: "adha"},
		{name: "testing patch with param name=adha", method: http.MethodPatch, path: "/test/p?name=adha", statusCode: test_status, expected: "adha"},

		{name: "testing GET route with POST method should return method not allowed 405", method: http.MethodPost, path: "/test-get", statusCode: http.StatusMethodNotAllowed, expected: "Method Not Allowed\n"},
		{name: "testing POST route with GET method should return method not allowed 405", method: http.MethodGet, path: "/test-post", statusCode: http.StatusMethodNotAllowed, expected: "Method Not Allowed\n"},
		{name: "testing PUT route with POST method should return method not allowed 405", method: http.MethodPost, path: "/test-put", statusCode: http.StatusMethodNotAllowed, expected: "Method Not Allowed\n"},
		{name: "testing DELETE route with GET method should return method not allowed 405", method: http.MethodGet, path: "/test-delete", statusCode: http.StatusMethodNotAllowed, expected: "Method Not Allowed\n"},
		{name: "testing PATCH route with PUT method should return method not allowed 405", method: http.MethodPut, path: "/test-patch", statusCode: http.StatusMethodNotAllowed, expected: "Method Not Allowed\n"},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.method, test.path, nil)
			rr := httptest.NewRecorder()
			serverHandler.ServeHTTP(rr, req)

			res := rr.Result()
			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, test.statusCode, res.StatusCode)
			assert.Equal(t, test.expected, string(body))
		})
	}

	t.Run("server returns 404 for non existent route", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/non-existent-route", nil)
		rr := httptest.NewRecorder()
		serverHandler.ServeHTTP(rr, req)

		res := rr.Result()
		body, _ := io.ReadAll(res.Body)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
		assert.Equal(t, "404 page not found\n", string(body))
	})

	t.Run("server handles same routes and different http methods", func(t *testing.T) {

		sampleHandler := func(ctx *stk.Context) {
			method := ctx.Request.Method
			ctx.Writer.WriteHeader(test_status)
			ctx.Writer.Write([]byte(method))
		}

		s.Get("/get-and-post", sampleHandler)
		s.Post("/get-and-post", sampleHandler)

		req := httptest.NewRequest(http.MethodGet, "/get-and-post", nil)
		rr := httptest.NewRecorder()
		serverHandler.ServeHTTP(rr, req)

		res := rr.Result()
		body, _ := io.ReadAll(res.Body)
		assert.Equal(t, test_status, res.StatusCode)
		assert.Equal(t, "GET", string(body))

		req = httptest.NewRequest(http.MethodPost, "/get-and-post", nil)
		rr = httptest.NewRecorder()
		serverHandler.ServeHTTP(rr, req)

		res = rr.Result()
		body, _ = io.ReadAll(res.Body)
		assert.Equal(t, test_status, res.StatusCode)
		assert.Equal(t, "POST", string(body))

	})

}

// Test middlewares

func TestMiddlewares(t *testing.T) {

	firstMiddleware := func(next stk.HandlerFunc) stk.HandlerFunc {
		return func(ctx *stk.Context) {
			ctx.Writer.Header().Add("X-FirstMiddleware", "true")
			next(ctx)
		}
	}

	secondMiddleware := func(next stk.HandlerFunc) stk.HandlerFunc {
		return func(ctx *stk.Context) {
			ctx.Writer.Header().Add("X-SecondMiddleware", "true")
			next(ctx)
		}
	}

	middlewareStatusCode := func(next stk.HandlerFunc) stk.HandlerFunc {
		return func(ctx *stk.Context) {
			ctx.Status(http.StatusBadRequest).JSONResponse("error")
		}
	}

	myHandler := func(ctx *stk.Context) {
		ctx.Status(http.StatusOK).JSONResponse("ok")
	}

	t.Run("server with two middlewares", func(t *testing.T) {
		config := &stk.ServerConfig{
			Port:           "8080",
			RequestLogging: true,
		}
		s := stk.NewServer(config)

		s.Use(firstMiddleware)
		s.Use(secondMiddleware)

		s.Get("/", myHandler)

		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		s.Router.ServeHTTP(w, req)

		resp := w.Result()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		assert.Equal(t, "true", resp.Header.Get("X-FirstMiddleware"))

		assert.Equal(t, "true", resp.Header.Get("X-SecondMiddleware"))
	})

	t.Run("server with no middlewares", func(t *testing.T) {
		config := &stk.ServerConfig{
			Port:           "8080",
			RequestLogging: true,
		}
		s := stk.NewServer(config)

		s.Get("/", myHandler)

		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		s.Router.ServeHTTP(w, req)

		resp := w.Result()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		assert.Equal(t, "", resp.Header.Get("X-FirstMiddleware"))
		assert.Equal(t, "", resp.Header.Get("X-SecondMiddleware"))
	})

	t.Run("middleware can write status code and body", func(t *testing.T) {
		config := &stk.ServerConfig{
			Port:           "8080",
			RequestLogging: true,
		}
		s := stk.NewServer(config)

		s.Use(middlewareStatusCode)
		s.Get("/", myHandler)

		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		s.Router.ServeHTTP(w, req)

		resp := w.Result()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, "\"error\"", string(body))

	})

}

func TestServerLogger(t *testing.T) {
	t.Run("Server initializes logger without passing one", func(t *testing.T) {
		config := &stk.ServerConfig{
			Port:           "8080",
			RequestLogging: true,
		}
		s := stk.NewServer(config)

		assert.NotNil(t, s.Logger)
	})
}

func TestNormalizePort(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"0.0.0.0:8080", "0.0.0.0:8080"},
		{"8080", "0.0.0.0:8080"},
		{":8080", "0.0.0.0:8080"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := stk.NormalizePort(tc.input)
			if result != tc.expected {
				t.Errorf("For input %s, expected %s, but got %s", tc.input, tc.expected, result)
			}
		})
	}
}
