package middleware

import (
	"net/http"
	"strings"

	"github.com/adharshmk96/stk/gsk"
)

var (
	defaultAllowMethods = []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "PATCH"}
	// "POST, GET, OPTIONS, PUT, DELETE, PATCH"
	defaultAllowHeaders = []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"}
	// "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"
	defaultAllowCredentials = "true"
)

const (
	defaultCORSOrigin = "same-origin"

	AccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	AccessControlAllowMethods     = "Access-Control-Allow-Methods"
	AccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	AccessControlAllowCredentials = "Access-Control-Allow-Credentials"
)

type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	AllowAll       bool
}

func CORS(config ...CORSConfig) gsk.Middleware {
	var corsConfig CORSConfig
	if len(config) > 0 {
		corsConfig = config[0]
	} else {
		corsConfig = CORSConfig{
			AllowedOrigins: []string{defaultCORSOrigin},
			AllowedMethods: defaultAllowMethods,
			AllowedHeaders: defaultAllowHeaders,
			AllowAll:       false,
		}
	}

	allowedMethods := strings.Join(defaultAllowMethods, ", ")
	if len(corsConfig.AllowedMethods) != 0 {
		allowedMethods = strings.Join(corsConfig.AllowedMethods, ", ")
	}
	allowedHeaders := strings.Join(defaultAllowHeaders, ", ")
	if len(corsConfig.AllowedHeaders) != 0 {
		allowedHeaders = strings.Join(corsConfig.AllowedHeaders, ", ")
	}

	return func(next gsk.HandlerFunc) gsk.HandlerFunc {
		return func(c gsk.Context) {
			allowedOrigins := getAllowedOrigins(corsConfig.AllowedOrigins)

			origin := c.Origin()
			// Check if the origin is in the allowedOrigins list
			isAllowed := false
			for _, allowedOrigin := range allowedOrigins {
				if allowedOrigin == "same-origin" || allowedOrigin == "*" || origin == allowedOrigin {
					isAllowed = true
					break
				}
			}

			if !corsConfig.AllowAll && !isAllowed {
				c.Status(http.StatusForbidden)
				c.SetHeader("Content-Type", "text/plain")
				c.RawResponse([]byte("Forbidden"))
				return
			}

			// Set CORS headers
			headers := c.GetWriter().Header()

			headers.Set(AccessControlAllowOrigin, origin)
			headers.Set(AccessControlAllowMethods, allowedMethods)
			headers.Set(AccessControlAllowHeaders, allowedHeaders)
			headers.Set(AccessControlAllowCredentials, defaultAllowCredentials)

			next(c)

		}
	}
}

func getAllowedOrigins(origins []string) []string {
	var allowedOrigins []string

	if len(origins) == 0 {
		allowedOrigins = []string{defaultCORSOrigin}
	} else {
		allowedOrigins = origins
	}
	return allowedOrigins
}
