package middleware

import (
	"fmt"
	"net/http"

	"github.com/adharshmk96/stk"
)

const (
	defaultCORSOrigin   = "same-origin"
	defaultAllowMethods = "POST, GET, OPTIONS, PUT, DELETE, PATCH"
	defaultAllowHeaders = "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"

	AccessControlAllowOrigin  = "Access-Control-Allow-Origin"
	AccessControlAllowMethods = "Access-Control-Allow-Methods"
	AccessControlAllowHeaders = "Access-Control-Allow-Headers"
)

func CORS(next stk.HandlerFunc) stk.HandlerFunc {
	return func(c *stk.Context) {
		allowedOrigins := getAllowedOrigins(c)

		origin := c.Request.Header.Get("Host")
		// Check if the origin is in the allowedOrigins list
		isAllowed := false
		for _, allowedOrigin := range allowedOrigins {
			if allowedOrigin == "same-origin" || allowedOrigin == "*" || origin == allowedOrigin {
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

		fmt.Println("origin")
		fmt.Println(origin)

		// Set CORS headers
		headers := c.Writer.Header()
		// TODO: Make this configurable
		headers.Set(AccessControlAllowOrigin, origin)
		headers.Set(AccessControlAllowMethods, defaultAllowMethods)
		headers.Set(AccessControlAllowHeaders, defaultAllowHeaders)

		next(c)
	}
}

func getAllowedOrigins(c *stk.Context) []string {
	var allowedOrigins []string

	if c.AllowedOrigins == nil {
		allowedOrigins = []string{defaultCORSOrigin}
	} else {
		allowedOrigins = c.AllowedOrigins
	}
	return allowedOrigins
}
