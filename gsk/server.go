package gsk

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/adharshmk96/stk/logging"
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

type Server struct {
	Router      *httprouter.Router
	Middlewares []Middleware
	Config      *ServerConfig
	Logger      *logrus.Logger
}

// NewServer creates a new server instance
func NewServer(config *ServerConfig) *Server {
	if config.Logger == nil {
		config.Logger = logging.NewLogrusLogger()
	}

	newSTKServer := &Server{
		Router:      httprouter.New(),
		Middlewares: []Middleware{},
		Config:      config,
		Logger:      config.Logger,
	}

	return newSTKServer
}

// Start starts the server on the configured port
func (s *Server) Start() {
	startingPort := NormalizePort(s.Config.Port)
	s.Logger.WithField("port", startingPort).Info("starting server")
	err := http.ListenAndServe(startingPort, s.Router)
	if err != nil {
		s.Logger.WithError(err).Error("error starting server")
		panic(err)
	}
}

// Use adds a middleware to the server
// usage example:
// server.Use(stk.RequestLogger())
func (s *Server) Use(middleware Middleware) {
	s.Middlewares = append(s.Middlewares, middleware)
}

func (s *Server) Get(path string, handler HandlerFunc) {
	s.Router.GET(path, wrapHandlerFunc(s.applyMiddleware(handler), s.Config))
}

func (s *Server) Post(path string, handler HandlerFunc) {
	s.Router.POST(path, wrapHandlerFunc(s.applyMiddleware(handler), s.Config))
}

func (s *Server) Put(path string, handler HandlerFunc) {
	s.Router.PUT(path, wrapHandlerFunc(s.applyMiddleware(handler), s.Config))
}

func (s *Server) Delete(path string, handler HandlerFunc) {
	s.Router.DELETE(path, wrapHandlerFunc(s.applyMiddleware(handler), s.Config))
}

func (s *Server) Patch(path string, handler HandlerFunc) {
	s.Router.PATCH(path, wrapHandlerFunc(s.applyMiddleware(handler), s.Config))
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

		handlerContext := &context{
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
