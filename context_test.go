package stk_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adharshmk96/stk"
	"github.com/stretchr/testify/assert"
)

func TestStatus(t *testing.T) {
	config := &stk.ServerConfig{
		Port:           "8080",
		RequestLogging: false,
	}
	s := stk.NewServer(config)

	t.Run("status and jsonresponse methods by chaining", func(t *testing.T) {

		request, _ := http.NewRequest("GET", "/", nil)
		responseRec := httptest.NewRecorder()

		s.Get("/", func(c stk.Context) {
			c.Status(http.StatusTeapot).JSONResponse("Hello, this is a JSON response!")
		})

		s.Router.ServeHTTP(responseRec, request)

		assert.Equal(t, http.StatusTeapot, responseRec.Code)

	})

	t.Run("using status method only", func(t *testing.T) {

		request, _ := http.NewRequest("GET", "/st", nil)
		responseRec := httptest.NewRecorder()

		s.Get("/st", func(c stk.Context) {
			c.Status(http.StatusTeapot)
		})

		s.Router.ServeHTTP(responseRec, request)

		assert.Equal(t, http.StatusTeapot, responseRec.Code)

	})

	t.Run("status and json response method without chaining", func(t *testing.T) {

		request, _ := http.NewRequest("GET", "/js", nil)
		responseRec := httptest.NewRecorder()

		s.Get("/js", func(c stk.Context) {
			c.Status(http.StatusBadGateway)
			c.JSONResponse("Hello, this is a JSON response!")
		})

		s.Router.ServeHTTP(responseRec, request)

		assert.Equal(t, http.StatusBadGateway, responseRec.Code)

	})

	t.Run("json response method individually gives 200", func(t *testing.T) {

		request, _ := http.NewRequest("GET", "/jso", nil)
		responseRec := httptest.NewRecorder()

		s.Get("/jso", func(c stk.Context) {
			c.JSONResponse("Hello, this is a JSON response!")
		})

		s.Router.ServeHTTP(responseRec, request)

		assert.Equal(t, http.StatusOK, responseRec.Code)

	})

	t.Run("json response method befire status", func(t *testing.T) {

		request, _ := http.NewRequest("GET", "/json", nil)
		responseRec := httptest.NewRecorder()

		s.Get("/json", func(c stk.Context) {
			c.JSONResponse("Hello, this is a JSON response!")
			c.Status(http.StatusBadGateway)
		})

		s.Router.ServeHTTP(responseRec, request)

		assert.Equal(t, http.StatusBadGateway, responseRec.Code)

	})

}

