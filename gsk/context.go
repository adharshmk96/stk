package gsk

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

type gskContext struct {
	// server
	request *http.Request
	writer  http.ResponseWriter

	// request params
	params Params

	// logging
	logger        *logrus.Logger
	bodySizeLimit int64

	// to write response
	responseStatus int
	responseBody   []byte
}

type Map map[string]interface{}

type Context interface {
	// http objects
	GetRequest() *http.Request
	GetWriter() http.ResponseWriter
	GetPath() string

	Origin() string

	// get data from request
	GetStatusCode() int
	GetParam(key string) string
	GetQueryParam(key string) string
	DecodeJSONBody(v interface{}) error

	// set data for response
	Status(status int) Context
	SetHeader(string, string)
	RawResponse(raw []byte)
	JSONResponse(data interface{})
	StringResponse(data string)

	// cookies
	SetCookie(cookie *http.Cookie)
	GetCookie(name string) (*http.Cookie, error)

	// logger
	Logger() *logrus.Logger

	// Internals
	eject() gskContext
}

func (c *gskContext) GetRequest() *http.Request {
	return c.request
}

func (c *gskContext) GetWriter() http.ResponseWriter {
	return c.writer
}

func (c *gskContext) GetPath() string {
	return c.request.URL.Path
}

func (c *gskContext) Origin() string {
	return c.request.Header.Get("Origin")
}

// GetParam gets the params within the path mentioned as a wildcard
func (c *gskContext) GetParam(key string) string {
	return c.params.ByName(key)
}

// GetQueryParam gets the query parameters passed eg: /?name=value
func (c *gskContext) GetQueryParam(key string) string {
	return c.request.URL.Query().Get(key)
}

func (c *gskContext) DecodeJSONBody(v interface{}) error {
	// TODO: config from server
	bodySizeLimit := int64(c.bodySizeLimit << 20) // 1 MB

	if c.request.Body == nil {
		return ErrInvalidJSON
	}

	// Set a maximum limit for the request body size to avoid possible malicious requests
	c.request.Body = http.MaxBytesReader(c.writer, c.request.Body, bodySizeLimit)

	// Manually check if the request body size exceeds the limit
	body, err := io.ReadAll(c.request.Body)
	if err != nil {
		c.writer.Header().Set("Content-Type", "application/json")
		http.Error(c.writer, ErrBodyTooLarge.Error(), http.StatusRequestEntityTooLarge)
		return ErrBodyTooLarge
	}

	// Decode the JSON body into the provided interface
	err = json.Unmarshal(body, v)

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
func (c *gskContext) Status(status int) Context {
	c.responseStatus = status
	return c
}

// sets response header key value
func (c *gskContext) SetHeader(key string, value string) {
	c.writer.Header().Add(key, value)
}

// JSONResponse marshals the provided interface into JSON and writes it to the response writer
// If there is an error in marshalling the JSON, an internal server error is returned
func (c *gskContext) JSONResponse(data interface{}) {
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

// StringResponse writes the provided string to the response writer
func (c *gskContext) StringResponse(data string) {
	c.writer.Header().Set("Content-Type", "text/plain")
	c.responseBody = []byte(data)
}

// sets response body as byte array
func (c *gskContext) RawResponse(raw []byte) {
	c.responseBody = raw
}

// get the logger
func (c *gskContext) Logger() *logrus.Logger {
	return c.logger
}

// Set cookie using http.Cookie
func (c *gskContext) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.writer, cookie)
}

// Get cookie using cookie name
func (c *gskContext) GetCookie(name string) (*http.Cookie, error) {
	return c.request.Cookie(name)
}

// Get the status code set for the response
func (c *gskContext) GetStatusCode() int {
	return c.responseStatus
}

// returns a copy of the context, now it's safe to use
func (c *gskContext) eject() gskContext {
	return *c
}
