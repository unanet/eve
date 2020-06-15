package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

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
	environments, err := s.manager.Environments(r.Context())
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, environments)
}

func (s CrudController) services(w http.ResponseWriter, r *http.Request) {

}
