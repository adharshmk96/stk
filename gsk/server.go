package gsk

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/adharshmk96/stk/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type HandlerFunc func(Context)

type ServerConfig struct {
	Port           string
	RequestLogging bool
	AllowedOrigins []string
	Logger         *logrus.Logger
}

type server struct {
	httpServer  *http.Server
	Router      *httprouter.Router
	Middlewares []Middleware
	Config      *ServerConfig
	Logger      *logrus.Logger
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

// NewServer creates a new server instance
func NewServer(config *ServerConfig) Server {
	if config.Logger == nil {
		config.Logger = logging.NewLogrusLogger()
	}

	newSTKServer := &server{
		httpServer: &http.Server{
			Addr: fmt.Sprintf(":%s", config.Port),
		},
		Router:      httprouter.New(),
		Middlewares: []Middleware{},
		Config:      config,
		Logger:      config.Logger,
	}

	return newSTKServer
}

// Start starts the server on the configured port
func (s *server) Start() {
	startingPort := NormalizePort(s.Config.Port)
	s.Logger.WithField("port", startingPort).Info("starting server")
	err := s.httpServer.ListenAndServe()
	if err != nil {
		s.Logger.WithError(err).Error("error starting server")
		panic(err)
	}
}

// Shuts down the server
func (s *server) Shutdown() error {
	s.Logger.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.httpServer.Shutdown(ctx)
}

// Use adds a middleware to the server
// usage example:
// server.Use(stk.RequestLogger())
func (s *server) Use(middleware Middleware) {
	s.Middlewares = append(s.Middlewares, middleware)
}

func (s *server) Get(path string, handler HandlerFunc) {
	s.Router.GET(path, wrapHandlerFunc(s.applyMiddleware(handler), s.Config))
}

func (s *server) Post(path string, handler HandlerFunc) {
	s.Router.POST(path, wrapHandlerFunc(s.applyMiddleware(handler), s.Config))
}

func (s *server) Put(path string, handler HandlerFunc) {
	s.Router.PUT(path, wrapHandlerFunc(s.applyMiddleware(handler), s.Config))
}

func (s *server) Delete(path string, handler HandlerFunc) {
	s.Router.DELETE(path, wrapHandlerFunc(s.applyMiddleware(handler), s.Config))
}

func (s *server) Patch(path string, handler HandlerFunc) {
	s.Router.PATCH(path, wrapHandlerFunc(s.applyMiddleware(handler), s.Config))
}

func (s *server) GetRouter() *httprouter.Router {
	return s.Router
}

// wrapHandlerFunc wraps the handler function with the httprouter.Handle
// this is done to pass the httprouter.Params to the handler
// and also to log the incoming request
func wrapHandlerFunc(handler HandlerFunc, config *ServerConfig) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		startTime := time.Now()

		if config.RequestLogging {
			config.Logger.WithFields(logrus.Fields{
				"method": r.Method,
				"url":    r.URL.String(),
			}).Info("incoming_request")
		}

		handlerContext := &gskContext{
			params:         p,
			request:        r,
			writer:         w,
			logger:         config.Logger,
			allowedOrigins: config.AllowedOrigins,
		}
		handler(handlerContext)

		if handlerContext.responseStatus != 0 {
			w.WriteHeader(handlerContext.responseStatus)
		} else {
			// Default to 200 OK
			w.WriteHeader(http.StatusOK)
		}

		if handlerContext.responseBody != nil {
			w.Write(handlerContext.responseBody)
		} else {
			w.Write([]byte(""))
		}

		if config.RequestLogging {
			timeTaken := time.Since(startTime).Milliseconds()
			config.Logger.WithFields(logrus.Fields{
				"method":    r.Method,
				"url":       r.URL.String(),
				"status":    handlerContext.responseStatus,
				"timeTaken": fmt.Sprintf("%d ms", timeTaken),
			}).Info("response_served")
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
