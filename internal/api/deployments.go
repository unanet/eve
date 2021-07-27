package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/unanet/eve/internal/service/crud"
)

type DeploymentsController struct {
	manager *crud.Manager
}

func NewDeploymentsController(manager *crud.Manager) *DeploymentsController {
	return &DeploymentsController{
		manager: manager,
	}
}

func (c DeploymentsController) Setup(r *Routers) {
	r.Auth.Get("/deployments/{deployment}", c.deployment)
}

func (c DeploymentsController) deployment(w http.ResponseWriter, r *http.Request) {
	deploymentID := chi.URLParam(r, "deployment")

	deployment, err := c.manager.Deployment(r.Context(), deploymentID)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, deployment)
}
