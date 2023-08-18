package middleware

import (
	"fmt"
	"time"

	"github.com/adharshmk96/stk/gsk"
)

func RequestLogger(next gsk.HandlerFunc) gsk.HandlerFunc {
	return func(c *gsk.Context) {
		startTime := time.Now()
		c.Logger().Info(
			"incoming_request",
			"method", c.Request.Method,
			"url", c.Request.URL.String(),
		)

		next(c)

		timeTaken := time.Since(startTime).Milliseconds()
		c.Logger().Info(
			"request_served",
			"method", c.Request.Method,
			"url", c.Request.URL.String(),
			"status", c.GetStatusCode(),
			"timeTaken", fmt.Sprintf("%d ms", timeTaken),
		)
	}
}
