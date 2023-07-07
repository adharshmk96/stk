package gsk

type routeGroup struct {
	server      *server
	pathPrefix  string
	middlewares []Middleware
}

type RouteGroup interface {
	Use(middleware Middleware)

	Get(path string, handler HandlerFunc)
	Post(path string, handler HandlerFunc)
	Put(path string, handler HandlerFunc)
	Delete(path string, handler HandlerFunc)
	Patch(path string, handler HandlerFunc)
	Handle(method string, path string, handler HandlerFunc)

	RouteGroup(path string) RouteGroup
}

func (rg *routeGroup) Use(middleware Middleware) {
	rg.middlewares = append(rg.middlewares, middleware)
}

func (rg *routeGroup) RouteGroup(path string) RouteGroup {
	return &routeGroup{
		server:      rg.server,
		pathPrefix:  rg.pathPrefix + path,
		middlewares: rg.middlewares,
	}
}

func (rg *routeGroup) Get(path string, handler HandlerFunc) {
	rg.server.Get(rg.pathPrefix+path, applyMiddlewares(rg.middlewares, handler))
}

func (rg *routeGroup) Post(path string, handler HandlerFunc) {
	rg.server.Post(rg.pathPrefix+path, applyMiddlewares(rg.middlewares, handler))
}

func (rg *routeGroup) Put(path string, handler HandlerFunc) {
	rg.server.Put(rg.pathPrefix+path, applyMiddlewares(rg.middlewares, handler))
}

func (rg *routeGroup) Delete(path string, handler HandlerFunc) {
	rg.server.Delete(rg.pathPrefix+path, applyMiddlewares(rg.middlewares, handler))
}

func (rg *routeGroup) Patch(path string, handler HandlerFunc) {
	rg.server.Patch(rg.pathPrefix+path, applyMiddlewares(rg.middlewares, handler))
}

func (rg *routeGroup) Handle(method string, path string, handler HandlerFunc) {
	rg.server.Handle(method, rg.pathPrefix+path, applyMiddlewares(rg.middlewares, handler))
}
