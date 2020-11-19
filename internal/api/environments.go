package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"gitlab.unanet.io/devops/eve/internal/service/crud"
	"gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/eve/pkg/json"
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
	r.Get("/environments/{environment}", c.environment)
	r.Post("/environments/{environment}", c.updateEnvironment)
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
