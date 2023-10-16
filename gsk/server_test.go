package gsk_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/adharshmk96/stk/gsk"
)

// func TestServer_Start(t *testing.T) {
// 	buffer := new(bytes.Buffer)
// 	logger := slog.New(slog.NewJSONHandler(buffer, nil))

// 	// Create a ServerConfig
// 	config := &gsk.ServerConfig{
// 		Port:   "8888",
// 		Logger: logger,
// 	}

// 	// Create a new server
// 	server := gsk.New(config)

// 	// Start the server in a goroutine
// 	go server.Start()

// 	// Wait for the server to start
// 	time.Sleep(2 * time.Second)

// 	// Check if the server is running
// 	_, err := net.DialTimeout("tcp", "localhost:"+config.Port, 1*time.Second)
// 	if err != nil {
// 		t.Fatalf("Expected server to start, but it didn't: %v", err)
// 	}
// }

// Test server routes

func TestServerRoutes(t *testing.T) {
	config := &gsk.ServerConfig{
		Port: "8888",
	}
	s := gsk.New(config)

	test_status := http.StatusNoContent

	sampleHandler := func(gc *gsk.Context) {
		method := gc.Request.Method
		gc.Status(test_status)
		gc.RawResponse([]byte(method))
	}

	s.Get("/test-get", sampleHandler)
	s.Post("/test-post", sampleHandler)
	s.Put("/test-put", sampleHandler)
	s.Delete("/test-delete", sampleHandler)
	s.Patch("/test-patch", sampleHandler)

	queryParamHandler := func(gc *gsk.Context) {
		gc.Status(test_status)
		gc.RawResponse([]byte(gc.QueryParam("name")))
	}

	s.Get("/test/p", queryParamHandler)
	s.Post("/test/p", queryParamHandler)
	s.Put("/test/p", queryParamHandler)
	s.Delete("/test/p", queryParamHandler)
	s.Patch("/test/p", queryParamHandler)

	paramsHandler := func(gc *gsk.Context) {
		gc.Status(test_status)
		gc.RawResponse([]byte(gc.Param("id")))
	}

	s.Get("/test/d/:id", paramsHandler)
	s.Post("/test/d/:id", paramsHandler)
	s.Put("/test/d/:id", paramsHandler)
	s.Delete("/test/d/:id", paramsHandler)
	s.Patch("/test/d/:id", paramsHandler)

	testCases := []struct {
		name       string
		method     string
		path       string
		statusCode int
		expected   string
	}{
		{name: "get returns 200", method: http.MethodGet, path: "/test-get", statusCode: test_status, expected: "GET"},
		{name: "post returns 200", method: http.MethodPost, path: "/test-post", statusCode: test_status, expected: "POST"},
		{name: "put returns 200", method: http.MethodPut, path: "/test-put", statusCode: test_status, expected: "PUT"},
		{name: "delete returns 200", method: http.MethodDelete, path: "/test-delete", statusCode: test_status, expected: "DELETE"},
		{name: "patch returns 200", method: http.MethodPatch, path: "/test-patch", statusCode: test_status, expected: "PATCH"},

		{name: "get with dynamic route", method: http.MethodGet, path: "/test/d/123", statusCode: test_status, expected: "123"},
		{name: "post with dynamic route", method: http.MethodPost, path: "/test/d/123", statusCode: test_status, expected: "123"},
		{name: "put with dynamic route", method: http.MethodPut, path: "/test/d/123", statusCode: test_status, expected: "123"},
		{name: "delete with dynamic route", method: http.MethodDelete, path: "/test/d/123", statusCode: test_status, expected: "123"},
		{name: "patch with dynamic route", method: http.MethodPatch, path: "/test/d/123", statusCode: test_status, expected: "123"},

		{name: "get with param name=adha", method: http.MethodGet, path: "/test/p?name=adha", statusCode: test_status, expected: "adha"},
		{name: "post with param name=adha", method: http.MethodPost, path: "/test/p?name=adha", statusCode: test_status, expected: "adha"},
		{name: "put with param name=adha", method: http.MethodPut, path: "/test/p?name=adha", statusCode: test_status, expected: "adha"},
		{name: "delete with param name=adha", method: http.MethodDelete, path: "/test/p?name=adha", statusCode: test_status, expected: "adha"},
		{name: "patch with param name=adha", method: http.MethodPatch, path: "/test/p?name=adha", statusCode: test_status, expected: "adha"},

		{name: "GET route with POST method should return method not allowed 405", method: http.MethodPost, path: "/test-get", statusCode: http.StatusMethodNotAllowed, expected: "Method Not Allowed\n"},
		{name: "POST route with GET method should return method not allowed 405", method: http.MethodGet, path: "/test-post", statusCode: http.StatusMethodNotAllowed, expected: "Method Not Allowed\n"},
		{name: "PUT route with POST method should return method not allowed 405", method: http.MethodPost, path: "/test-put", statusCode: http.StatusMethodNotAllowed, expected: "Method Not Allowed\n"},
		{name: "DELETE route with GET method should return method not allowed 405", method: http.MethodGet, path: "/test-delete", statusCode: http.StatusMethodNotAllowed, expected: "Method Not Allowed\n"},
		{name: "PATCH route with PUT method should return method not allowed 405", method: http.MethodPut, path: "/test-patch", statusCode: http.StatusMethodNotAllowed, expected: "Method Not Allowed\n"},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			rr, _ := s.Test(test.method, test.path, nil)
			res := rr.Result()
			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, test.statusCode, res.StatusCode)
			assert.Equal(t, test.expected, string(body))
		})
	}

	t.Run("server returns 404 for non existent route", func(t *testing.T) {
		rr, _ := s.Test(http.MethodGet, "/non-existent-route", nil)
		res := rr.Result()
		body, _ := io.ReadAll(res.Body)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
		assert.Equal(t, "404 page not found\n", string(body))
	})

	t.Run("server handles same routes and different http methods", func(t *testing.T) {

		sampleHandler := func(gc *gsk.Context) {
			method := gc.Request.Method
			gc.Status(test_status)
			gc.RawResponse([]byte(method))
		}

		s.Get("/get-and-post", sampleHandler)
		s.Post("/get-and-post", sampleHandler)

		rr, _ := s.Test(http.MethodGet, "/get-and-post", nil)

		res := rr.Result()
		body, _ := io.ReadAll(res.Body)
		assert.Equal(t, test_status, res.StatusCode)
		assert.Equal(t, "GET", string(body))

		rr, _ = s.Test(http.MethodPost, "/get-and-post", nil)
		res = rr.Result()
		body, _ = io.ReadAll(res.Body)
		assert.Equal(t, test_status, res.StatusCode)
		assert.Equal(t, "POST", string(body))

	})

	t.Run("server handles same routes and different http methods with dynamic routes", func(t *testing.T) {

		sampleHandler := func(gc *gsk.Context) {
			method := gc.Request.Method
			gc.Status(test_status)
			gc.RawResponse([]byte(method))
		}

		s.Get("/get-and-post/:id", sampleHandler)
		s.Post("/get-and-post/:id", sampleHandler)

		rr, _ := s.Test(http.MethodGet, "/get-and-post/123", nil)
		res := rr.Result()
		body, _ := io.ReadAll(res.Body)
		assert.Equal(t, test_status, res.StatusCode)
		assert.Equal(t, "GET", string(body))

		rr, _ = s.Test(http.MethodPost, "/get-and-post/123", nil)

		res = rr.Result()
		body, _ = io.ReadAll(res.Body)
		assert.Equal(t, test_status, res.StatusCode)
		assert.Equal(t, "POST", string(body))

	})

	t.Run("server with overlapping routes", func(t *testing.T) {

		getHandler := func(gc *gsk.Context) {
			response := "get"
			gc.Status(http.StatusOK)
			gc.RawResponse([]byte(response))
		}

		getThatHandler := func(gc *gsk.Context) {
			response := "get-that"
			gc.Status(http.StatusOK)
			gc.RawResponse([]byte(response))
		}

		s.Get("/get", getHandler)
		s.Get("/get/:that", getThatHandler)

		rr, _ := s.Test(http.MethodGet, "/get", nil)

		res := rr.Result()
		body, _ := io.ReadAll(res.Body)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "get", string(body))

		rr, _ = s.Test(http.MethodGet, "/get/that", nil)

		res = rr.Result()
		body, _ = io.ReadAll(res.Body)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "get-that", string(body))

	})

}

