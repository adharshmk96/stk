package gsk

import "strings"

type RouteGroup struct {
	server      *Server
	pathPrefix  string
	middlewares []Middleware
}

func (rg *RouteGroup) Use(middleware Middleware) {
	rg.middlewares = append(rg.middlewares, middleware)
}

func (rg *RouteGroup) RouteGroup(path string) *RouteGroup {
	path = strings.TrimSuffix(path, "/")
	return &RouteGroup{
		server:      rg.server,
		pathPrefix:  rg.pathPrefix + path,
		middlewares: rg.middlewares,
	}
}

func (rg *RouteGroup) Get(path string, handler HandlerFunc) {
	rg.server.Get(rg.pathPrefix+path, applyMiddlewares(rg.middlewares, handler))
}

func (rg *RouteGroup) Post(path string, handler HandlerFunc) {
	rg.server.Post(rg.pathPrefix+path, applyMiddlewares(rg.middlewares, handler))
}

func (rg *RouteGroup) Put(path string, handler HandlerFunc) {
	rg.server.Put(rg.pathPrefix+path, applyMiddlewares(rg.middlewares, handler))
}

func (rg *RouteGroup) Delete(path string, handler HandlerFunc) {
	rg.server.Delete(rg.pathPrefix+path, applyMiddlewares(rg.middlewares, handler))
}

func (rg *RouteGroup) Patch(path string, handler HandlerFunc) {
	rg.server.Patch(rg.pathPrefix+path, applyMiddlewares(rg.middlewares, handler))
}

func (rg *RouteGroup) Handle(method string, path string, handler HandlerFunc) {
	rg.server.Handle(method, rg.pathPrefix+path, applyMiddlewares(rg.middlewares, handler))
}
