package api

import (
	"net/http"

	"gitlab.unanet.io/devops/eve/internal/service/plans"
	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/go/pkg/json"

	"github.com/go-chi/render"
)

type DeploymentPlansController struct {
	planGenerator *plans.PlanGenerator
}

func NewDeploymentPlansController(planGenerator *plans.PlanGenerator) *DeploymentPlansController {
	return &DeploymentPlansController{
		planGenerator: planGenerator,
	}
}

func (c DeploymentPlansController) Setup(r *Routers) {
	r.Auth.Post("/deployment-plans", c.createDeploymentPlan)
}

func (c DeploymentPlansController) createDeploymentPlan(w http.ResponseWriter, r *http.Request) {
	var options eve.DeploymentPlanOptions
	if err := json.ParseBody(r, &options); err != nil {
		render.Respond(w, r, err)
		return
	}

	err := c.planGenerator.QueuePlan(r.Context(), &options)
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
