package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.unanet.io/devops/eve/internal/service/crud"
	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/go/pkg/json"

	"github.com/go-chi/render"
)

type DeploymentsCronController struct {
	manager *crud.Manager
}

func NewDeploymentsCronController(manager *crud.Manager) *DeploymentsCronController {
	return &DeploymentsCronController{
		manager: manager,
	}
}

func (c DeploymentsCronController) Setup(r *Routers) {
	r.Auth.Get("/deployment-crons", c.deploymentCrons)
	r.Auth.Put("/deployment-crons/{deploymentCronJob}", c.updateDeploymentCron)
	r.Auth.Post("/deployment-crons", c.createDeploymentCron)
	r.Auth.Put("/deployment-crons/{deploymentCron}", c.updateDeploymentCron)
	r.Auth.Delete("/deployment-crons/{deploymentCronJob}", c.deleteDeploymentCron)
}

func (c DeploymentsCronController) deploymentCrons(w http.ResponseWriter, r *http.Request) {

	results, err := c.manager.DeploymentCronJobs(r.Context())

	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, results)
}

func (c DeploymentsCronController) createDeploymentCron(w http.ResponseWriter, r *http.Request) {

	var m eve.DeploymentCronJob
	if err := json.ParseBody(r, &m); err != nil {
		render.Respond(w, r, err)
		return
	}

	err := c.manager.CreateDeploymentCronJob(r.Context(), &m)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, m)
}

func (c DeploymentsCronController) updateDeploymentCron(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "deploymentCronJob")

	var m eve.DeploymentCronJob
	if err := json.ParseBody(r, &m); err != nil {
		render.Respond(w, r, err)
		return
	}

	m.ID = id

	err := c.manager.UpdateDeploymentCronJob(r.Context(), &m)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, m)
}

func (c DeploymentsCronController) deleteDeploymentCron(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "deploymentCronJob")

	if err := c.manager.DeleteDeploymentCronJob(r.Context(), id); err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}
