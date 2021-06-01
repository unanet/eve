package api

import (
	"net/http"
	"strconv"

	"gitlab.unanet.io/devops/eve/internal/service/crud"
	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/go/pkg/errors"
	"gitlab.unanet.io/devops/go/pkg/json"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type EnvironmentController struct {
	manager *crud.Manager
}

func NewEnvironmentController(manager *crud.Manager) *EnvironmentController {
	return &EnvironmentController{
		manager: manager,
	}
}

func (c EnvironmentController) Setup(r chi.Router) {
	r.Get("/environments", c.environments)
	r.Post("/environments", c.createEnvironment)
	r.Get("/environments/{environment}", c.environment)
	r.Post("/environments/{environment}", c.updateEnvironment)
	//r.Delete("/environments/{environment}", c.deleteEnvironment)
}

func (c EnvironmentController) environments(w http.ResponseWriter, r *http.Request) {
	environments, err := c.manager.Environments(r.Context())
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, environments)
}

func (c EnvironmentController) environment(w http.ResponseWriter, r *http.Request) {
	if environmentID := chi.URLParam(r, "environment"); environmentID != "" {
		environment, err := c.manager.Environment(r.Context(), environmentID)
		if err != nil {
			render.Respond(w, r, err)
			return
		}
		render.Respond(w, r, environment)
	} else {
		render.Respond(w, r, errors.NotFoundf("environment not found"))
		return
	}
}

func (c EnvironmentController) updateEnvironment(w http.ResponseWriter, r *http.Request) {
	environmentID := chi.URLParam(r, "environment")
	intID, err := strconv.Atoi(environmentID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid environment in route"))
		return
	}

	var environment eve.Environment
	if err = json.ParseBody(r, &environment); err != nil {
		render.Respond(w, r, err)
		return
	}

	environment.ID = intID
	rs, err := c.manager.UpdateEnvironment(r.Context(), &environment)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, rs)
}


func (c EnvironmentController) createEnvironment(w http.ResponseWriter, r *http.Request) {

	var m eve.Environment
	if err := json.ParseBody(r, &m); err != nil {
		render.Respond(w, r, err)
		return
	}

	err := c.manager.CreateEnvironment(r.Context(), &m)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, m)
}

func (c EnvironmentController) deleteEnvironment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "environment")
	intID, err := strconv.Atoi(id)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid environment in route"))
		return
	}

	if err = c.manager.DeleteEnvironment(r.Context(), intID); err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}

