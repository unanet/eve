package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"gitlab.unanet.io/devops/eve/internal/service/plans"
	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/go/pkg/json"
)

type DeploymentsController struct {
	planGenerator *plans.PlanGenerator
}

func NewDeploymentPlanController(planGenerator *plans.PlanGenerator) *DeploymentsController {
	return &DeploymentsController{
		planGenerator: planGenerator,
	}
}

func (c DeploymentsController) Setup(r chi.Router) {
	r.Post("/deployment-plans", c.createDeploymentPlan)
}

func (c DeploymentsController) createDeploymentPlan(w http.ResponseWriter, r *http.Request) {
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
