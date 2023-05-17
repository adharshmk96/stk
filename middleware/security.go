package middleware

import (
	"github.com/adharshmk96/stk"
)

func SecurityHeaders(next stk.HandlerFunc) stk.HandlerFunc {
	return func(c *stk.Context) {
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
