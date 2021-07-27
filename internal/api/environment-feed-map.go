package api

import (
	"github.com/unanet/eve/pkg/eve"
	"net/http"

	"github.com/unanet/eve/internal/service/crud"
	"github.com/unanet/go/pkg/json"

	"github.com/go-chi/render"
)

type EnvironmentFeedMapController struct {
	manager *crud.Manager
}

func NewEnvironmentFeedMapController(manager *crud.Manager) *EnvironmentFeedMapController {
	return &EnvironmentFeedMapController{
		manager: manager,
	}
}

func (c EnvironmentFeedMapController) Setup(r *Routers) {
	r.Auth.Get("/environment-feed-maps", c.environmentFeedMaps)
	r.Auth.Post("/environment-feed-maps", c.createEnvironmentFeedMaps)
	r.Auth.Put("/environment-feed-maps", c.updateEnvironmentFeedMap)
	r.Auth.Delete("/environment-feed-maps", c.deleteEnvironmentFeedMap)
}

func (c EnvironmentFeedMapController) environmentFeedMaps(w http.ResponseWriter, r *http.Request) {

	results, err := c.manager.EnvironmentFeedMaps(r.Context())

	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, results)
}

func (c EnvironmentFeedMapController) createEnvironmentFeedMaps(w http.ResponseWriter, r *http.Request) {

	var m eve.EnvironmentFeedMap
	if err := json.ParseBody(r, &m); err != nil {
		render.Respond(w, r, err)
		return
	}

	err := c.manager.CreateEnvironmentFeedMap(r.Context(), &m)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, m)
}

func (c EnvironmentFeedMapController) updateEnvironmentFeedMap(w http.ResponseWriter, r *http.Request) {

	var m eve.EnvironmentFeedMap
	if err := json.ParseBody(r, &m); err != nil {
		render.Respond(w, r, err)
		return
	}

	err := c.manager.UpdateEnvironmentFeedMap(r.Context(), &m)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, m)
}

func (c EnvironmentFeedMapController) deleteEnvironmentFeedMap(w http.ResponseWriter, r *http.Request) {

	var m eve.EnvironmentFeedMap
	if err := json.ParseBody(r, &m); err != nil {
		render.Respond(w, r, err)
		return
	}

	err := c.manager.DeleteEnvironmentFeedMap(r.Context(), &m)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
	render.Respond(w, r, m)
}
