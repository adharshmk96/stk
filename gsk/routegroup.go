package gsk

type routeGroup struct {
	server      *server
	pathPrefix  string
	middlewares []Middleware
}

type RouteGroup interface {
	Use(middleware Middleware)
	applyMiddlewares(handler HandlerFunc) HandlerFunc

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
	rg.server.Get(rg.pathPrefix+path, rg.applyMiddlewares(handler))
}

func (rg *routeGroup) Post(path string, handler HandlerFunc) {
	rg.server.Post(rg.pathPrefix+path, rg.applyMiddlewares(handler))
}

func (rg *routeGroup) Put(path string, handler HandlerFunc) {
	rg.server.Put(rg.pathPrefix+path, rg.applyMiddlewares(handler))
}

func (rg *routeGroup) Delete(path string, handler HandlerFunc) {
	rg.server.Delete(rg.pathPrefix+path, rg.applyMiddlewares(handler))
}

func (rg *routeGroup) Patch(path string, handler HandlerFunc) {
	rg.server.Patch(rg.pathPrefix+path, rg.applyMiddlewares(handler))
}

func (rg *routeGroup) Handle(method string, path string, handler HandlerFunc) {
	rg.server.Handle(method, rg.pathPrefix+path, rg.applyMiddlewares(handler))
}

// Here we apply the RouteGroup middlewares to a handler
func (rg *routeGroup) applyMiddlewares(handler HandlerFunc) HandlerFunc {
	finalHandler := handler
	for i := len(rg.middlewares) - 1; i >= 0; i-- {
		finalHandler = rg.middlewares[i](finalHandler)
	}
	return finalHandler
}

// Extend the server interface and server struct
func (s *server) RouteGroup(path string) RouteGroup {
	return &routeGroup{
		server:      s,
		pathPrefix:  path,
		middlewares: s.middlewares,
	}
}
