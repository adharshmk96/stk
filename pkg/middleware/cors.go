package middleware

import (
	"net/http"

	"github.com/adharshmk96/stk/gsk"
)

const (
	defaultCORSOrigin   = "same-origin"
	defaultAllowMethods = "POST, GET, OPTIONS, PUT, DELETE, PATCH"
	defaultAllowHeaders = "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"

	AccessControlAllowOrigin  = "Access-Control-Allow-Origin"
	AccessControlAllowMethods = "Access-Control-Allow-Methods"
	AccessControlAllowHeaders = "Access-Control-Allow-Headers"
)

func CORS(next gsk.HandlerFunc) gsk.HandlerFunc {
	return func(c gsk.Context) {
		allowedOrigins := getAllowedOrigins(c.GetAllowedOrigins())

		origin := c.GetRequest().Header.Get("Host")
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
			c.SetHeader("Content-Type", "text/plain")
			c.RawResponse([]byte("Forbidden"))
			return
		}

		// Set CORS headers
		headers := c.GetWriter().Header()
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
