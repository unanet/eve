package deployments

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"gitlab.unanet.io/devops/eve/pkg/json"
)

type Controller struct {
}

type DeploymentRequest struct {
	Environment string   `json:"environment"`
	Namespaces  []string `json:"namespaces"`
	Services    []string `json:"services"`
}

func (dr DeploymentRequest) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &dr,
		validation.Field(&dr.Environment, validation.Required))
}

func New() *Controller {
	return &Controller{}
}

func (c Controller) Setup(r chi.Router) {
	r.Post("/deployments", c.createDeployment)
}

func (c Controller) createDeployment(w http.ResponseWriter, r *http.Request) {
	var dr DeploymentRequest
	if err := json.ParseBody(r, &dr); err != nil {
		render.Respond(w, r, err)
		return
	}
}
