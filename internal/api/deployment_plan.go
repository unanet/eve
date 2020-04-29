package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"gitlab.unanet.io/devops/eve/internal/service"
	"gitlab.unanet.io/devops/eve/pkg/json"
)

type DeploymentsController struct {
	planGenerator *service.DeploymentPlanGenerator
}

type DeploymentPlanType string

func NewDeploymentPlanController(planGenerator *service.DeploymentPlanGenerator) *DeploymentsController {
	return &DeploymentsController{
		planGenerator: planGenerator,
	}
}

func (c DeploymentsController) Setup(r chi.Router) {
	r.Post("/deployment-plans", c.createDeployment)
}

func (c DeploymentsController) createDeployment(w http.ResponseWriter, r *http.Request) {
	var options service.DeploymentPlanOptions
	if err := json.ParseBody(r, &options); err != nil {
		render.Respond(w, r, err)
		return
	}

	err := c.planGenerator.QueueDeploymentPlan(r.Context(), &options)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	if len(options.Messages) > 0 {
		render.Status(r, http.StatusPartialContent)
	} else {
		render.Status(r, http.StatusAccepted)
	}
	render.Respond(w, r, options)
}
