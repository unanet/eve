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

type AnnotationController struct {
	manager *crud.Manager
}

func NewAnnotationController(manager *crud.Manager) *AnnotationController {
	return &AnnotationController{
		manager: manager,
	}
}

func (c AnnotationController) Setup(r chi.Router) {
	r.Get("/annotations", c.annotation)
	r.Put("/annotations", c.upsertAnnotation)
	r.Patch("/annotations", c.upsertMergeAnnotation)
	r.Delete("/annotations/{annotation}/{key}", c.deleteAnnotationKey)
	r.Delete("/annotations/{annotation}", c.deleteAnnotation)
	r.Get("/annotations/{annotation}", c.getAnnotation)

	r.Put("/annotations/{annotation}/service-maps", c.upsertAnnotationServiceMap)
	r.Delete("/annotations/{annotation}/service-maps/{description}", c.deleteServiceAnnotationMap)
	r.Get("/annotations/{annotation}/service-maps", c.getServiceAnnotationMapsByAnnotationID)

	r.Put("/annotations/{annotation}/job-maps", c.upsertAnnotationJobMap)
	r.Delete("/annotations/{annotation}/job-maps/{description}", c.deleteJobAnnotationMap)
	r.Get("/annotations/{annotation}/job-maps", c.getJobAnnotationMapsByAnnotationID)
}

func (c AnnotationController) annotation(w http.ResponseWriter, r *http.Request) {
	result, err := c.manager.Annotations(r.Context())
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, result)
}

func (c AnnotationController) upsertMergeAnnotation(w http.ResponseWriter, r *http.Request) {
	var m eve.Annotation
	if err := json.ParseBody(r, &m); err != nil {
		render.Respond(w, r, err)
		return
	}

	err := c.manager.UpsertMergeAnnotation(r.Context(), &m)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.Respond(w, r, m)
}

func (c AnnotationController) upsertAnnotation(w http.ResponseWriter, r *http.Request) {
	var m eve.Annotation
	if err := json.ParseBody(r, &m); err != nil {
		render.Respond(w, r, err)
		return
	}

	err := c.manager.CreateAnnotation(r.Context(), &m)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, m)
}

func (c AnnotationController) deleteAnnotationKey(w http.ResponseWriter, r *http.Request) {
	annotationID := chi.URLParam(r, "annotation")
	intID, err := strconv.Atoi(annotationID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid annotation route parameter, required int value"))
		return
	}

	key := chi.URLParam(r, "key")
	if key == "" {
		render.Respond(w, r, errors.BadRequest("invalid key route parameter"))
		return
	}

	annotation, err := c.manager.DeleteAnnotationKey(r.Context(), intID, key)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, annotation)
}

func (c AnnotationController) deleteAnnotation(w http.ResponseWriter, r *http.Request) {
	annotationID := chi.URLParam(r, "annotation")
	intID, err := strconv.Atoi(annotationID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid annotation route parameter, required int value"))
		return
	}
	err = c.manager.DeleteAnnotation(r.Context(), intID)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (c AnnotationController) getAnnotation(w http.ResponseWriter, r *http.Request) {
	annotationID := chi.URLParam(r, "annotation")
	annotation, err := c.manager.GetAnnotation(r.Context(), annotationID)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, annotation)
}

func (c AnnotationController) upsertAnnotationServiceMap(w http.ResponseWriter, r *http.Request) {
	annotationID := chi.URLParam(r, "annotation")
	intID, err := strconv.Atoi(annotationID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid annotation route parameter, required int value"))
		return
	}

	var m eve.AnnotationServiceMap
	if err = json.ParseBody(r, &m); err != nil {
		render.Respond(w, r, err)
		return
	}

	m.AnnotationID = intID

	err = c.manager.UpsertAnnotationServiceMap(r.Context(), &m)
	if err != nil {
		render.Respond(w, r, err)
		return
	}
	render.Respond(w, r, m)
}

func (c AnnotationController) deleteServiceAnnotationMap(w http.ResponseWriter, r *http.Request) {
	annotationID := chi.URLParam(r, "annotation")
	intID, err := strconv.Atoi(annotationID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid annotation route parameter, required int value"))
		return
	}

	description := chi.URLParam(r, "description")
	if description == "" {
		render.Respond(w, r, errors.BadRequest("invalid description route parameter"))
		return
	}

	err = c.manager.DeleteAnnotationServiceMap(r.Context(), intID, description)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (c AnnotationController) getServiceAnnotationMapsByAnnotationID(w http.ResponseWriter, r *http.Request) {
	annotation := chi.URLParam(r, "annotation")
	annotationID, err := strconv.Atoi(annotation)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid annotation route parameter, required int value"))
		return
	}
	result, err := c.manager.ServiceAnnotationMapsByAnnotationID(r.Context(), annotationID)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, result)
}

func (c AnnotationController) upsertAnnotationJobMap(w http.ResponseWriter, r *http.Request) {
	annotationID := chi.URLParam(r, "annotation")
	intID, err := strconv.Atoi(annotationID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid annotation route parameter, required int value"))
		return
	}

	var m eve.AnnotationJobMap
	if err = json.ParseBody(r, &m); err != nil {
		render.Respond(w, r, err)
		return
	}

	m.AnnotationID = intID

	err = c.manager.UpsertAnnotationJobMap(r.Context(), &m)
	if err != nil {
		render.Respond(w, r, err)
		return
	}
	render.Respond(w, r, m)
}

func (c AnnotationController) deleteJobAnnotationMap(w http.ResponseWriter, r *http.Request) {
	annotationID := chi.URLParam(r, "annotation")
	intID, err := strconv.Atoi(annotationID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid annotation route parameter, required int value"))
		return
	}

	description := chi.URLParam(r, "description")
	if description == "" {
		render.Respond(w, r, errors.BadRequest("invalid description route parameter"))
		return
	}

	err = c.manager.DeleteAnnotationJobMap(r.Context(), intID, description)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (c AnnotationController) getJobAnnotationMapsByAnnotationID(w http.ResponseWriter, r *http.Request) {
	annotation := chi.URLParam(r, "annotation")
	annotationID, err := strconv.Atoi(annotation)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid annotation route parameter, required int value"))
		return
	}
	result, err := c.manager.JobAnnotationMapsByAnnotationID(r.Context(), annotationID)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, result)
}
