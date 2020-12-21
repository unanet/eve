package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"gitlab.unanet.io/devops/eve/internal/service/crud"
	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/go/pkg/errors"
	"gitlab.unanet.io/devops/go/pkg/json"
)

type LabelController struct {
	manager *crud.Manager
}

func NewLabelController(manager *crud.Manager) *LabelController {
	return &LabelController{
		manager: manager,
	}
}

func (c LabelController) Setup(r chi.Router) {
	r.Get("/labels", c.label)
	r.Put("/labels", c.upsertLabel)
	r.Patch("/labels", c.upsertMergeLabel)
	r.Delete("/labels/{label}/{key}", c.deleteLabelKey)
	r.Delete("/labels/{label}", c.deleteLabel)
	r.Get("/labels/{label}", c.getLabel)

	r.Put("/labels/{label}/service-maps", c.upsertLabelServiceMap)
	r.Delete("/labels/{label}/service-maps/{description}", c.deleteServiceLabelMap)
	r.Get("/labels/{label}/service-maps", c.getServiceLabelMapsByLabelID)

	r.Put("/labels/{label}/job-maps", c.upsertLabelJobMap)
	r.Delete("/labels/{label}/job-maps/{description}", c.deleteJobLabelMap)
	r.Get("/labels/{label}/job-maps", c.getJobLabelMapsByLabelID)
}

func (c LabelController) label(w http.ResponseWriter, r *http.Request) {
	result, err := c.manager.Labels(r.Context())
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, result)
}

func (c LabelController) upsertMergeLabel(w http.ResponseWriter, r *http.Request) {
	var m eve.Label
	if err := json.ParseBody(r, &m); err != nil {
		render.Respond(w, r, err)
		return
	}

	err := c.manager.UpsertMergeLabel(r.Context(), &m)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.Respond(w, r, m)
}

func (c LabelController) upsertLabel(w http.ResponseWriter, r *http.Request) {
	var m eve.Label
	if err := json.ParseBody(r, &m); err != nil {
		render.Respond(w, r, err)
		return
	}

	err := c.manager.CreateLabel(r.Context(), &m)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, m)
}

func (c LabelController) deleteLabelKey(w http.ResponseWriter, r *http.Request) {
	labelID := chi.URLParam(r, "label")
	intID, err := strconv.Atoi(labelID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid label route parameter, required int value"))
		return
	}

	key := chi.URLParam(r, "key")
	if key == "" {
		render.Respond(w, r, errors.BadRequest("invalid key route parameter"))
		return
	}

	label, err := c.manager.DeleteLabelKey(r.Context(), intID, key)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, label)
}

func (c LabelController) deleteLabel(w http.ResponseWriter, r *http.Request) {
	labelID := chi.URLParam(r, "label")
	intID, err := strconv.Atoi(labelID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid label route parameter, required int value"))
		return
	}
	err = c.manager.DeleteLabel(r.Context(), intID)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (c LabelController) getLabel(w http.ResponseWriter, r *http.Request) {
	labelID := chi.URLParam(r, "label")
	label, err := c.manager.GetLabel(r.Context(), labelID)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, label)
}

func (c LabelController) upsertLabelServiceMap(w http.ResponseWriter, r *http.Request) {
	labelID := chi.URLParam(r, "label")
	intID, err := strconv.Atoi(labelID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid label route parameter, required int value"))
		return
	}

	var m eve.LabelServiceMap
	if err = json.ParseBody(r, &m); err != nil {
		render.Respond(w, r, err)
		return
	}

	m.LabelID = intID

	err = c.manager.UpsertLabelServiceMap(r.Context(), &m)
	if err != nil {
		render.Respond(w, r, err)
		return
	}
	render.Respond(w, r, m)
}

func (c LabelController) deleteServiceLabelMap(w http.ResponseWriter, r *http.Request) {
	labelID := chi.URLParam(r, "label")
	intID, err := strconv.Atoi(labelID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid label route parameter, required int value"))
		return
	}

	description := chi.URLParam(r, "description")
	if description == "" {
		render.Respond(w, r, errors.BadRequest("invalid description route parameter"))
		return
	}

	err = c.manager.DeleteLabelServiceMap(r.Context(), intID, description)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (c LabelController) getServiceLabelMapsByLabelID(w http.ResponseWriter, r *http.Request) {
	label := chi.URLParam(r, "label")
	labelID, err := strconv.Atoi(label)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid label route parameter, required int value"))
		return
	}
	result, err := c.manager.ServiceLabelMapsByLabelID(r.Context(), labelID)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, result)
}

func (c LabelController) upsertLabelJobMap(w http.ResponseWriter, r *http.Request) {
	labelID := chi.URLParam(r, "label")
	intID, err := strconv.Atoi(labelID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid label route parameter, required int value"))
		return
	}

	var m eve.LabelJobMap
	if err = json.ParseBody(r, &m); err != nil {
		render.Respond(w, r, err)
		return
	}

	m.LabelID = intID

	err = c.manager.UpsertLabelJobMap(r.Context(), &m)
	if err != nil {
		render.Respond(w, r, err)
		return
	}
	render.Respond(w, r, m)
}

func (c LabelController) deleteJobLabelMap(w http.ResponseWriter, r *http.Request) {
	labelID := chi.URLParam(r, "label")
	intID, err := strconv.Atoi(labelID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid label route parameter, required int value"))
		return
	}

	description := chi.URLParam(r, "description")
	if description == "" {
		render.Respond(w, r, errors.BadRequest("invalid description route parameter"))
		return
	}

	err = c.manager.DeleteLabelJobMap(r.Context(), intID, description)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (c LabelController) getJobLabelMapsByLabelID(w http.ResponseWriter, r *http.Request) {
	label := chi.URLParam(r, "label")
	labelID, err := strconv.Atoi(label)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid label route parameter, required int value"))
		return
	}
	result, err := c.manager.ServiceLabelMapsByLabelID(r.Context(), labelID)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, result)
}
