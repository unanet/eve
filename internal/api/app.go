package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve/pkg/log"
	"gitlab.unanet.io/devops/eve/pkg/metrics"

	"gitlab.unanet.io/devops/eve/internal/config"

	"github.com/go-chi/chi"

	"gitlab.unanet.io/devops/eve/internal/controller/ping"
	"gitlab.unanet.io/devops/eve/pkg/artifactory"
	"gitlab.unanet.io/devops/eve/pkg/middleware"
	"gitlab.unanet.io/devops/eve/pkg/mux"
)

type App struct {
	r           chi.Router
	Controllers []mux.EveController
	Artifactory *artifactory.Client
	server      *http.Server
	mServer     *http.Server
	done        chan bool
	sigChannel  chan os.Signal
}

func NewApp() (*App, error) {
	client, err := artifactory.NewClient(config.Values.ArtifactoryConfig)
	if err != nil {
		return nil, err
	}
	router := chi.NewMux()
	return &App{
		r: router,
		Controllers: []mux.EveController{
			ping.New(),
		},
		Artifactory: client,
		server: &http.Server{
			ReadTimeout:  time.Duration(5) * time.Second,
			WriteTimeout: time.Duration(30) * time.Second,
			IdleTimeout:  time.Duration(90) * time.Second,
			Addr:         fmt.Sprintf(":%d", config.Values.Port),
			Handler:      router,
		},
		done:       make(chan bool),
		sigChannel: make(chan os.Signal, 1024),
	}, nil
}

// Handle SIGNALS
func (a *App) sigHandler() {
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

func (a *App) gracefulShutdown() {
	// Pause the Context for `ShutdownTimeoutSecs` config value
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(120)*time.Second)
	defer cancel()

	// Turn off keepalive
	a.server.SetKeepAlivesEnabled(false)

	// Attempt to shutdown cleanly
	if err := a.server.Shutdown(ctx); err != nil {
		panic("HTTP EVE-API Server Failed Graceful Shutdown")
	}

	if err := a.mServer.Shutdown(ctx); err != nil {
		panic("HTTP EVE Metrics Server Failed Graceful Shutdown")
	}
	close(a.done)
}

func (a *App) Start() {
	a.setup()
	a.mServer = metrics.StartMetricsServer(a.done)

	signal.Notify(a.sigChannel, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	go a.sigHandler()

	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Logger.Panic("Failed to Start Server", zap.Error(err))
	}

	<-a.done
	log.Logger.Info("Eve-API Shutdown")
}

func (a *App) setup() {
	middleware.SetupMiddleware(a.r, 60*time.Second)
	for _, c := range a.Controllers {
		c.Setup(a.r)
	}
}
