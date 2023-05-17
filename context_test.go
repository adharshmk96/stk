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
	"github.com/julienschmidt/httprouter"
)

func TestStatus(t *testing.T) {
	config := &stk.ServerConfig{
		Port:           "8080",
		RequestLogging: true,
	}
	s := stk.NewServer(config)

	t.Run("test status and jsonresponse methods by chaining", func(t *testing.T) {

		request, _ := http.NewRequest("GET", "/", nil)
		responseRec := httptest.NewRecorder()

		s.Get("/", func(c *stk.Context) {
			c.Status(http.StatusTeapot).JSONResponse("Hello, this is a JSON response!")
		})

		s.Router.ServeHTTP(responseRec, request)

		if responseRec.Code != http.StatusTeapot {
			t.Errorf("Expected status code to be %d but got %d", http.StatusTeapot, responseRec.Code)
		}

	})

	t.Run("test status method only", func(t *testing.T) {

		request, _ := http.NewRequest("GET", "/st", nil)
		responseRec := httptest.NewRecorder()

		s.Get("/st", func(c *stk.Context) {
			c.Status(http.StatusTeapot)
		})

		s.Router.ServeHTTP(responseRec, request)

		if responseRec.Code != http.StatusTeapot {
			t.Errorf("Expected status code to be %d but got %d", http.StatusTeapot, responseRec.Code)
		}

	})

	t.Run("test status and json response method without chaining", func(t *testing.T) {

		request, _ := http.NewRequest("GET", "/js", nil)
		responseRec := httptest.NewRecorder()

		s.Get("/js", func(c *stk.Context) {
			c.Status(http.StatusBadGateway)
			c.JSONResponse("Hello, this is a JSON response!")
		})

		s.Router.ServeHTTP(responseRec, request)

		if responseRec.Code != http.StatusBadGateway {
			t.Errorf("Expected status code to be %d but got %d", http.StatusBadGateway, responseRec.Code)
		}

	})

	t.Run("test json response method individually gives 200", func(t *testing.T) {

		request, _ := http.NewRequest("GET", "/jso", nil)
		responseRec := httptest.NewRecorder()

		s.Get("/jso", func(c *stk.Context) {
			c.JSONResponse("Hello, this is a JSON response!")
		})

		s.Router.ServeHTTP(responseRec, request)

		if responseRec.Code != http.StatusOK {
			t.Errorf("Expected status code to be %d but got %d", http.StatusOK, responseRec.Code)
		}

	})

}

type TestPayload struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func TestJSONResponse(t *testing.T) {

	config := &stk.ServerConfig{
		Port:           "8080",
		RequestLogging: true,
	}
	s := stk.NewServer(config)

	testCases := []struct {
		path        string
		name        string
		data        interface{}
		expectedErr error
	}{
		{
			path:        "/",
			name:        "Invalid JSON",
			data:        make(chan int),
			expectedErr: stk.ErrInternalServer,
		},
		{
			path: "/s",
			name: "Struct data",
			data: TestPayload{
				Message: "Hello, this is a JSON response!",
				Status:  http.StatusOK,
			},
			expectedErr: nil,
		},
		{
			path: "/m",
			name: "Map data",
			data: map[string]interface{}{
				"message": "Hello, this is a JSON response!",
				"status":  http.StatusOK,
			},
			expectedErr: nil,
		},
		{
			path:        "/n",
			name:        "nil",
			data:        nil,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			request, _ := http.NewRequest("GET", tc.path, nil)
			responseRec := httptest.NewRecorder()

			s.Get(tc.path, func(c *stk.Context) {
				c.JSONResponse(tc.data)
			})

			s.Router.ServeHTTP(responseRec, request)

			if tc.expectedErr != nil {
				expectedErr := tc.expectedErr.Error()
				if responseRec.Body.String() != string(expectedErr) {
					t.Errorf("Expected error to be %q but got %q", expectedErr, responseRec.Body.String())
				}
			} else {
				expectedJSON, _ := json.Marshal(tc.data)
				if responseRec.Body.String() != string(expectedJSON) {
					t.Errorf("Expected JSON data to be %q but got %q", expectedJSON, responseRec.Body.String())
				}
			}
		})
	}
}

// TestDecodeJSONBody tests the DecodeJSONBody method in the Context struct.
func TestDecodeJSONBody(t *testing.T) {
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

			context := stk.Context{
				Request:        req,
				Writer:         resp,
				Params:         httprouter.Params{},
				ResponseStatus: 0,
			}

			var res SampleStruct
			err := context.DecodeJSONBody(&res)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("Expected error to be '%v', got '%v'", tt.expectedErr, err)
			}

			if err == nil && res != tt.expectedResult {
				t.Errorf("Expected result to be '%v', got '%v'", tt.expectedResult, res)
			}
		})
	}
}

type SampleStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
