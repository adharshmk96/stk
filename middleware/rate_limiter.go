package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/adharshmk96/stk"
)

type rateLimiter struct {
	requestsPerInterval int
	interval            time.Duration
	accessCounter       map[string]int
	mux                 *sync.Mutex
}

func NewRateLimiter(requestsPerInterval int, interval time.Duration) stk.Middleware {

	rl := &rateLimiter{
		requestsPerInterval: requestsPerInterval,
		interval:            interval,
		accessCounter:       make(map[string]int),
		mux:                 &sync.Mutex{},
	}

	return func(next stk.HandlerFunc) stk.HandlerFunc {
		return func(c *stk.Context) {
			clientIP := c.Request.RemoteAddr
			rl.mux.Lock()
			defer rl.mux.Unlock()

			if cnt, ok := rl.accessCounter[clientIP]; ok {
				if cnt >= rl.requestsPerInterval {
					c.Status(http.StatusTooManyRequests)
					c.JSONResponse(stk.Map{
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
}
