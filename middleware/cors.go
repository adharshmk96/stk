package middleware

import (
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
		allowedOrigins := getAllowedOrigins(c.GetAllowedOrigins())

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
			c.RawResponse([]byte("Forbidden"))
			return
		}

		// Set CORS headers
		headers := c.Writer.Header()
		// TODO: Make this configurable
		headers.Set(AccessControlAllowOrigin, origin)
		headers.Set(AccessControlAllowMethods, defaultAllowMethods)
		headers.Set(AccessControlAllowHeaders, defaultAllowHeaders)

		next(c)
	}
}

func getAllowedOrigins(origins []string) []string {
	var allowedOrigins []string

	if origins == nil || len(origins) == 0 {
		allowedOrigins = []string{defaultCORSOrigin}
	} else {
		allowedOrigins = origins
	}
	return allowedOrigins
}
