package gsk_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/adharshmk96/stk/gsk"
	"github.com/adharshmk96/stk/testutils"
	"github.com/stretchr/testify/assert"
)

func TestStatus(t *testing.T) {
	config := &gsk.ServerConfig{
		Port: "8888",
	}
	s := gsk.New(config)

	t.Run("sets status and jsonresponse methods by chaining", func(t *testing.T) {

		s.Get("/", func(c *gsk.Context) {
			c.Status(http.StatusTeapot).JSONResponse("Hello, this is a JSON response!")
		})

		rr, _ := s.Test("GET", "/", nil)
		assert.Equal(t, http.StatusTeapot, rr.Code)

	})

	t.Run("using status method sets http response", func(t *testing.T) {

		s.Get("/st", func(c *gsk.Context) {
			c.Status(http.StatusTeapot)
		})

		rr, _ := s.Test("GET", "/st", nil)
		assert.Equal(t, http.StatusTeapot, rr.Code)

	})

	t.Run("set status and json response without chaining", func(t *testing.T) {

		s.Get("/js", func(c *gsk.Context) {
			c.Status(http.StatusBadGateway)
			c.JSONResponse("Hello, this is a JSON response!")
		})

		rr, _ := s.Test("GET", "/js", nil)
		assert.Equal(t, http.StatusBadGateway, rr.Code)

	})

	t.Run("json response method default gives 200", func(t *testing.T) {

		s.Get("/jso", func(c *gsk.Context) {
			c.JSONResponse("Hello, this is a JSON response!")
		})

		rr, _ := s.Test("GET", "/jso", nil)
		assert.Equal(t, http.StatusOK, rr.Code)

	})

	t.Run("json response method used before status gives proper response", func(t *testing.T) {

		s.Get("/json", func(c *gsk.Context) {
			c.JSONResponse("Hello, this is a JSON response!")
			c.Status(http.StatusBadGateway)
		})

		rr, _ := s.Test("GET", "/json", nil)
		assert.Equal(t, http.StatusBadGateway, rr.Code)

	})

}

type TestPayload struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func TestJSONResponse(t *testing.T) {

	config := &gsk.ServerConfig{
		Port: "8888",
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

			s.Get(tc.path, func(c *gsk.Context) {
				c.Status(tc.status).JSONResponse(tc.data)
			})

			rr, _ := s.Test("GET", tc.path, nil)

			if tc.expectedErr != nil {
				expectedErr := tc.expectedErr.Error()
				assert.Equal(t, rr.Body.String(), expectedErr)
			} else {
				expectedJSON, _ := json.Marshal(tc.data)
				assert.Equal(t, rr.Body.String(), string(expectedJSON))
			}

			assert.Equal(t, tc.status, rr.Code)
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

	t.Run("decodes valid json", func(t *testing.T) {
		server := gsk.New(&gsk.ServerConfig{})

		server.Post("/", func(c *gsk.Context) {
			var res SampleStruct
			err := c.DecodeJSONBody(&res)

			if err != nil {
				t.Errorf("Expected error to be nil, got '%v'", err)
			}

			if res != (SampleStruct{Name: "John", Age: 30}) {
				t.Errorf("Expected result to be '%v', got '%v'", SampleStruct{Name: "John", Age: 30}, res)
			}
		})

		server.Test("POST", "/", io.NopCloser(bytes.NewReader(jsonBody)))

	})

	t.Run("returns error on invalid json", func(t *testing.T) {
		server := gsk.New(&gsk.ServerConfig{})

		server.Post("/", func(c *gsk.Context) {
			var res SampleStruct
			err := c.DecodeJSONBody(&res)

			if err == nil {
				t.Errorf("Expected error to be '%v', got '%v'", gsk.ErrInvalidJSON, err)
			}

			if res != (SampleStruct{}) {
				t.Errorf("Expected result to be '%v', got '%v'", SampleStruct{}, res)
			}
		})

		server.Test("POST", "/", io.NopCloser(bytes.NewReader([]byte(`{"name":"John",,,"age":30}`))))

	})

	t.Run("decodes to empty struct on empty json", func(t *testing.T) {
		server := gsk.New(&gsk.ServerConfig{})
		server.Post("/", func(c *gsk.Context) {
			var res SampleStruct
			err := c.DecodeJSONBody(&res)

			if err == nil {
				t.Errorf("Expected error to be '%v', got '%v'", gsk.ErrInvalidJSON, err)
			}

			if res != (SampleStruct{}) {
				t.Errorf("Expected result to be '%v', got '%v'", SampleStruct{}, res)
			}
		})

		server.Test("POST", "/", io.NopCloser(bytes.NewReader([]byte(``))))
	})

	t.Run("returns err if body is nil", func(t *testing.T) {
		server := gsk.New(&gsk.ServerConfig{})
		server.Post("/", func(c *gsk.Context) {
			var res SampleStruct
			err := c.DecodeJSONBody(&res)

			if err == nil {
				t.Errorf("Expected error to be '%v', got '%v'", gsk.ErrInvalidJSON, err)
			}

			if res != (SampleStruct{}) {
				t.Errorf("Expected result to be '%v', got '%v'", SampleStruct{}, res)
			}
		})

		server.Test("POST", "/", nil)
	})

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

			server := gsk.New(&gsk.ServerConfig{
				BodySizeLimit: tt.bodySizeLimit,
			})

			server.Post("/", func(c *gsk.Context) {
				var res SampleStruct
				err := c.DecodeJSONBody(&res)

				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("Expected error to be '%v', got '%v'", tt.expectedErr, err)
				}

				if err == nil && res != tt.expectedResult {
					t.Errorf("Expected result to be '%v', got '%v'", tt.expectedResult, res)
				}
			})

			server.Test("POST", "/", body)

		})
	}
}

