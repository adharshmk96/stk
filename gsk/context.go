package gsk

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type context struct {
	request *http.Request
	writer  http.ResponseWriter

	params httprouter.Params

	logger *logrus.Logger

	allowedOrigins []string
	responseStatus int
	responseBody   []byte
}

type Map map[string]interface{}

type Context interface {
	// http objects
	GetRequest() *http.Request
	GetWriter() http.ResponseWriter

	// get data from request
	GetParam(key string) string
	GetQueryParam(key string) string
	GetAllowedOrigins() []string
	DecodeJSONBody(v interface{}) error

	// set data for response
	Status(status int) Context
	SetHeader(string, string)
	RawResponse(raw []byte)
	JSONResponse(data interface{})

	// cookies
	SetCookie(cookie *http.Cookie)
	GetCookie(name string) (*http.Cookie, error)

	// logger
	Logger() *logrus.Logger
}

func (c *context) GetRequest() *http.Request {
	return c.request
}

func (c *context) GetWriter() http.ResponseWriter {
	return c.writer
}

// GetParam gets the params within the path mentioned as a wildcard
func (c *context) GetParam(key string) string {
	return c.params.ByName(key)
}

// GetQueryParam gets the query parameters passed eg: /?name=value
func (c *context) GetQueryParam(key string) string {
	return c.request.URL.Query().Get(key)
}

func (c *context) GetAllowedOrigins() []string {
	return c.allowedOrigins
}

func (c *context) DecodeJSONBody(v interface{}) error {
	bodySizeLimit := int64(1 << 20) // 1 MB

	// Set a maximum limit for the request body size to avoid possible malicious requests
	c.request.Body = http.MaxBytesReader(c.writer, c.request.Body, bodySizeLimit)

	// Manually check if the request body size exceeds the limit
	if c.request.ContentLength > bodySizeLimit {
		c.writer.Header().Set("Content-Type", "application/json")
		http.Error(c.writer, ErrBodyTooLarge.Error(), http.StatusRequestEntityTooLarge)
		return ErrBodyTooLarge
	}

	// Decode the JSON body into the provided interface
	decoder := json.NewDecoder(c.request.Body)
	err := decoder.Decode(v)

	defer c.request.Body.Close()

	// Check if there is an error in decoding the JSON, and return a user-friendly error message
	if err != nil {
		if err == io.EOF {
			return ErrInvalidJSON
		} else if _, ok := err.(*json.SyntaxError); ok {
			return ErrInvalidJSON
		}
		// If the error is not an EOF or syntax error, return the original error
		return err
	}

	return nil
}

// Status sets the status code of the response
func (c *context) Status(status int) Context {
	c.responseStatus = status
	return c
}

// sets response header key value
func (c *context) SetHeader(key string, value string) {
	c.writer.Header().Add(key, value)
}

// JSONResponse marshals the provided interface into JSON and writes it to the response writer
// If there is an error in marshalling the JSON, an internal server error is returned
func (c *context) JSONResponse(data interface{}) {
	response, err := json.Marshal(data)
	// Set the content type to JSON
	c.writer.Header().Set("Content-Type", "application/json")

	// Check if there is an error in marshalling the JSON (internal server error)
	if err != nil {
		c.responseStatus = http.StatusInternalServerError
		c.responseBody = []byte(ErrInternalServer.Error())
		return
	}

	c.responseBody = response
}

// sets response body as byte array
func (c *context) RawResponse(raw []byte) {
	c.responseBody = raw
}

func (c *context) Logger() *logrus.Logger {
	return c.logger
}

func (c *context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.writer, cookie)
}

func (c *context) GetCookie(name string) (*http.Cookie, error) {
	return c.request.Cookie(name)
}
