package deployments

import (
	"net/http"

	"github.com/go-chi/chi"

	"gitlab.unanet.io/devops/eve/internal/api/common"
)

type Controller struct {
	common.Base
}

func New() *Controller {
	return &Controller{}
}

func (c Controller) Setup(r chi.Router) {
	r.Route("/deployments", func(r chi.Router) {
		r.Post("/", c.createDeployment)
		r.Get("/", c.deploymentPlan)
	})
}

func (c Controller) createDeployment(w http.ResponseWriter, r *http.Request) {

}

func (c Controller) deploymentPlan(w http.ResponseWriter, r *http.Request) {

}
