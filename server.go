package stk

import (
	"net/http"

	"github.com/adharshmk96/stk/logging"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type HandlerFunc func(*Context)

type ServerConfig struct {
	Port           string
	RequestLogging bool
	CORS           bool
	Logger         *zap.Logger
}

type Server struct {
	Router      *httprouter.Router
	Middlewares []Middleware
	Config      *ServerConfig
	Logger      *zap.Logger
}

// NewServer creates a new server instance
func NewServer(config *ServerConfig) *Server {
	if config.Logger == nil {
		config.Logger = logging.NewZapLogger()
	}

	newSTKServer := &Server{
		Router:      httprouter.New(),
		Middlewares: []Middleware{},
		Config:      config,
		Logger:      config.Logger,
	}

	newSTKServer.Use(SecurityHeaders)
	newSTKServer.Use(CORS)

	return newSTKServer
}

// Start starts the server on the configured port
func (s *Server) Start() {
	s.Logger.Info("starting server", zap.String("port", s.Config.Port))
	err := http.ListenAndServe(":"+s.Config.Port, s.Router)
	if err != nil {
		s.Logger.Panic("error starting server", zap.Error(err))
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

	logger := logging.NewZapLogger()

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		if config.RequestLogging {
			logger.Info("incoming request",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
			)
		}

		handlerContext := &Context{
			Params:  p,
			Request: r,
			Writer:  w,
			Logger:  logger,
		}
		handler(handlerContext)

		if handlerContext.ResponseStatus != 0 {
			w.WriteHeader(handlerContext.ResponseStatus)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}
