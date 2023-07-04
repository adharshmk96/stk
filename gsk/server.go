package gsk

import (
	"context"
	"net/http"
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

	// GetRouter
	GetRouter() *httprouter.Router
}

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

// Shuts down the server
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
}

func (s *server) Get(path string, handler HandlerFunc) {
	s.router.GET(path, wrapHandlerFunc(s.applyMiddleware(handler), s))
}

func (s *server) Post(path string, handler HandlerFunc) {
	s.router.POST(path, wrapHandlerFunc(s.applyMiddleware(handler), s))
}

func (s *server) Put(path string, handler HandlerFunc) {
	s.router.PUT(path, wrapHandlerFunc(s.applyMiddleware(handler), s))
}

func (s *server) Delete(path string, handler HandlerFunc) {
	s.router.DELETE(path, wrapHandlerFunc(s.applyMiddleware(handler), s))
}

func (s *server) Patch(path string, handler HandlerFunc) {
	s.router.PATCH(path, wrapHandlerFunc(s.applyMiddleware(handler), s))
}

func (s *server) GetRouter() *httprouter.Router {
	return s.router
}

// wrapHandlerFunc wraps the handler function with the httprouter.Handle
// this is done to pass the httprouter.Params to the handler
// and also to log the incoming request
func wrapHandlerFunc(handler HandlerFunc, s *server) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		handlerContext := &gskContext{
			params:        p,
			request:       r,
			writer:        w,
			logger:        s.config.Logger,
			bodySizeLimit: s.config.BodySizeLimit,
		}

		handler(handlerContext)

		ctx := handlerContext.eject()

		if ctx.responseStatus != 0 {
			w.WriteHeader(ctx.responseStatus)
		} else {
			// Default to 200 OK
			w.WriteHeader(http.StatusOK)
		}

		if ctx.responseBody != nil {
			w.Write(ctx.responseBody)
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
