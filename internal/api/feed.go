package api

import (
	"gitlab.unanet.io/devops/eve/pkg/eve"
	"net/http"
	"strconv"

	"gitlab.unanet.io/devops/eve/internal/service/crud"
	"gitlab.unanet.io/devops/go/pkg/errors"
	"gitlab.unanet.io/devops/go/pkg/json"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type FeedController struct {
	manager *crud.Manager
}

func NewFeedController(manager *crud.Manager) *FeedController {
	return &FeedController{
		manager: manager,
	}
}

func (c FeedController) Setup(r chi.Router) {
	r.Get("/feeds", c.feed)
	r.Post("/feeds", c.create)
	r.Put("/feeds/{feedID}", c.update)
	r.Delete("/feeds/{feedID}", c.delete)
}

func (c FeedController) feed(w http.ResponseWriter, r *http.Request) {

	results, err := c.manager.Feeds(r.Context())

	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, results)
}

func (c FeedController) create(w http.ResponseWriter, r *http.Request) {

	var m eve.Feed
	if err := json.ParseBody(r, &m); err != nil {
		render.Respond(w, r, err)
		return
	}

	err := c.manager.CreateFeed(r.Context(), &m)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, m)
}

func (c FeedController) update(w http.ResponseWriter, r *http.Request) {

	clusterID := chi.URLParam(r, "feedID")
	intID, err := strconv.Atoi(clusterID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid feedID in route"))
		return
	}

	var m eve.Feed
	if err := json.ParseBody(r, &m); err != nil {
		render.Respond(w, r, err)
		return
	}

	m.ID = intID

	err = c.manager.UpdateFeed(r.Context(), &m)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, m)
}

func (c FeedController) delete(w http.ResponseWriter, r *http.Request) {
	// TODO conversation is needed about if this is needed or do we do a soft delete
	render.Status(r, http.StatusNotImplemented)
	return

	clusterID := chi.URLParam(r, "feedID")
	intID, err := strconv.Atoi(clusterID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid feedID in route"))
		return
	}

	if err = c.manager.DeleteFeed(r.Context(), intID); err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}
