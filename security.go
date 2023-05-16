package stk

import (
	"net/http"
	"strings"
)

const (
	defaultCORSOrigin = "same-origin"
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
		allowedOrigins := getAllowedOrigins(c)

		origin := c.Request.Header.Get("Origin")
		// Check if the origin is in the allowedOrigins list
		isAllowed := false
		for _, allowedOrigin := range allowedOrigins {
			originWithoutPort := strings.Split(origin, ":")[0]
			if allowedOrigin == "same-origin" && origin == c.Request.Host {
				isAllowed = true
				origin = allowedOrigin
				break
			}
			if allowedOrigin == "*" {
				isAllowed = true
				break
			}
			if origin == allowedOrigin || origin == originWithoutPort {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.Status(http.StatusForbidden)
			c.Writer.Header().Set("Content-Type", "text/plain")
			c.Writer.Write([]byte("Forbidden"))
			return
		}

		// Set CORS headers
		headers := c.Writer.Header()
		// TODO: Make this configurable
		headers.Set("Access-Control-Allow-Origin", origin)
		headers.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		headers.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		next(c)
	}
}

func getAllowedOrigins(c *Context) []string {
	var allowedOrigins []string

	if c.AllowedOrigins == nil {
		allowedOrigins = []string{defaultCORSOrigin}
	} else {
		allowedOrigins = c.AllowedOrigins
	}
	return allowedOrigins
}
