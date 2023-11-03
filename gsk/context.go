package gsk

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

// TODO: experiment with a response object

type Context struct {
	// server
	Request *http.Request
	Writer  http.ResponseWriter

	// request
	params        Params
	bodySizeLimit int64

	// logging
	logger *slog.Logger

	// response
	responseStatus  int
	responseBody    []byte
	responseWritten bool
}

type Map map[string]interface{}

// Methods to handle request

// Param gets the params within the path mentioned as a wildcard
func (c *Context) Param(key string) string {
	return c.params.ByName(key)
}

// QueryParam gets the query parameters passed eg: /?name=value
func (c *Context) QueryParam(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Context) Path() string {
	return c.Request.URL.Path
}

func (c *Context) Origin() string {
	return c.Request.Header.Get("Origin")
}

// Get cookie using cookie name
func (c *Context) GetCookie(name string) (*http.Cookie, error) {
	return c.Request.Cookie(name)
}

func (c *Context) DecodeJSONBody(v interface{}) error {
	bodySizeLimit := int64(c.bodySizeLimit << 20) // 1 MB

	if c.Request.Body == nil {
		return ErrInvalidJSON
	}

	// Set a maximum limit for the request body size to avoid possible malicious requests
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, bodySizeLimit)

	// Manually check if the request body size exceeds the limit
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Writer.Header().Set("Content-Type", "application/json")
		http.Error(c.Writer, ErrBodyTooLarge.Error(), http.StatusRequestEntityTooLarge)
		return ErrBodyTooLarge
	}

	// Decode the JSON body into the provided interface
	err = json.Unmarshal(body, v)

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

// Methods related to response

// For formdata and multipart formdata, use c.Request.ParseMultipartForm ...
// sets response header key value
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Add(key, value)
}

// Set cookie using http.Cookie
func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Writer, cookie)
}

// Status sets the status code of the response
func (c *Context) Status(status int) *Context {
	c.responseStatus = status
	return c
}

// Redirect redirects the request to the provided URL
func (c *Context) Redirect(url string) {
	if c.responseStatus == 0 {
		c.responseStatus = http.StatusTemporaryRedirect
	}
	status := c.responseStatus
	c.responseWritten = true
	// Set the default redirect status code
	http.Redirect(c.Writer, c.Request, url, status)
}

// JSONResponse marshals the provided interface into JSON and writes it to the response writer
// If there is an error in marshalling the JSON, an internal server error is returned
func (c *Context) JSONResponse(data interface{}) {
	var err error
	// Set the content type to JSON
	c.Writer.Header().Set("Content-Type", "application/json")
	c.responseBody, err = json.Marshal(data)

	// Check if there is an error in marshalling the JSON (internal server error)
	if err != nil {
		c.responseStatus = http.StatusInternalServerError
		c.responseBody = []byte(ErrInternalServer.Error())
	}
}

// TemplateResponse renders the provided template with the provided data
// and writes it to the response writer with content type text/html
func (c *Context) TemplateResponse(template *Tpl) {
	var err error
	c.Writer.Header().Set("Content-Type", "text/html")
	c.responseBody, err = template.Render(DEFAULT_TEMPLATE_VARIABLES)
	if err != nil {
		c.responseStatus = http.StatusInternalServerError
		c.responseBody = []byte(ErrInternalServer.Error())
	}

}

// StringResponse writes the provided string to the response writer
func (c *Context) StringResponse(data string) {
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.responseBody = []byte(data)
}

// sets response body as byte array
func (c *Context) RawResponse(raw []byte) {
	c.responseBody = raw
}

// TODO: support for template rendering, file response, stream response

// Methods to get context values

// get the logger
func (c *Context) Logger() *slog.Logger {
	return c.logger
}

// Get the status code set for the response
func (c *Context) GetStatusCode() int {
	return c.responseStatus
}

// returns a copy of the context, now it's safe to use
func (c *Context) eject() Context {
	return *c
}
