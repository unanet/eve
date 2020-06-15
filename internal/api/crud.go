package api

import (
	"net/http"

	"github.com/go-chi/chi"

	"gitlab.unanet.io/devops/eve/internal/service/crud"
)

type CrudController struct {
	manager *crud.Manager
}

func NewCrudController(manager *crud.Manager) *CrudController {
	return &CrudController{
		manager: manager,
	}
}

func (s CrudController) Setup(r chi.Router) {
	r.Get("/environments", s.environments)

	r.Get("/services", s.services)
}

func (s CrudController) environments(w http.ResponseWriter, r *http.Request) {
	r.URL.Query().Get("version")

	s.manager.Environments()

	//err := c.planGenerator.QueuePlan(r.Context(), &options)
	//if err != nil {
	//	render.Respond(w, r, err)
	//	return
	//}
	//
	//if len(options.Messages) > 0 {
	//	render.Status(r, http.StatusPartialContent)
	//} else {
	//	render.Status(r, http.StatusAccepted)
	//}
	//render.Respond(w, r, options)
}

func (s CrudController) services(w http.ResponseWriter, r *http.Request) {

}
