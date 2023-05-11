package stk

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type Context struct {
	Request *http.Request
	Writer  http.ResponseWriter

	Params httprouter.Params

	Logger *zap.Logger

	ResponseStatus int
}

type Map map[string]interface{}

// Status sets the status code of the response
func (c *Context) Status(status int) *Context {
	c.ResponseStatus = status
	return c
}

// JSONResponse marshals the provided interface into JSON and writes it to the response writer
// If there is an error in marshalling the JSON, an internal server error is returned
func (c *Context) JSONResponse(data interface{}) {
	response, err := json.Marshal(data)
	// Set the content type to JSON
	c.Writer.Header().Set("Content-Type", "application/json")

	// Check if there is an error in marshalling the JSON (internal server error)
	if err != nil {
		c.ResponseStatus = http.StatusInternalServerError
		http.Error(c.Writer, ErrInternalServer.Error(), http.StatusInternalServerError)
		return
	}

	c.Writer.Write(response)
}

func (c *Context) GetParam(key string) string {
	return c.Params.ByName(key)
}

func (c *Context) GetQueryParam(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Context) DecodeJSONBody(v interface{}) error {
	bodySizeLimit := int64(1 << 20) // 1 MB

	// Set a maximum limit for the request body size to avoid possible malicious requests
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, bodySizeLimit)

	// Manually check if the request body size exceeds the limit
	if c.Request.ContentLength > bodySizeLimit {
		c.Writer.Header().Set("Content-Type", "application/json")
		http.Error(c.Writer, ErrBodyTooLarge.Error(), http.StatusRequestEntityTooLarge)
		return ErrBodyTooLarge
	}

	// Decode the JSON body into the provided interface
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(v)

	defer c.Request.Body.Close()

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
