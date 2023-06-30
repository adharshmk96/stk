package gsk

// This implementation of middleware will enable middleware chaining
type Middleware func(HandlerFunc) HandlerFunc

// applyMiddleware applies all the middlewares to the handler
// in the reverse order, chaining the middlewares independently
func (s *Server) applyMiddleware(handler HandlerFunc) HandlerFunc {
	updatedHandler := handler
	for i := len(s.Middlewares) - 1; i >= 0; i-- {
		updatedHandler = s.Middlewares[i](updatedHandler)
	}
	return updatedHandler
}
