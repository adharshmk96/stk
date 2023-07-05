package gsk

// This implementation of middleware will enable middleware chaining
type Middleware func(HandlerFunc) HandlerFunc

// applyMiddlewares applies all the middlewares to the handler
// in the reverse order, chaining the middlewares independently
func (s *server) applyMiddlewares(handler HandlerFunc) HandlerFunc {
	updatedHandler := handler
	for i := len(s.middlewares) - 1; i >= 0; i-- {
		updatedHandler = s.middlewares[i](updatedHandler)
	}
	return updatedHandler
}