type TestPayload struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func TestJSONResponse(t *testing.T) {

	config := &stk.ServerConfig{
		Port:           "8080",
		RequestLogging: false,
	}
	s := stk.NewServer(config)

	testCases := []struct {
		path        string
		name        string
		data        interface{}
		status      int
		expectedErr error
	}{
		{
			path:        "/",
			name:        "Invalid JSON",
			data:        make(chan int),
			status:      http.StatusInternalServerError,
			expectedErr: stk.ErrInternalServer,
		},
		{
			path: "/s",
			name: "Struct data",
			data: TestPayload{
				Message: "Hello, this is a JSON response!",
				Status:  http.StatusOK,
			},
			status:      http.StatusOK,
			expectedErr: nil,
		},
		{
			path: "/m",
			name: "Map data",
			data: map[string]interface{}{
				"message": "Hello, this is a JSON response!",
				"status":  http.StatusOK,
			},
			status:      http.StatusOK,
			expectedErr: nil,
		},
		{
			path:        "/n",
			name:        "nil",
			data:        nil,
			status:      http.StatusOK,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			request, _ := http.NewRequest("GET", tc.path, nil)
			responseRec := httptest.NewRecorder()

			s.Get(tc.path, func(c stk.Context) {
				c.Status(tc.status).JSONResponse(tc.data)
			})

			s.Router.ServeHTTP(responseRec, request)

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

	tests := []struct {
		name           string
		reqBody        string
		expectedErr    error
		expectedResult SampleStruct
	}{
		{
			name:        "Valid JSON",
			reqBody:     `{"name":"John","age":30}`,
			expectedErr: nil,
			expectedResult: SampleStruct{
				Name: "John",
				Age:  30,
			},
		},
		{
			name:           "Invalid JSON",
			reqBody:        `{"name":"John",,,"age":30}`,
			expectedErr:    stk.ErrInvalidJSON,
			expectedResult: SampleStruct{},
		},
		{
			name:           "Empty JSON",
			reqBody:        "",
			expectedErr:    stk.ErrInvalidJSON,
			expectedResult: SampleStruct{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := io.NopCloser(bytes.NewReader([]byte(tt.reqBody)))
			req := httptest.NewRequest("POST", "/", body)
			resp := httptest.NewRecorder()

			server := stk.NewServer(&stk.ServerConfig{})

			server.Post("/", func(c stk.Context) {
				var res SampleStruct
				err := c.DecodeJSONBody(&res)

				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("Expected error to be '%v', got '%v'", tt.expectedErr, err)
				}

				if err == nil && res != tt.expectedResult {
					t.Errorf("Expected result to be '%v', got '%v'", tt.expectedResult, res)
				}
			})

			server.Router.ServeHTTP(resp, req)

		})
	}
}

func TestGetAllowedOrigins(t *testing.T) {
	t.Run("returns configured origins", func(t *testing.T) {
		config := &stk.ServerConfig{
			Port:           "8080",
			RequestLogging: false,
			AllowedOrigins: []string{"http://localhost:8080", "http://localhost:8081"},
		}
		s := stk.NewServer(config)

		interMiddleware := func(next stk.HandlerFunc) stk.HandlerFunc {
			return func(c stk.Context) {
				assert.Equal(t, c.GetAllowedOrigins(), config.AllowedOrigins)
			}
		}

		s.Use(interMiddleware)

		s.Get("/", func(c stk.Context) {
			assert.Equal(t, c.GetAllowedOrigins(), config.AllowedOrigins)
		})

		request, _ := http.NewRequest("GET", "/", nil)
		responseRec := httptest.NewRecorder()

		s.Router.ServeHTTP(responseRec, request)

	})
}

func TestRawResponse(t *testing.T) {
	t.Run("sets raw response", func(t *testing.T) {
		config := &stk.ServerConfig{
			Port:           "8080",
			RequestLogging: false,
		}
		s := stk.NewServer(config)

		s.Get("/", func(c stk.Context) {
			c.RawResponse([]byte("Hello, this is a raw response!"))
		})

		request, _ := http.NewRequest("GET", "/", nil)
		responseRec := httptest.NewRecorder()

		s.Router.ServeHTTP(responseRec, request)

		assert.Equal(t, "Hello, this is a raw response!", responseRec.Body.String())
		assert.Equal(t, http.StatusOK, responseRec.Code)
	})
}

func TestGetMethod(t *testing.T) {
	t.Run("returns correct request method", func(t *testing.T) {
		config := &stk.ServerConfig{
			Port:           "8080",
			RequestLogging: false,
		}
		s := stk.NewServer(config)

		s.Get("/", func(c stk.Context) {
			assert.Equal(t, http.MethodGet, c.GetRequest().Method)
		})

		request, _ := http.NewRequest("GET", "/", nil)
		responseRec := httptest.NewRecorder()

		s.Router.ServeHTTP(responseRec, request)

	})
}

func TestSetHeader(t *testing.T) {
	t.Run("adds header to the response", func(t *testing.T) {
		config := &stk.ServerConfig{
			Port:           "8080",
			RequestLogging: false,
		}
		s := stk.NewServer(config)

		s.Get("/", func(c stk.Context) {
			c.SetHeader("X-Header", "Added")
		})

		request, _ := http.NewRequest("GET", "/", nil)
		responseRec := httptest.NewRecorder()

		s.Router.ServeHTTP(responseRec, request)

		assert.Equal(t, responseRec.Header().Get("X-Header"), "Added")
	})
}

func TestContext(t *testing.T) {
	t.Run("context is passed by reference", func(t *testing.T) {
		config := &stk.ServerConfig{
			Port:           "8080",
			RequestLogging: false,
		}
		s1 := stk.NewServer(config)
		s2 := stk.NewServer(config)

		var context1 stk.Context
		var context2 stk.Context
		var context3 stk.Context
		var context4 stk.Context

		s1.Use(func(next stk.HandlerFunc) stk.HandlerFunc {
			return func(c stk.Context) {
				context1 = c
				next(c)
			}
		})

		s2.Use(func(next stk.HandlerFunc) stk.HandlerFunc {
			return func(c stk.Context) {
				context2 = c
				next(c)
			}
		})

		s1.Get("/", func(c stk.Context) {
			context3 = c
			c.SetHeader("X-Header", "Added")
		})

		s2.Get("/", func(c stk.Context) {
			context4 = c
		})

		request, _ := http.NewRequest("GET", "/", nil)
		responseRec1 := httptest.NewRecorder()
		responseRec2 := httptest.NewRecorder()

		s1.Router.ServeHTTP(responseRec1, request)
		s2.Router.ServeHTTP(responseRec2, request)

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

func TestSetCookie(t *testing.T) {
	t.Run("adds cookie to the response", func(t *testing.T) {
		config := &stk.ServerConfig{
			Port:           "8080",
			RequestLogging: false,
		}
		s := stk.NewServer(config)

		cookie := &http.Cookie{
			Name:  "X-Cookie",
			Value: "Added",
			Path:  "/",
		}

		s.Get("/", func(c stk.Context) {
			c.SetCookie(cookie)
		})

		request, _ := http.NewRequest("GET", "/", nil)
		responseRec := httptest.NewRecorder()

		s.Router.ServeHTTP(responseRec, request)

		assert.Equal(t, responseRec.Header().Get("Set-Cookie"), "X-Cookie=Added; Path=/")
	})
}
