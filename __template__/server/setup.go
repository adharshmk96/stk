package server

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/adharshmk96/stk/gsk"
	"github.com/adharshmk96/stk/pkg/middleware"
	"github.com/adharshmk96/stktemplate/server/infra"
	svrmw "github.com/adharshmk96/stktemplate/server/middleware"
	"github.com/adharshmk96/stktemplate/server/routing"
)

func StartHttpServer(port string) (*gsk.Server, chan bool) {

	logger := infra.GetLogger()

	serverConfig := &gsk.ServerConfig{
		Port:   port,
		Logger: logger,
	}

	server := gsk.New(serverConfig)

	rateLimiter := svrmw.RateLimiter()
	server.Use(rateLimiter)
	server.Use(middleware.RequestLogger)
	server.Use(middleware.CORS(middleware.CORSConfig{
		AllowAll: true,
	}))

	infra.LoadDefaultConfig()

	routing.SetupTemplateRoutes(server)
	routing.SetupApiRoutes(server)

	server.Start()

	// graceful shutdown
	done := make(chan bool)

	// A go routine that listens for os signals
	// it will block until it receives a signal
	// once it receives a signal, it will shut down close the done channel
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		if err := server.Shutdown(); err != nil {
			logger.Error(err.Error())
		}

		close(done)
	}()

	return server, done
}
