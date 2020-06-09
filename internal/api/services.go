package api

//
//import (
//	"net/http"
//
//	"github.com/go-chi/chi"
//	"github.com/go-chi/render"
//
//	"gitlab.unanet.io/devops/eve/internal/service/deployments"
//	"gitlab.unanet.io/devops/eve/pkg/json"
//)
//
//type ServiceController struct {
//	planGenerator *deployments.DeploymentPlanGenerator
//}
//
//func NewServiceController(planGenerator *deployments.DeploymentPlanGenerator) *ServiceController {
//	return &ServiceController{
//		planGenerator: planGenerator,
//	}
//}
//
//func (s ServiceController) Setup(r chi.Router) {
//	r.Post("/services", s.getServices)
//}
//
//func (s ServiceController) getServices(w http.ResponseWriter, r *http.Request) {
//
//}
//
//func (s ServiceController) createDeployment(w http.ResponseWriter, r *http.Request) {
//	var options deployments.DeploymentPlanOptions
//	if err := json.ParseBody(r, &options); err != nil {
//		render.Respond(w, r, err)
//		return
//	}
//
//	err := c.planGenerator.QueueDeploymentPlan(r.Context(), &options)
//	if err != nil {
//		render.Respond(w, r, err)
//		return
//	}
//
//	if len(options.Messages) > 0 {
//		render.Status(r, http.StatusPartialContent)
//	} else {
//		render.Status(r, http.StatusAccepted)
//	}
//	render.Respond(w, r, options)
//}
