package gsk_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adharshmk96/stk/gsk"
	"github.com/stretchr/testify/assert"
)

func TestStatus(t *testing.T) {
	config := &gsk.ServerConfig{
		Port: "8080",
	}
	s := gsk.New(config)

	t.Run("sets status and jsonresponse methods by chaining", func(t *testing.T) {

		request, _ := http.NewRequest("GET", "/", nil)
		responseRec := httptest.NewRecorder()

		s.Get("/", func(c gsk.Context) {
			c.Status(http.StatusTeapot).JSONResponse("Hello, this is a JSON response!")
		})

		s.GetRouter().ServeHTTP(responseRec, request)

		assert.Equal(t, http.StatusTeapot, responseRec.Code)

	})

	t.Run("using status method sets http response", func(t *testing.T) {

		request, _ := http.NewRequest("GET", "/st", nil)
		responseRec := httptest.NewRecorder()

		s.Get("/st", func(c gsk.Context) {
			c.Status(http.StatusTeapot)
		})

		s.GetRouter().ServeHTTP(responseRec, request)

		assert.Equal(t, http.StatusTeapot, responseRec.Code)

	})

	t.Run("set status and json response without chaining", func(t *testing.T) {

		request, _ := http.NewRequest("GET", "/js", nil)
		responseRec := httptest.NewRecorder()

		s.Get("/js", func(c gsk.Context) {
			c.Status(http.StatusBadGateway)
			c.JSONResponse("Hello, this is a JSON response!")
		})

		s.GetRouter().ServeHTTP(responseRec, request)

		assert.Equal(t, http.StatusBadGateway, responseRec.Code)

	})

	t.Run("json response method default gives 200", func(t *testing.T) {

		request, _ := http.NewRequest("GET", "/jso", nil)
		responseRec := httptest.NewRecorder()

		s.Get("/jso", func(c gsk.Context) {
			c.JSONResponse("Hello, this is a JSON response!")
		})

		s.GetRouter().ServeHTTP(responseRec, request)

		assert.Equal(t, http.StatusOK, responseRec.Code)

	})

	t.Run("json response method used before status gives proper response", func(t *testing.T) {

		request, _ := http.NewRequest("GET", "/json", nil)
		responseRec := httptest.NewRecorder()

		s.Get("/json", func(c gsk.Context) {
			c.JSONResponse("Hello, this is a JSON response!")
			c.Status(http.StatusBadGateway)
		})

		s.GetRouter().ServeHTTP(responseRec, request)

		assert.Equal(t, http.StatusBadGateway, responseRec.Code)

	})

}

type TestPayload struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func TestJSONResponse(t *testing.T) {

	config := &gsk.ServerConfig{
		Port: "8080",
	}
	s := gsk.New(config)

	testCases := []struct {
		path        string
		name        string
		data        interface{}
		status      int
		expectedErr error
	}{
		{
			path:        "/",
			name:        "returns error for invalid data",
			data:        make(chan int),
			status:      http.StatusInternalServerError,
			expectedErr: gsk.ErrInternalServer,
		},
		{
			path: "/s",
			name: "structure data is parsed",
			data: TestPayload{
				Message: "Hello, this is a JSON response!",
				Status:  http.StatusOK,
			},
			status:      http.StatusOK,
			expectedErr: nil,
		},
		{
			path: "/m",
			name: "map data is parsed",
			data: map[string]interface{}{
				"message": "Hello, this is a JSON response!",
				"status":  http.StatusOK,
			},
			status:      http.StatusOK,
			expectedErr: nil,
		},
		{
			path:        "/n",
			name:        "nil data is parsed",
			data:        nil,
			status:      http.StatusOK,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			request, _ := http.NewRequest("GET", tc.path, nil)
			responseRec := httptest.NewRecorder()

			s.Get(tc.path, func(c gsk.Context) {
				c.Status(tc.status).JSONResponse(tc.data)
			})

			s.GetRouter().ServeHTTP(responseRec, request)

			if tc.expectedErr != nil {
				expectedErr := tc.expectedErr.Error()
				assert.Equal(t, responseRec.Body.String(), expectedErr)
			} else {
				expectedJSON, _ := json.Marshal(tc.data)
				assert.Equal(t, responseRec.Body.String(), string(expectedJSON))
			}

			assert.Equal(t, tc.status, responseRec.Code)
		})
	}
}