// Test middlewares

func TestMiddlewares(t *testing.T) {

	firstMiddleware := func(next gsk.HandlerFunc) gsk.HandlerFunc {
		return func(gc *gsk.Context) {
			gc.SetHeader("X-FirstMiddleware", "true")
			next(gc)
		}
	}

	secondMiddleware := func(next gsk.HandlerFunc) gsk.HandlerFunc {
		return func(gc *gsk.Context) {
			gc.SetHeader("X-SecondMiddleware", "true")
			next(gc)
		}
	}

	middlewareStatusCode := func(next gsk.HandlerFunc) gsk.HandlerFunc {
		return func(gc *gsk.Context) {
			gc.Status(http.StatusBadRequest).JSONResponse("error")
		}
	}

	myHandler := func(gc *gsk.Context) {
		gc.Status(http.StatusOK).JSONResponse("ok")
	}

	laterMiddleware := func(next gsk.HandlerFunc) gsk.HandlerFunc {
		return func(gc *gsk.Context) {
			gc.SetHeader("X-LaterMiddleware", "true")
			next(gc)
		}
	}

	t.Run("server with two middlewares", func(t *testing.T) {
		config := &gsk.ServerConfig{
			Port: "8888",
		}
		s := gsk.New(config)

		s.Use(firstMiddleware)
		s.Use(secondMiddleware)

		s.Get("/", myHandler)

		w, _ := s.Test("GET", "/", nil)

		resp := w.Result()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		assert.Equal(t, "true", resp.Header.Get("X-FirstMiddleware"))

		assert.Equal(t, "true", resp.Header.Get("X-SecondMiddleware"))
	})

	t.Run("server with no middlewares", func(t *testing.T) {
		config := &gsk.ServerConfig{
			Port: "8888",
		}
		s := gsk.New(config)

		s.Get("/", myHandler)

		w, _ := s.Test("GET", "/", nil)

		resp := w.Result()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		assert.Equal(t, "", resp.Header.Get("X-FirstMiddleware"))
		assert.Equal(t, "", resp.Header.Get("X-SecondMiddleware"))
	})

	t.Run("middleware can write status code and body", func(t *testing.T) {
		config := &gsk.ServerConfig{
			Port: "8888",
		}
		s := gsk.New(config)

		s.Use(middlewareStatusCode)
		s.Get("/", myHandler)

		rr, _ := s.Test("GET", "/", nil)
		resp := rr.Result()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, "\"error\"", string(body))

	})

	t.Run("middleware blocks certain routes", func(t *testing.T) {
		blockerMiddleware := func(next gsk.HandlerFunc) gsk.HandlerFunc {
			return func(gc *gsk.Context) {
				if gc.Request.URL.Path == "/blocked" {
					gc.Status(http.StatusForbidden).JSONResponse("blocked")
					return
				}
				next(gc)
			}
		}

		config := &gsk.ServerConfig{
			Port: "8888",
		}

		s := gsk.New(config)

		s.Use(blockerMiddleware)
		s.Get("/", myHandler)
		s.Get("/blocked", myHandler)

		rr, _ := s.Test("GET", "/", nil)
		resp := rr.Result()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		rr, _ = s.Test("GET", "/blocked", nil)

		resp = rr.Result()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)

		assert.Equal(t, "\"blocked\"", string(body))
	})

	t.Run("middleware will be applied only for routes below it", func(t *testing.T) {
		config := &gsk.ServerConfig{
			Port: "8888",
		}
		s := gsk.New(config)

		s.Use(firstMiddleware)
		s.Use(secondMiddleware)

		s.Get("/", myHandler)

		w, _ := s.Test("GET", "/", nil)

		s.Use(laterMiddleware)

		s.Get("/later", myHandler)

		w2, _ := s.Test("GET", "/later", nil)

		resp := w.Result()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "true", resp.Header.Get("X-FirstMiddleware"))
		assert.Equal(t, "true", resp.Header.Get("X-SecondMiddleware"))
		assert.Equal(t, "", resp.Header.Get("X-LaterMiddleware"))

		resp2 := w2.Result()

		assert.Equal(t, http.StatusOK, resp2.StatusCode)
		assert.Equal(t, "true", resp2.Header.Get("X-FirstMiddleware"))
		assert.Equal(t, "true", resp2.Header.Get("X-SecondMiddleware"))
		assert.Equal(t, "true", resp2.Header.Get("X-LaterMiddleware"))
	})

}

