package gsk

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

type HandlerFunc func(*Context)

type ServerConfig struct {
	Port   string
	Logger *slog.Logger
	// Input
	BodySizeLimit int64

	// Static
	StaticPath string
	StaticDir  string
}

type Server struct {
	httpServer  *http.Server
	router      Router
	middlewares []Middleware
	// configurations
	config *ServerConfig
}

// New creates a new server instance
// Configurations can be passed as a parameter and It's optional
// If no configurations are passed, default values are used
// Returning server struct pointer because its more performant
// interface{} is slower than concrete types, because it slows down garbage collector
func New(userconfig ...*ServerConfig) *Server {
	config := initConfig(userconfig...)

	startingPort := NormalizePort(config.Port)
	router := newGskRouter()

	// Serve static files
	router.ServeFiles(config.StaticPath+"/*filepath", http.Dir(config.StaticDir))

	newSTKServer := &Server{
		httpServer: &http.Server{
			Addr:    startingPort,
			Handler: router,
		},
		router:      router,
		middlewares: []Middleware{},
		config:      config,
	}

	return newSTKServer
}

// Start starts the server on the configured port
func (s *Server) Start() {
	startingPort := NormalizePort(s.config.Port)
	s.config.Logger.Info("starting server", "port", startingPort)
	err := s.httpServer.ListenAndServe()
	if err != nil {
		s.config.Logger.Error("error starting server", "error", err)
		panic(err)
	}
}

// Shuts down the server, use for graceful shutdown
// Eg Usage:
/*
// indicate that the server is shutting down
done := make(chan bool)

// A go routine that listens for os signals
// it will block until it receives a signal
// once it receives a signal, it will shutdown close the done channel
go func() {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint

	if err := server.Shutdown(); err != nil {
		logger.Error(err)
	}

	close(done)
}()

return server, done
*/
func (s *Server) Shutdown() error {
	s.config.Logger.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.httpServer.Shutdown(ctx)
}

// Use adds a middleware to the server
// usage example:
// server.Use(stk.RequestLogger())
// NOTE: Middlewares will be applied when the route is registered
// SO Make sure to register the routes after adding the middlewares
func (s *Server) Use(mw Middleware) {
	s.middlewares = append(s.middlewares, mw)

	// Add preflight handler if CORS middleware is used
	// this is a hack to make sure that the preflight handler works with CORS Middleware
	router := s.router.Router()
	router.GlobalOPTIONS = wrapHandlerFunc(s, applyMiddlewares(s.middlewares, preFlightHandler))
}

// Register handlers for the HTTP methods
// usage example:
// server.Get("/test", func(c stk.Context) { gc.Status(http.StatusOK).JSONResponse("OK") })
func (s *Server) Get(path string, handler HandlerFunc) {
	s.Handle(http.MethodGet, path, handler)
}

func (s *Server) Post(path string, handler HandlerFunc) {
	s.Handle(http.MethodPost, path, handler)
}

func (s *Server) Put(path string, handler HandlerFunc) {
	s.Handle(http.MethodPut, path, handler)
}

func (s *Server) Delete(path string, handler HandlerFunc) {
	s.Handle(http.MethodDelete, path, handler)
}

func (s *Server) Patch(path string, handler HandlerFunc) {
	s.Handle(http.MethodPatch, path, handler)
}

func preFlightHandler(gc *Context) {
	gc.Status(http.StatusNoContent)
}

func (s *Server) Handle(method string, path string, handler HandlerFunc) {
	if path != "/" {
		path = strings.TrimSuffix(path, "/")
	}
	s.router.HandlerFunc(method, path, wrapHandlerFunc(s, handler))
}

// RouteGroup returns a new RouteGroup instance
// RouteGroup is used to register routes with the same path prefix
// It will also ensure that the middlewares are applied to the routes exclusively
// usage example:
// rg := server.RouteGroup("/api")
// rg.Get("/users", func(c stk.Context) { gc.Status(http.StatusOK).JSONResponse("OK") })
func (s *Server) RouteGroup(path string) *RouteGroup {
	// Trim suffix /
	path = strings.TrimSuffix(path, "/")

	return &RouteGroup{
		server:     s,
		pathPrefix: path,
	}
}

type TestParams struct {
	Cookies []*http.Cookie
	Headers map[string]string
}

// Helper function to test the server
// Usage example:
// w, err := server.Test("GET", "/test", nil)
func (s *Server) Test(method string, route string, body io.Reader, testParams ...TestParams) (httptest.ResponseRecorder, error) {

	req, err := http.NewRequest(method, route, body)

	if len(testParams) > 0 {
		for _, cookie := range testParams[0].Cookies {
			req.AddCookie(cookie)
		}

		for key, value := range testParams[0].Headers {
			req.Header.Set(key, value)
		}
	}

	if len(testParams) > 0 && len(testParams[0].Cookies) > 0 {
		for _, cookie := range testParams[0].Cookies {
			req.AddCookie(cookie)
		}
	}

	if len(testParams) > 0 && len(testParams[0].Headers) > 0 {
		for key, value := range testParams[0].Headers {
			req.Header.Set(key, value)
		}
	}

	w := httptest.NewRecorder()
	if err != nil {
		return *w, err
	}
	s.router.ServeHTTP(w, req)
	return *w, nil
}

// wrapHandlerFunc wraps the handler function with the router.Handle
// this is done to pass the gsk context to the handler function
func wrapHandlerFunc(s *Server, handler HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		p := s.router.ParamsFromContext(r.Context())

		handlerContext := &Context{
			params:        p,
			Request:       r,
			Writer:        w,
			logger:        s.config.Logger,
			bodySizeLimit: s.config.BodySizeLimit,
		}

		finalHandler := applyMiddlewares(s.middlewares, handler)
		finalHandler(handlerContext)

	}
}

func NormalizePort(val string) string {
	var result string
	if strings.ContainsAny(val, ".") {
		result = val
	} else if strings.HasPrefix(val, ":") {
		result = "0.0.0.0" + val
	} else {
		result = "0.0.0.0:" + val
	}
	return result
}
