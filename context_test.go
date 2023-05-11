package stk_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adharshmk96/stk"
	"github.com/julienschmidt/httprouter"
)

type TestPayload struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func TestJSONResponse(t *testing.T) {
	testCases := []struct {
		name        string
		data        interface{}
		expectedErr error
	}{
		{
			name:        "Invalid JSON",
			data:        make(chan int),
			expectedErr: stk.ErrInternalServer,
		},
		{
			name: "Struct data",
			data: TestPayload{
				Message: "Hello, this is a JSON response!",
				Status:  http.StatusOK,
			},
			expectedErr: nil,
		},
		{
			name: "Map data",
			data: map[string]interface{}{
				"message": "Hello, this is a JSON response!",
				"status":  http.StatusOK,
			},
			expectedErr: nil,
		},
		{
			name:        "nil",
			data:        nil,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request, _ := http.NewRequest("GET", "/", nil)
			responseRec := httptest.NewRecorder()

			router := httprouter.New()
			router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
				context := &stk.Context{
					Request: r,
					Writer:  w,
				}
				context.JSONResponse(tc.data)
			})

			router.ServeHTTP(responseRec, request)

			if tc.expectedErr != nil {
				expectedErr := tc.expectedErr.Error() + "\n"
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
			body := ioutil.NopCloser(bytes.NewReader([]byte(tt.reqBody)))
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
