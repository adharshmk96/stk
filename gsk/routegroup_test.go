package gsk_test

import (
	"net/http"
	"testing"

	"github.com/adharshmk96/stk/gsk"
	"github.com/stretchr/testify/assert"
)

func TestRouteGroup(t *testing.T) {
	t.Run("route group registers routes with the correct path prefix", func(t *testing.T) {
		server := gsk.New()

		handler := func(c gsk.Context) {
			c.Status(http.StatusTeapot).StringResponse(c.GetPath())
		}

		server.Get("/users", handler)

		rg := server.RouteGroup("/api")
		rg.Get("/users", handler)

		r1, err := server.Test("GET", "/api/users", nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusTeapot, r1.Code)
		assert.Equal(t, "/api/users", r1.Body.String())

		r2, err := server.Test("GET", "/users", nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusTeapot, r2.Code)
		assert.Equal(t, "/users", r2.Body.String())
	})

	t.Run("route group registers routes with the correct middlewares", func(t *testing.T) {
		server := gsk.New()

		handler := func(c gsk.Context) {
			c.Status(http.StatusTeapot).StringResponse(c.GetPath())
		}

		globalMiddleware := func(next gsk.HandlerFunc) gsk.HandlerFunc {
			return func(c gsk.Context) {
				c.SetHeader("X-Global", "global")
				next(c)
			}
		}

		authMiddleware := func(next gsk.HandlerFunc) gsk.HandlerFunc {
			return func(c gsk.Context) {
				c.SetHeader("X-Auth", "auth")
				next(c)
			}
		}

		adminMiddleware := func(next gsk.HandlerFunc) gsk.HandlerFunc {
			return func(c gsk.Context) {
				c.SetHeader("X-Admin", "admin")
				next(c)
			}
		}

		server.Use(globalMiddleware)
		server.Get("/users", handler)

		rg := server.RouteGroup("/auth")
		rg.Use(authMiddleware)
		rg.Get("/users", handler)

		rg2 := rg.RouteGroup("/admin")
		rg2.Use(adminMiddleware)
		rg2.Get("/users", handler)

		rg.Get("/me", handler)

		rgp := server.RouteGroup("/public")
		rgp.Get("/users", handler)

		r1, err := server.Test("GET", "/users", nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusTeapot, r1.Code)
		assert.Equal(t, "/users", r1.Body.String())
		assert.Equal(t, "global", r1.Header().Get("X-Global"))

		r2, err := server.Test("GET", "/auth/users", nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusTeapot, r2.Code)
		assert.Equal(t, "/auth/users", r2.Body.String())
		assert.Equal(t, "global", r2.Header().Get("X-Global"))
		assert.Equal(t, "auth", r2.Header().Get("X-Auth"))

		r3, err := server.Test("GET", "/auth/admin/users", nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusTeapot, r3.Code)
		assert.Equal(t, "/auth/admin/users", r3.Body.String())
		assert.Equal(t, "global", r3.Header().Get("X-Global"))
		assert.Equal(t, "auth", r3.Header().Get("X-Auth"))
		assert.Equal(t, "admin", r3.Header().Get("X-Admin"))

		r4, err := server.Test("GET", "/auth/me", nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusTeapot, r4.Code)
		assert.Equal(t, "/auth/me", r4.Body.String())
		assert.Equal(t, "global", r4.Header().Get("X-Global"))
		assert.Equal(t, "auth", r4.Header().Get("X-Auth"))
		assert.Equal(t, "", r4.Header().Get("X-Admin"))

		r5, err := server.Test("GET", "/public/users", nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusTeapot, r5.Code)
		assert.Equal(t, "/public/users", r5.Body.String())
		assert.Equal(t, "global", r5.Header().Get("X-Global"))
		assert.Equal(t, "", r5.Header().Get("X-Auth"))
		assert.Equal(t, "", r5.Header().Get("X-Admin"))

	})
}
