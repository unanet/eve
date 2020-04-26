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
	planGenerator *service.DeploymentPlanGenerator
}

type DeploymentPlanType string

const (
	ApplicationDeploymentPlan DeploymentPlanType = "application"
	MigrationDeploymentPlan   DeploymentPlanType = "migration"
)

type DeploymentRequest struct {
	Environment string                      `json:"environment"`
	Namespaces  service.StringList          `json:"namespaces"`
	Services    service.ArtifactDefinitions `json:"services"`
	ForceDeploy bool                        `json:"force_deploy"`
	CallbackURL string                      `json:"callback_url"`
	Type        DeploymentPlanType          `json:"type"`
	DryRun      bool                        `json:"dry_run"`
}

func (dr DeploymentRequest) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &dr,
		validation.Field(&dr.Environment, validation.Required),
		validation.Field(&dr.Type, validation.Required, validation.In(ApplicationDeploymentPlan, MigrationDeploymentPlan)))
}

func NewDeploymentPlanController(planGenerator *service.DeploymentPlanGenerator) *DeploymentsController {
	return &DeploymentsController{
		planGenerator: planGenerator,
	}
}

func (c DeploymentsController) Setup(r chi.Router) {
	r.Post("/deployment-plans", c.createDeployment)
}

func (c DeploymentsController) createDeployment(w http.ResponseWriter, r *http.Request) {
	var dr DeploymentRequest
	if err := json.ParseBody(r, &dr); err != nil {
		render.Respond(w, r, err)
		return
	}

	planOptions := service.DeploymentPlanOptions{
		Environment:      dr.Environment,
		NamespaceAliases: dr.Namespaces,
		Artifacts:        dr.Services,
		ForceDeploy:      false,
		DryRun:           dr.DryRun,
		CallbackURL:      dr.CallbackURL,
	}

	var plan *service.DeploymentPlan
	var err error

	if dr.Type == ApplicationDeploymentPlan {
		plan, err = c.planGenerator.GenerateApplicationPlan(r.Context(), planOptions)
	} else {
		plan, err = c.planGenerator.GenerateMigrationPlan(r.Context(), planOptions)
	}

	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, plan)
}
