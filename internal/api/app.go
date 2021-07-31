package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/unanet/eve/internal/config"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"go.uber.org/zap"

	"github.com/unanet/go/pkg/log"
	"github.com/unanet/go/pkg/metrics"
	"github.com/unanet/go/pkg/middleware"
	"github.com/unanet/go/pkg/version"
)

var (
	Branch     = version.Branch
	SHA        = version.SHA
	ShortSHA   = version.ShortSHA
	Author     = version.Author
	BuildHost  = version.BuildHost
	Version    = version.Version
	Date       = version.Date
	Prerelease = version.Prerelease
)

type Api struct {
	r           chi.Router
	controllers []Controller
	server      *http.Server
	mServer     *http.Server
	//idSvc       *identity.Service
	//enforcer    *casbin.Enforcer
	done       chan bool
	sigChannel chan os.Signal
	config     *config.Config
	onShutdown []func()
	adminToken string
}

func NewApi(
	controllers []Controller,
	c config.Config,
) (*Api, error) {
	router := chi.NewMux()

	return &Api{
		adminToken: c.AdminToken,
		//idSvc:       svc,
		//enforcer:    enforcer,
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(300)*time.Second)
	defer cancel()

	// Turn off keepalive
	a.server.SetKeepAlivesEnabled(false)
	a.mServer.SetKeepAlivesEnabled(false)

	// Attempt to shutdown cleanly
	for _, x := range a.onShutdown {
		x()
	}
	if err := a.mServer.Shutdown(ctx); err != nil {
		panic("HTTP Metrics Server Failed Graceful Shutdown")
	}
	if err := a.server.Shutdown(ctx); err != nil {
		panic("HTTP API Server Failed Graceful Shutdown")
	}

	// not much to do here
	_ = log.Logger.Sync()

	close(a.done)
}

// Start starts the Mux Service Listeners (API/Metrics)
func (a *Api) Start(onShutdown ...func()) {
	a.setup()
	a.onShutdown = onShutdown
	a.mServer = metrics.StartMetricsServer(a.config.MetricsPort)

	signal.Notify(a.sigChannel, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	go a.sigHandler()
	log.Logger.Info("API Listener", zap.Int("port", a.config.Port))
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Logger.Panic("Failed to Start Server", zap.Error(err))
	}

	<-a.done
	log.Logger.Info("Service Shutdown")
}

func (a *Api) setup() {
	middleware.SetupMiddleware(a.r, 60*time.Second)
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	a.r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD", "CONNECT", "TRACE", "PATCH"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	authenticated := a.r.Group(nil)
	authenticated.Use(a.authenticationMiddleware())

	for _, c := range a.controllers {
		c.Setup(&Routers{
			Auth:      authenticated,
			Anonymous: a.r,
		})
	}
}

// TODO: Refactor into pkg
// This code is duped in cloud-api
func (a *Api) authenticationMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Admin token, you shall PASS!!!
			//if jwtauth.TokenFromHeader(r) == a.adminToken {
			//	next.ServeHTTP(w, r.WithContext(ctx))
			//	return
			//}
			//
			//claims, err := a.idSvc.TokenVerification(r)
			//if err != nil {
			//	middleware.Log(ctx).Debug("failed token verification", zap.Error(err))
			//	render.Respond(w, r, err)
			//	return
			//}
			//
			//middleware.Log(ctx).Debug("incoming auth claims", zap.Any("claims", claims))
			////role := extractRole(ctx, claims)
			//role := "admin"
			//middleware.Log(ctx).Debug("extracted auth role", zap.String("role", role))
			//
			//grantedAccess, err := a.enforcer.Enforce(role, r.URL.Path, r.Method)
			//if err != nil {
			//	middleware.Log(ctx).Error("casbin enforced resulted in an error", zap.Error(err))
			//	render.Status(r, 500)
			//	return
			//}
			//
			//if !grantedAccess {
			//	middleware.Log(ctx).Debug("not authorized")
			//	render.Respond(w, r, errors.NewRestError(403, "Forbidden"))
			//	return
			//}

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}

// TODO: Refactor into pkg
// This code is duped in cloud-api
//func extractRole(ctx context.Context, claims jwt.MapClaims) string {
//	middleware.Log(ctx).Debug("extract role from incoming claims", zap.Any("claims", claims))
//	cfg := config.GetConfig()
//	if ra, ok := claims["resource_access"].(map[string]interface{}); ok {
//		if ca, ok := ra[cfg.Identity.ClientID].(map[string]interface{}); ok {
//			if roles, ok := ca["roles"].([]interface{}); ok {
//				middleware.Log(ctx).Debug("incoming claim roles slice found", zap.Any("role", roles))
//				if found, role := checkArrayForRoles(ctx, roles); found {
//					return role
//				}
//			}
//		}
//	}
//
//	if r, ok := claims["role"].(string); ok {
//		middleware.Log(ctx).Debug("incoming claim role", zap.String("role", r))
//		return r
//	}
//
//	if groups, ok := claims["groups"].([]interface{}); ok {
//		middleware.Log(ctx).Debug("incoming claim groups slice found", zap.Any("groups", groups))
//		if found, role := checkArrayForRoles(ctx, groups); found {
//			return role
//		}
//	}
//	middleware.Log(ctx).Debug("unknown role extracted")
//	return "unknown"
//}

// func checkArrayForRoles(ctx context.Context, strings []interface{}) (bool, string) {
// 	if contains(strings, "admin") {
// 		middleware.Log(ctx).Debug("incoming claim contains admin role")
// 		return true, string(AdminRole)
// 	}
// 	if contains(strings, "user") {
// 		middleware.Log(ctx).Debug("incoming claim contains user role")
// 		return true, string(UserRole)
// 	}
// 	if contains(strings, "service") {
// 		middleware.Log(ctx).Debug("incoming claim contains service role")
// 		return true, string(ServiceRole)
// 	}
// 	if contains(strings, "guest") {
// 		middleware.Log(ctx).Debug("incoming claim contains guest role")
// 		return true, string(GuestRole)
// 	}

// 	middleware.Log(ctx).Debug("unknown role extracted")
// 	return false, ""
// }

// func contains(s []interface{}, e string) bool {
// 	for _, a := range s {
// 		if a == e {
// 			return true
// 		}
// 	}
// 	return false
// }

// const (
// 	AdminRole   Role = "admin"
// 	UserRole    Role = "user"
// 	ServiceRole Role = "service"
// 	GuestRole   Role = "guest"
// )

// type Role string
