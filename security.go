package stk

import (
	"net/http"
)

func SecurityHeaders(next HandlerFunc) HandlerFunc {
	return func(c *Context) {
		headers := map[string]string{
			"X-Content-Type-Options":            "nosniff",
			"X-Frame-Options":                   "SAMEORIGIN",
			"X-XSS-Protection":                  "1; mode=block",
			"Referrer-Policy":                   "strict-origin-when-cross-origin",
			"Content-Security-Policy":           "default-src 'self';",
			"X-Permitted-Cross-Domain-Policies": "master-only",
			"Strict-Transport-Security":         "max-age=31536000; includeSubDomains",
		}

		for key, value := range headers {
			c.Writer.Header().Set(key, value)
		}

		next(c)
	}
}

func CORS(next HandlerFunc) HandlerFunc {
	return func(c *Context) {
		// Set CORS headers
		headers := c.Writer.Header()
		// TODO: Make this configurable
		headers.Set("Access-Control-Allow-Origin", "*")
		headers.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		headers.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// For preflight requests, send only the headers and terminate the middleware chain
		if c.Request.Method == "OPTIONS" {
			c.Writer.WriteHeader(http.StatusNoContent)
			return
		}

		next(c)
	}
}