// TestDecodeJSONBody tests the DecodeJSONBody method in the Context struct.
func TestDecodeJSONBody(t *testing.T) {

	type SampleStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	jsonBody, _ := json.Marshal(SampleStruct{Name: "John", Age: 30})

	tests := []struct {
		name           string
		reqBody        string
		expectedErr    error
		expectedResult SampleStruct
	}{
		{
			name:        "decodes valid json",
			reqBody:     string(jsonBody),
			expectedErr: nil,
			expectedResult: SampleStruct{
				Name: "John",
				Age:  30,
			},
		},
		{
			name:           "returns error on invalid json",
			reqBody:        `{"name":"John",,,"age":30}`,
			expectedErr:    gsk.ErrInvalidJSON,
			expectedResult: SampleStruct{},
		},
		{
			name:           "decodes to empty struct on empty json",
			reqBody:        "",
			expectedErr:    gsk.ErrInvalidJSON,
			expectedResult: SampleStruct{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := io.NopCloser(bytes.NewReader([]byte(tt.reqBody)))
			req := httptest.NewRequest("POST", "/", body)
			resp := httptest.NewRecorder()

			server := gsk.New(&gsk.ServerConfig{})

			server.Post("/", func(c gsk.Context) {
				var res SampleStruct
				err := c.DecodeJSONBody(&res)

				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("Expected error to be '%v', got '%v'", tt.expectedErr, err)
				}

				if err == nil && res != tt.expectedResult {
					t.Errorf("Expected result to be '%v', got '%v'", tt.expectedResult, res)
				}
			})

			server.GetRouter().ServeHTTP(resp, req)

		})
	}
}

// generate10KBArray generates a byte array of 10KB size
func generate2MB() string {
	size := 1 << 20 * 2 // 1KB is 1 << 10, so 10KB is 10 times that
	b := make([]byte, size)

	// Let's fill the array with some value, e.g. 1
	for i := range b {
		b[i] = 'a'
	}

	return string(b)
}

func TestDecodeJSONBodySizeLimit(t *testing.T) {
	type SampleStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	tests := []struct {
		name           string
		reqBody        string
		bodySizeLimit  int64
		expectedErr    error
		expectedResult SampleStruct
	}{
		{
			name:          "decodes valid json",
			reqBody:       generate2MB(),
			bodySizeLimit: 1,
			expectedErr:   gsk.ErrBodyTooLarge,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := io.NopCloser(bytes.NewReader([]byte(tt.reqBody)))
			req := httptest.NewRequest("POST", "/", body)
			resp := httptest.NewRecorder()

			server := gsk.New(&gsk.ServerConfig{
				BodySizeLimit: tt.bodySizeLimit,
			})

			server.Post("/", func(c gsk.Context) {
				var res SampleStruct
				err := c.DecodeJSONBody(&res)

				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("Expected error to be '%v', got '%v'", tt.expectedErr, err)
				}

				if err == nil && res != tt.expectedResult {
					t.Errorf("Expected result to be '%v', got '%v'", tt.expectedResult, res)
				}
			})

			server.GetRouter().ServeHTTP(resp, req)

		})
	}
}

func TestRawResponse(t *testing.T) {
	t.Run("sets raw response", func(t *testing.T) {
		config := &gsk.ServerConfig{
			Port: "8080",
		}
		s := gsk.New(config)

		s.Get("/", func(c gsk.Context) {
			c.RawResponse([]byte("Hello, this is a raw response!"))
		})

		request, _ := http.NewRequest("GET", "/", nil)
		responseRec := httptest.NewRecorder()

		s.GetRouter().ServeHTTP(responseRec, request)

		assert.Equal(t, "Hello, this is a raw response!", responseRec.Body.String())
		assert.Equal(t, http.StatusOK, responseRec.Code)
	})
}

func TestGetRequestMethod(t *testing.T) {
	t.Run("returns correct request method", func(t *testing.T) {
		config := &gsk.ServerConfig{
			Port: "8080",
		}
		s := gsk.New(config)

		s.Get("/", func(c gsk.Context) {
			assert.Equal(t, http.MethodGet, c.GetRequest().Method)
		})

		request, _ := http.NewRequest("GET", "/", nil)
		responseRec := httptest.NewRecorder()

		s.GetRouter().ServeHTTP(responseRec, request)

	})
}

