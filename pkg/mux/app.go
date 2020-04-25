package mux

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve/pkg/log"
	"gitlab.unanet.io/devops/eve/pkg/metrics"
	"gitlab.unanet.io/devops/eve/pkg/middleware"
)

var (
	Version = "unset"
)

type Api struct {
	r           chi.Router
	controllers []EveController
	server      *http.Server
	mServer     *http.Server
	done        chan bool
	sigChannel  chan os.Signal
	config      *Config
}

func NewApi(controllers []EveController, c Config) (*Api, error) {
	router := chi.NewMux()
	return &Api{
		r:           router,
		config:      &c,
		controllers: controllers,
		server: &http.Server{
			ReadTimeout:  time.Duration(5) * time.Second,
			WriteTimeout: time.Duration(30) * time.Second,
			IdleTimeout:  time.Duration(90) * time.Second,
			Addr:         fmt.Sprintf(":%d", c.Port),
			Handler:      router,
		},
		done:       make(chan bool),
		sigChannel: make(chan os.Signal, 1024),
	}, nil
}

// Handle SIGNALS
func (a *Api) sigHandler() {
	for {
		sig := <-a.sigChannel
		switch sig {
		case syscall.SIGHUP:
			log.Logger.Warn("SIGHUP hit, Nothing supports this currently")
		case os.Interrupt, syscall.SIGTERM, syscall.SIGINT:
			log.Logger.Info("Caught Shutdown Signal", zap.String("signal", sig.String()))
			a.gracefulShutdown()
		}
	}
}

func (a *Api) gracefulShutdown() {
	// Pause the Context for `ShutdownTimeoutSecs` config value
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(120)*time.Second)
	defer cancel()

	// Turn off keepalive
	a.server.SetKeepAlivesEnabled(false)

	// Attempt to shutdown cleanly
	if err := a.server.Shutdown(ctx); err != nil {
		panic("HTTP Server Failed Graceful Shutdown")
	}

	if err := a.mServer.Shutdown(ctx); err != nil {
		panic("HTTP Metrics Server Failed Graceful Shutdown")
	}
	close(a.done)
}

// Start starts the Mux Service Listeners (API/Metrics)
func (a *Api) Start() {
	a.setup()
	a.mServer = metrics.StartMetricsServer(a.done, a.config.MetricsPort)
	log.Logger.Info(fmt.Sprintf("start %v metrics listener", a.config.ServiceName), zap.Int("port", a.config.MetricsPort))

	signal.Notify(a.sigChannel, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	go a.sigHandler()

	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Logger.Panic("Failed to Start Server", zap.Error(err))
	}

	log.Logger.Info(fmt.Sprintf("start %v api listener", a.config.ServiceName), zap.Int("port", a.config.Port))
	<-a.done
	log.Logger.Info("Service Shutdown", zap.String("service_name", a.config.ServiceName))
}

func (a *Api) setup() {
	middleware.SetupMiddleware(a.r, 60*time.Second)
	for _, c := range a.controllers {
		c.Setup(a.r)
	}
}
