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

type MetadataController struct {
	manager *crud.Manager
}

func NewMetadataController(manager *crud.Manager) *MetadataController {
	return &MetadataController{
		manager: manager,
	}
}

func (c MetadataController) Setup(r chi.Router) {
	r.Get("/metadata", c.metadata)
	r.Put("/metadata", c.upsertMetadata)
	r.Patch("/metadata", c.upsertMergeMetadata)
	r.Delete("/metadata/{metadata}/{key}", c.deleteMetadataKey)
	r.Delete("/metadata/{metadata}", c.deleteMetadata)
	r.Get("/metadata/{metadata}", c.getMetadata)
	r.Put("/metadata/{metadata}/service", c.upsertMetadataServiceMap)
	r.Delete("/metadata/{metadata}/service/{description}", c.deleteServiceMetadataMap)

	r.Get("/services/{service}/metadata", c.getServiceMetadata)
}

func (c MetadataController) metadata(w http.ResponseWriter, r *http.Request) {
	serviceID := r.URL.Query().Get("service")
	namespaceName := r.URL.Query().Get("namespace")

	result, err := c.manager.Metadata(r.Context(), serviceID, namespaceName)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, result)
}

func (c MetadataController) upsertMergeMetadata(w http.ResponseWriter, r *http.Request) {
	var m eve.Metadata
	if err := json.ParseBody(r, &m); err != nil {
		render.Respond(w, r, err)
		return
	}

	err := c.manager.UpsertMergeMetadata(r.Context(), &m)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.Respond(w, r, m)
}

func (c MetadataController) upsertMetadata(w http.ResponseWriter, r *http.Request) {
	var m eve.Metadata
	if err := json.ParseBody(r, &m); err != nil {
		render.Respond(w, r, err)
		return
	}

	err := c.manager.CreateMetadata(r.Context(), &m)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, m)
}

func (c MetadataController) deleteMetadataKey(w http.ResponseWriter, r *http.Request) {
	metadataID := chi.URLParam(r, "metadata")
	intID, err := strconv.Atoi(metadataID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid metadata route parameter, required int value"))
		return
	}

	key := chi.URLParam(r, "key")
	if key == "" {
		render.Respond(w, r, errors.BadRequest("invalid key route parameter"))
		return
	}

	metadata, err := c.manager.DeleteMetadataKey(r.Context(), intID, key)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, metadata)
}

func (c MetadataController) deleteMetadata(w http.ResponseWriter, r *http.Request) {
	metadataID := chi.URLParam(r, "metadata")
	intID, err := strconv.Atoi(metadataID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid metadata route parameter, required int value"))
		return
	}
	err = c.manager.DeleteMetadata(r.Context(), intID)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (c MetadataController) getMetadata(w http.ResponseWriter, r *http.Request) {
	metadataID := chi.URLParam(r, "metadata")
	metadata, err := c.manager.GetMetadata(r.Context(), metadataID)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, metadata)
}

func (c MetadataController) upsertMetadataServiceMap(w http.ResponseWriter, r *http.Request) {
	metadataID := chi.URLParam(r, "metadata")
	intID, err := strconv.Atoi(metadataID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid metadata route parameter, required int value"))
		return
	}

	var m eve.MetadataServiceMap
	if err = json.ParseBody(r, &m); err != nil {
		render.Respond(w, r, err)
		return
	}

	m.MetadataID = intID

	err = c.manager.UpsertMetadataServiceMap(r.Context(), &m)
	if err != nil {
		render.Respond(w, r, err)
		return
	}
	render.Respond(w, r, m)
}

func (c MetadataController) deleteServiceMetadataMap(w http.ResponseWriter, r *http.Request) {
	metadataID := chi.URLParam(r, "metadata")
	intID, err := strconv.Atoi(metadataID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid metadata route parameter, required int value"))
		return
	}

	description := chi.URLParam(r, "description")
	if description == "" {
		render.Respond(w, r, errors.BadRequest("invalid description route parameter"))
		return
	}

	err = c.manager.DeleteMetadataServiceMap(r.Context(), intID, description)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (c MetadataController) getServiceMetadata(w http.ResponseWriter, r *http.Request) {
	service := chi.URLParam(r, "service")
	serviceID, err := strconv.Atoi(service)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid service route parameter, required int value"))
		return
	}
	result, err := c.manager.ServiceMetadata(r.Context(), serviceID)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, result)
}
