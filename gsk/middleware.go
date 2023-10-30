package gsk

// This implementation of middleware will enable middleware chaining
type Middleware func(HandlerFunc) HandlerFunc

// applyMiddlewares applies all the middlewares to the handler
// in the reverse order, chaining the middlewares independently
func applyMiddlewares(middlewares []Middleware, handler HandlerFunc) HandlerFunc {
	finalHandler := handler
	for i := len(middlewares) - 1; i >= 0; i-- {
		finalHandler = middlewares[i](finalHandler)
	}
	return finalHandler
}
