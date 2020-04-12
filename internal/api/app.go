package api

import (
	"net/http"
	"time"

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
}

func NewApp() (*App, error) {
	client, err := artifactory.NewClient(config.Values.ArtifactoryConfig)
	if err != nil {
		return nil, err
	}
	return &App{
		r: chi.NewMux(),
		Controllers: []mux.EveController{
			ping.New(),
		},
		Artifactory: client,
	}, nil
}

func (a *App) Start() {
	a.setup()
	metrics.StartMetrics()
	http.ListenAndServe(":8080", a.r)
}

func (a *App) setup() {
	middleware.SetupMiddleware(a.r, 60*time.Second)
	for _, c := range a.Controllers {
		c.Setup(a.r)
	}
}
