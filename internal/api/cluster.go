package api

import (
	"net/http"
	"strconv"

	"github.com/unanet/eve/internal/service/crud"
	"github.com/unanet/eve/pkg/eve"
	"github.com/unanet/go/pkg/errors"
	"github.com/unanet/go/pkg/json"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type ClusterController struct {
	manager *crud.Manager
}

func NewClusterController(manager *crud.Manager) *ClusterController {
	return &ClusterController{
		manager: manager,
	}
}

func (c ClusterController) Setup(r *Routers) {
	r.Auth.Get("/clusters", c.cluster)
	r.Auth.Post("/clusters", c.createCluster)
	r.Auth.Put("/clusters/{clusterID}", c.updateCluster)
	//r.Auth.Delete("/clusters/{clusterID}", c.deleteCluster)
}

func (c ClusterController) cluster(w http.ResponseWriter, r *http.Request) {

	results, err := c.manager.Clusters(r.Context())

	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, results)
}

func (c ClusterController) createCluster(w http.ResponseWriter, r *http.Request) {

	var m eve.Cluster
	if err := json.ParseBody(r, &m); err != nil {
		render.Respond(w, r, err)
		return
	}

	err := c.manager.CreateCluster(r.Context(), &m)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, m)
}

func (c ClusterController) updateCluster(w http.ResponseWriter, r *http.Request) {
	// TODO conversation is needed about if this is needed or do we do a soft delete
	render.Status(r, http.StatusNotImplemented)
	return

	clusterID := chi.URLParam(r, "clusterID")

	var m eve.Cluster
	if err := json.ParseBody(r, &m); err != nil {
		render.Respond(w, r, err)
		return
	}

	m.ID = clusterID

	err := c.manager.UpdateCluster(r.Context(), &m)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, m)
}

func (c ClusterController) deleteCluster(w http.ResponseWriter, r *http.Request) {
	// TODO conversation is needed about if this is needed or do we do a soft delete
	render.Status(r, http.StatusNotImplemented)
	return

	clusterID := chi.URLParam(r, "clusterID")
	intID, err := strconv.Atoi(clusterID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid id in route"))
		return
	}

	if err = c.manager.DeleteCluster(r.Context(), intID); err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}