func TestSetHeader(t *testing.T) {
	t.Run("adds header to the response", func(t *testing.T) {
		config := &gsk.ServerConfig{
			Port: "8080",
		}
		s := gsk.New(config)

		s.Get("/", func(c gsk.Context) {
			c.SetHeader("X-Header", "Added")
		})

		request, _ := http.NewRequest("GET", "/", nil)
		responseRec := httptest.NewRecorder()

		s.GetRouter().ServeHTTP(responseRec, request)

		assert.Equal(t, responseRec.Header().Get("X-Header"), "Added")
	})
}

func TestContext(t *testing.T) {
	t.Run("context is desn't overlap between handlers", func(t *testing.T) {
		config := &gsk.ServerConfig{
			Port: "8080",
		}
		s1 := gsk.New(config)
		s2 := gsk.New(config)

		if s1 == s2 {
			t.Errorf("Servers should be different")
		}

		var context1 gsk.Context
		var context2 gsk.Context
		var context3 gsk.Context
		var context4 gsk.Context

		s1.Use(func(next gsk.HandlerFunc) gsk.HandlerFunc {
			return func(c gsk.Context) {
				context1 = c
				next(c)
			}
		})

		s2.Use(func(next gsk.HandlerFunc) gsk.HandlerFunc {
			return func(c gsk.Context) {
				context2 = c
				next(c)
			}
		})

		s1.Get("/", func(c gsk.Context) {
			context3 = c
			c.SetHeader("X-Header", "Added")
		})

		s2.Get("/", func(c gsk.Context) {
			context4 = c
		})

		request, _ := http.NewRequest("GET", "/", nil)
		responseRec1 := httptest.NewRecorder()
		responseRec2 := httptest.NewRecorder()

		s1.GetRouter().ServeHTTP(responseRec1, request)
		s2.GetRouter().ServeHTTP(responseRec2, request)

		if context1 != context3 {
			t.Errorf("Expected context1 to be the same as context3")
		}
		if context2 != context4 {
			t.Errorf("Expected context2 to be the same as context4")
		}

		if context1 == context2 {
			t.Errorf("Expected context1 to be different from context2")
		}
		if context3 == context4 {
			t.Errorf("Expected context3 to be different from context4")
		}

		if responseRec1.Header().Get("X-Header") != "Added" {
			t.Errorf("Expected responseRec1 to have header 'X-Header'")
		}
		if responseRec2.Header().Get("X-Header") == "Added" {
			t.Errorf("Expected responseRec2 to not have header 'X-Header'")
		}
	})
}

func TestCookie(t *testing.T) {

	config := &gsk.ServerConfig{
		Port: "8080",
	}
	s := gsk.New(config)

	t.Run("SetCookie adds cookie to the response", func(t *testing.T) {

		cookie := &http.Cookie{
			Name:     "X-Cookie",
			Value:    "Added",
			Path:     "/",
			HttpOnly: true,
		}

		s.Get("/", func(c gsk.Context) {
			c.SetCookie(cookie)
		})

		request, _ := http.NewRequest("GET", "/", nil)
		responseRec := httptest.NewRecorder()

		s.GetRouter().ServeHTTP(responseRec, request)

		assert.Equal(t, responseRec.Header().Get("Set-Cookie"), "X-Cookie=Added; Path=/; HttpOnly")

	})

	t.Run("GetCookie gets cookie from the request", func(t *testing.T) {

		cookie := &http.Cookie{
			Name:  "X-Cookie",
			Value: "Added",
			Path:  "/",
		}

		request, _ := http.NewRequest("GET", "/c", nil)
		request.AddCookie(cookie)
		responseRec := httptest.NewRecorder()

		s.Get("/c", func(c gsk.Context) {
			reqCookie, _ := c.GetCookie("X-Cookie")
			assert.Equal(t, cookie.Value, reqCookie.Value)
			assert.Equal(t, cookie.Name, reqCookie.Name)
		})

		s.GetRouter().ServeHTTP(responseRec, request)

	})

	t.Run("GetCookie returns error if cookie is not found", func(t *testing.T) {

		request, _ := http.NewRequest("GET", "/ce", nil)
		responseRec := httptest.NewRecorder()

		s.Get("/ce", func(c gsk.Context) {
			_, err := c.GetCookie("X-Cookie")
			assert.Error(t, err)
		})

		s.GetRouter().ServeHTTP(responseRec, request)

	})

}
