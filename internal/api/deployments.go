package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"gitlab.unanet.io/devops/eve/internal/service"
	"gitlab.unanet.io/devops/eve/pkg/json"
)

type DeploymentsController struct {
	planGenerator *service.PlanGenerator
}

type DeploymentRequest struct {
	Environment string                     `json:"environment"`
	Namespaces  service.StringList         `json:"namespaces"`
	Services    service.ServiceDefinitions `json:"services"`
	ForceDeploy bool                       `json:"force_deploy"`
	DryRun      bool                       `json:"dry_run"`
}

func (dr DeploymentRequest) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &dr,
		validation.Field(&dr.Environment, validation.Required))
}

func NewDeploymentsController(planGenerator *service.PlanGenerator) *DeploymentsController {
	return &DeploymentsController{
		planGenerator: planGenerator,
	}
}

func (c DeploymentsController) Setup(r chi.Router) {
	r.Post("/deployments", c.createDeployment)
}

func (c DeploymentsController) createDeployment(w http.ResponseWriter, r *http.Request) {
	var dr DeploymentRequest
	if err := json.ParseBody(r, &dr); err != nil {
		render.Respond(w, r, err)
		return
	}

	plan, err := c.planGenerator.GenerateDeploymentPlan(r.Context(), service.PlanOptions{
		Environment:      dr.Environment,
		NamespaceAliases: dr.Namespaces,
		Services:         dr.Services,
		ForceDeploy:      false,
		DryRun:           dr.DryRun,
	})

	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, plan)
}
