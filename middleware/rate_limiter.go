package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/adharshmk96/stk/gsk"
)

type RateLimiter struct {
	requestsPerInterval int
	interval            time.Duration
	accessCounter       map[string]int
	mux                 *sync.Mutex
	Middleware          gsk.Middleware
}

func NewRateLimiter(requestsPerInterval int, interval time.Duration) *RateLimiter {

	rl := &RateLimiter{
		requestsPerInterval: requestsPerInterval,
		interval:            interval,
		accessCounter:       make(map[string]int),
		mux:                 &sync.Mutex{},
	}

	middleware := func(next gsk.HandlerFunc) gsk.HandlerFunc {
		return func(c gsk.Context) {
			clientIP := c.GetRequest().RemoteAddr
			rl.mux.Lock()
			defer rl.mux.Unlock()

			if cnt, ok := rl.accessCounter[clientIP]; ok {
				if cnt >= rl.requestsPerInterval {
					c.Status(http.StatusTooManyRequests).JSONResponse(gsk.Map{
						"error": "Too many requests. Please try again later.",
					})
					return
				}
				rl.accessCounter[clientIP]++
			} else {
				rl.accessCounter[clientIP] = 1
				go func(ip string) {
					time.Sleep(rl.interval)
					rl.mux.Lock()
					defer rl.mux.Unlock()
					delete(rl.accessCounter, ip)
				}(clientIP)
			}

			next(c)
		}
	}

	rl.Middleware = middleware

	return rl

}
