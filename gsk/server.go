package gsk

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/adharshmk96/stk/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type HandlerFunc func(Context)

type ServerConfig struct {
	Port   string
	Logger *logrus.Logger
	// Input
	BodySizeLimit int64
}

type server struct {
	httpServer  *http.Server
	router      *httprouter.Router
	middlewares []Middleware
	// configurations
	config *ServerConfig
}

type Server interface {
	// Start and Stop
	Start()
	Shutdown() error
	// Middleware
	Use(Middleware)

	// HTTP methods
	Get(string, HandlerFunc)
	Post(string, HandlerFunc)
	Put(string, HandlerFunc)
	Delete(string, HandlerFunc)
	Patch(string, HandlerFunc)

	// Other Server methods
	// Static(string, string)

	// Helpers
	Test(string, string, io.Reader, ...TestParams) (httptest.ResponseRecorder, error)

	// Internals
	Router() *httprouter.Router
	// to Apply middlewares
	applyMiddlewares(HandlerFunc) HandlerFunc
}

// Initialize the server configurations
// if no configurations are passed, default values are used
func initConfig(config ...*ServerConfig) *ServerConfig {
	var initConfig *ServerConfig
	if len(config) == 0 {
		initConfig = &ServerConfig{}
	} else {
		initConfig = config[0]
	}

	if initConfig.Port == "" {
		initConfig.Port = "8080"
	}

	if initConfig.Logger == nil {
		initConfig.Logger = logging.NewLogrusLogger()
	}

	if initConfig.BodySizeLimit == 0 {
		initConfig.BodySizeLimit = 1
	}

	return initConfig
}

// New creates a new server instance
// Configurations can be passed as a parameter and It's optional
// If no configurations are passed, default values are used
func New(userconfig ...*ServerConfig) Server {
	config := initConfig(userconfig...)

	startingPort := NormalizePort(config.Port)
	router := httprouter.New()

	newSTKServer := &server{
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
func (s *server) Start() {

	startingPort := NormalizePort(s.config.Port)
	s.config.Logger.WithField("port", startingPort).Info("starting server")
	err := s.httpServer.ListenAndServe()
	if err != nil {
		s.config.Logger.WithError(err).Error("error starting server")
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
func (s *server) Shutdown() error {
	s.config.Logger.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.httpServer.Shutdown(ctx)
}

// Use adds a middleware to the server
// usage example:
// server.Use(stk.RequestLogger())
func (s *server) Use(middleware Middleware) {
	s.middlewares = append(s.middlewares, middleware)

	// TODO: Fix this hack, this needs to be in a different place.
	s.router.GlobalOPTIONS = wrapHandlerFunc(
		s.applyMiddlewares(func(gc Context) {
			gc.Status(http.StatusNoContent)
		}),
		s,
	)
}

// Register handlers for the HTTP methods
// usage example:
// server.Get("/test", func(c stk.Context) { gc.Status(http.StatusOK).JSONResponse("OK") })
func (s *server) Get(path string, handler HandlerFunc) {
	s.router.HandlerFunc("GET", path, wrapHandlerFunc(s.applyMiddlewares(handler), s))
}

func (s *server) Post(path string, handler HandlerFunc) {
	s.router.HandlerFunc("POST", path, wrapHandlerFunc(s.applyMiddlewares(handler), s))
}

func (s *server) Put(path string, handler HandlerFunc) {
	s.router.HandlerFunc("PUT", path, wrapHandlerFunc(s.applyMiddlewares(handler), s))
}

func (s *server) Delete(path string, handler HandlerFunc) {
	s.router.HandlerFunc("DELETE", path, wrapHandlerFunc(s.applyMiddlewares(handler), s))
}

func (s *server) Patch(path string, handler HandlerFunc) {
	s.router.HandlerFunc("PATCH", path, wrapHandlerFunc(s.applyMiddlewares(handler), s))
}

type TestParams struct {
	Cookies []*http.Cookie
	Headers map[string]string
}

// Helper function to test the server
// Usage example:
// w, err := server.Test("GET", "/test", nil)
func (s *server) Test(method string, route string, body io.Reader, testParams ...TestParams) (httptest.ResponseRecorder, error) {

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

// Router returns the router instance
func (s *server) Router() *httprouter.Router {
	return s.router
}

// wrapHandlerFunc wraps the handler function with the httprouter.Handle
// this is done to pass the gsk context to the handler function
func wrapHandlerFunc(handler HandlerFunc, s *server) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		p := httprouter.ParamsFromContext(r.Context())

		handlerContext := &gskContext{
			params:        p,
			request:       r,
			writer:        w,
			logger:        s.config.Logger,
			bodySizeLimit: s.config.BodySizeLimit,
		}

		handler(handlerContext)

		gc := handlerContext.eject()

		if gc.responseStatus != 0 {
			w.WriteHeader(gc.responseStatus)
		} else {
			// Default to 200 OK
			w.WriteHeader(http.StatusOK)
		}

		if gc.responseBody != nil {
			w.Write(gc.responseBody)
		} else {
			w.Write([]byte(""))
		}

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