func TestServerLogger(t *testing.T) {
	t.Run("Server initializes logger without passing one", func(t *testing.T) {
		config := &gsk.ServerConfig{
			Port: "8888",
		}
		s := gsk.New(config)

		s.Get("/", func(gc *gsk.Context) {
			assert.NotNil(t, gc.Logger())
			gc.Status(http.StatusOK).JSONResponse("ok")
		})

	})
}

func TestNormalizePort(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"0.0.0.0:8888", "0.0.0.0:8888"},
		{"8888", "0.0.0.0:8888"},
		{":8888", "0.0.0.0:8888"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := gsk.NormalizePort(tc.input)
			if result != tc.expected {
				t.Errorf("For input %s, expected %s, but got %s", tc.input, tc.expected, result)
			}
		})
	}
}

func TestServer_Test(t *testing.T) {
	config := &gsk.ServerConfig{
		Port: "8888",
	}
	s := gsk.New(config)

	s.Get("/", func(gc *gsk.Context) {
		gc.Status(http.StatusAccepted).JSONResponse("ok")
	})

	w, err := s.Test("GET", "/", nil)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	assert.Equal(t, http.StatusAccepted, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	body, _ := io.ReadAll(w.Body)
	assert.Equal(t, "\"ok\"", string(body))
}

// func setupTestDir() {
// 	os.MkdirAll("./testdata", os.ModePerm)
// 	os.Create("./testdata/test.txt")
// }

// func teardownTestDir() {
// 	os.RemoveAll("./testdata")
// }

// func TestServer_Static(t *testing.T) {
// 	setupTestDir()
// 	defer teardownTestDir()

// 	config := &gsk.ServerConfig{
// 		Port: "8888",
// 	}
// 	s := gsk.New(config)

// 	s.Static("/static/*filepath", "./testdata")

// 	w, err := s.Test("GET", "/static/test.txt", nil)
// 	if err != nil {
// 		t.Errorf("Expected no error, but got %v", err)
// 	}

// 	assert.Equal(t, http.StatusOK, w.Code)
// 	assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))

// 	// body, _ := io.ReadAll(w.Body)
// 	// assert.Equal(t, "test", string(body))
// }
