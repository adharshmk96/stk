package middleware

import (
	"fmt"
	"time"

	"github.com/adharshmk96/stk/gsk"
	"github.com/sirupsen/logrus"
)

func RequestLogger(next gsk.HandlerFunc) gsk.HandlerFunc {
	return func(c gsk.Context) {
		startTime := time.Now()
		c.Logger().WithFields(logrus.Fields{
			"method": c.GetRequest().Method,
			"url":    c.GetRequest().URL.String(),
		}).Info("incoming_request")
		next(c)
		timeTaken := time.Since(startTime).Milliseconds()
		c.Logger().WithFields(logrus.Fields{
			"method":    c.GetRequest().Method,
			"url":       c.GetRequest().URL.String(),
			"status":    c.GetStatusCode(),
			"timeTaken": fmt.Sprintf("%d ms", timeTaken),
		}).Info("response_served")
	}
}
