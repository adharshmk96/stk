package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/adharshmk96/stk"
)

type RateLimiter struct {
	requestsPerInterval int
	interval            time.Duration
	AccessCounter       map[string]int
	mux                 *sync.Mutex
	Middleware          stk.Middleware
}

func NewRateLimiter(requestsPerInterval int, interval time.Duration) *RateLimiter {

	rl := &RateLimiter{
		requestsPerInterval: requestsPerInterval,
		interval:            interval,
		AccessCounter:       make(map[string]int),
		mux:                 &sync.Mutex{},
	}

	middleware := func(next stk.HandlerFunc) stk.HandlerFunc {
		return func(c *stk.Context) {
			clientIP := c.Request.RemoteAddr
			rl.mux.Lock()
			defer rl.mux.Unlock()

			if cnt, ok := rl.AccessCounter[clientIP]; ok {
				if cnt >= rl.requestsPerInterval {
					c.Status(http.StatusTooManyRequests).JSONResponse(stk.Map{
						"error": "Too many requests. Please try again later.",
					})
					return
				}
				rl.AccessCounter[clientIP]++
			} else {
				rl.AccessCounter[clientIP] = 1
				go func(ip string) {
					time.Sleep(rl.interval)
					rl.mux.Lock()
					defer rl.mux.Unlock()
					delete(rl.AccessCounter, ip)
				}(clientIP)
			}

			next(c)
		}
	}

	rl.Middleware = middleware

	return rl

}