func TestRawResponse(t *testing.T) {
	t.Run("sets raw response", func(t *testing.T) {
		config := &gsk.ServerConfig{
			Port: "8888",
		}
		s := gsk.New(config)

		s.Get("/", func(c *gsk.Context) {
			c.RawResponse([]byte("Hello, this is a raw response!"))
		})

		rr, _ := s.Test("GET", "/", nil)

		assert.Equal(t, "Hello, this is a raw response!", rr.Body.String())
		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

func TestGetRequestMethod(t *testing.T) {
	t.Run("returns correct request method", func(t *testing.T) {
		config := &gsk.ServerConfig{
			Port: "8888",
		}
		s := gsk.New(config)

		s.Get("/", func(c *gsk.Context) {
			assert.Equal(t, http.MethodGet, c.Request.Method)
		})

		s.Test("GET", "/", nil)

	})
}

func TestSetHeader(t *testing.T) {
	t.Run("adds header to the response", func(t *testing.T) {
		config := &gsk.ServerConfig{
			Port: "8888",
		}
		s := gsk.New(config)

		s.Get("/", func(c *gsk.Context) {
			c.SetHeader("X-Header", "Added")
		})

		rr, _ := s.Test("GET", "/", nil)

		assert.Equal(t, rr.Header().Get("X-Header"), "Added")
	})
}

func TestContext(t *testing.T) {
	t.Run("context is desn't overlap between handlers", func(t *testing.T) {
		config := &gsk.ServerConfig{
			Port: "8888",
		}
		s1 := gsk.New(config)
		s2 := gsk.New(config)

		if s1 == s2 {
			t.Errorf("Servers should be different")
		}

		var context1 *gsk.Context
		var context2 *gsk.Context
		var context3 *gsk.Context
		var context4 *gsk.Context

		s1.Use(func(next gsk.HandlerFunc) gsk.HandlerFunc {
			return func(c *gsk.Context) {
				context1 = c
				next(c)
			}
		})

		s2.Use(func(next gsk.HandlerFunc) gsk.HandlerFunc {
			return func(c *gsk.Context) {
				context2 = c
				next(c)
			}
		})

		s1.Get("/", func(c *gsk.Context) {
			context3 = c
			c.SetHeader("X-Header", "Added")
		})

		s2.Get("/", func(c *gsk.Context) {
			context4 = c
		})

		rr1, _ := s1.Test("GET", "/", nil)
		rr2, _ := s2.Test("GET", "/", nil)

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

		if rr1.Header().Get("X-Header") != "Added" {
			t.Errorf("Expected rr1 to have header 'X-Header'")
		}
		if rr2.Header().Get("X-Header") == "Added" {
			t.Errorf("Expected rr2 to not have header 'X-Header'")
		}
	})
}

func TestCookie(t *testing.T) {

	config := &gsk.ServerConfig{
		Port: "8888",
	}
	s := gsk.New(config)

	t.Run("SetCookie adds cookie to the response", func(t *testing.T) {

		cookie := &http.Cookie{
			Name:     "X-Cookie",
			Value:    "Added",
			Path:     "/",
			HttpOnly: true,
		}

		s.Get("/", func(c *gsk.Context) {
			c.SetCookie(cookie)
		})

		rr, _ := s.Test("GET", "/", nil)

		assert.Equal(t, rr.Header().Get("Set-Cookie"), "X-Cookie=Added; Path=/; HttpOnly")

	})

	t.Run("GetCookie gets cookie from the request", func(t *testing.T) {

		cookie := &http.Cookie{
			Name:  "X-Cookie",
			Value: "Added",
			Path:  "/",
		}
		testParams := gsk.TestParams{
			Cookies: []*http.Cookie{
				cookie,
			},
		}

		s.Test("GET", "/c", nil, testParams)

		s.Get("/c", func(c *gsk.Context) {
			reqCookie, _ := c.GetCookie("X-Cookie")
			assert.Equal(t, cookie.Value, reqCookie.Value)
			assert.Equal(t, cookie.Name, reqCookie.Name)
		})

	})

	t.Run("GetCookie returns error if cookie is not found", func(t *testing.T) {

		s.Test("GET", "/ce", nil)

		s.Get("/ce", func(c *gsk.Context) {
			_, err := c.GetCookie("X-Cookie")
			assert.Error(t, err)
		})

	})

}

func TestHelperFunctions(t *testing.T) {
	t.Run("test origin", func(t *testing.T) {
		config := &gsk.ServerConfig{
			Port: "8888",
		}
		s := gsk.New(config)

		s.Get("/", func(c *gsk.Context) {
			assert.Equal(t, "http://example.com", c.Origin())
		})

		s.Test("GET", "/", nil, gsk.TestParams{
			Headers: map[string]string{
				"Origin": "http://example.com",
			},
		})

	})

	t.Run("logger is returned", func(t *testing.T) {
		config := &gsk.ServerConfig{
			Port: "8888",
		}
		s := gsk.New(config)

		s.Get("/", func(c *gsk.Context) {
			assert.NotNil(t, c.Logger())
		})

		s.Test("GET", "/", nil)
	})

	t.Run("redirect method redirects to the given url", func(t *testing.T) {
		config := &gsk.ServerConfig{
			Port: "8888",
		}
		s := gsk.New(config)

		s.Get("/", func(c *gsk.Context) {
			c.Redirect("http://example.com")
		})

		rr, _ := s.Test("GET", "/", nil)

		assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
		assert.Equal(t, "http://example.com", rr.Header().Get("Location"))
	})

	t.Run("get cookie returns cookie from the request", func(t *testing.T) {
		config := &gsk.ServerConfig{
			Port: "8888",
		}
		s := gsk.New(config)

		s.Get("/", func(c *gsk.Context) {
			cookie, err := c.GetCookie("X-Cookie")
			assert.NoError(t, err)
			assert.Equal(t, "Added", cookie.Value)
		})

		cookie := &http.Cookie{
			Name:  "X-Cookie",
			Value: "Added",
			Path:  "/",
		}
		testParams := gsk.TestParams{
			Cookies: []*http.Cookie{
				cookie,
			},
		}

		s.Test("GET", "/", nil, testParams)
	})

	t.Run("get status returns status code of the response", func(t *testing.T) {
		config := &gsk.ServerConfig{
			Port: "8888",
		}
		s := gsk.New(config)

		s.Get("/", func(c *gsk.Context) {
			c.Status(http.StatusTeapot)
			statusCode := c.GetStatusCode()
			assert.Equal(t, http.StatusTeapot, statusCode)
		})

		s.Test("GET", "/", nil)

	})
}

func TestTemplateResponse(t *testing.T) {
	t.Run("renders template with data", func(t *testing.T) {
		templateContent := `<html>
	<head>
		<title>{{ .Var.Title }}</title>
	</head>
	<body>
		<h1>{{ .Var.Title }}</h1>
		<p>{{ .Var.Body }}</p>
	</body>
</html>`

		tempFile, removeDir := testutils.CreateTempFile(t, templateContent)
		defer removeDir()

		config := &gsk.ServerConfig{
			Port: "8888",
		}
		s := gsk.New(config)

		s.Get("/", func(c *gsk.Context) {
			c.Status(http.StatusOK).TemplateResponse(&gsk.Tpl{
				TemplatePath: tempFile,
				Variables: struct {
					Title string
					Body  string
				}{
					Title: "Hello",
					Body:  "World",
				},
			})
		})

		rr, _ := s.Test("GET", "/", nil)

		expected := `<html>
	<head>
		<title>Hello</title>
	</head>
	<body>
		<h1>Hello</h1>
		<p>World</p>
	</body>
</html>`

		assert.Equal(t, expected, rr.Body.String())
		assert.Equal(t, http.StatusOK, rr.Code)

	})

	t.Run("returns error if template is not found", func(t *testing.T) {
		config := &gsk.ServerConfig{
			Port: "8888",
		}
		s := gsk.New(config)

		s.Get("/", func(c *gsk.Context) {
			c.Status(http.StatusOK).TemplateResponse(&gsk.Tpl{
				TemplatePath: "/tmp/nonext",
				Variables: struct {
					Title string
					Body  string
				}{
					Title: "Hello",
					Body:  "World",
				},
			})
		})

		rr, _ := s.Test("GET", "/", nil)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}
