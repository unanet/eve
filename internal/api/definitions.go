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

type DefinitionsController struct {
	manager *crud.Manager
}

func NewDefinitionsController(manager *crud.Manager) *DefinitionsController {
	return &DefinitionsController{
		manager: manager,
	}
}

func (c DefinitionsController) Setup(r *Routers) {
	r.Auth.Get("/definitions", c.definitions)
	r.Auth.Put("/definitions", c.upsertDefinition)
	r.Auth.Patch("/definitions", c.upsertMergeDefinition)

	// Sorted here to go above thee definitions by id
	r.Auth.Get("/definitions/job-maps", c.definitionJobMaps)
	r.Auth.Get("/definitions/service-maps", c.definitionServiceMaps)

	r.Auth.Delete("/definitions/{definition}/{key}", c.deleteDefinitionKey)
	r.Auth.Delete("/definitions/{definition}", c.deleteDefinition)
	r.Auth.Get("/definitions/{definition}", c.getDefinition)

	r.Auth.Put("/definitions/{definition}/service-maps", c.upsertDefinitionServiceMap)
	r.Auth.Delete("/definitions/{definition}/service-maps/{description}", c.deleteServiceDefinitionMap)
	r.Auth.Get("/definitions/{definition}/service-maps", c.getServiceDefinitionMapsByDefinitionID)

	r.Auth.Put("/definitions/{definition}/job-maps", c.upsertDefinitionJobMap)
	r.Auth.Delete("/definitions/{definition}/job-maps/{description}", c.deleteJobDefinitionMap)
	r.Auth.Get("/definitions/{definition}/job-maps", c.getJobDefinitionMapsByDefinitionID)

	r.Auth.Get("/definition-types", c.definitionTypes)
	r.Auth.Post("/definition-types", c.createDefinitionType)
	r.Auth.Put("/definition-types/{definitionType}", c.updateDefinitionType)
	//r.Auth.Delete("/definition-types/{definitionType}", c.deleteDefinitionType)
}

func (c DefinitionsController) definitions(w http.ResponseWriter, r *http.Request) {
	result, err := c.manager.Definitions(r.Context())
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, result)
}

func (c DefinitionsController) upsertDefinition(w http.ResponseWriter, r *http.Request) {
	var m eve.Definition
	if err := json.ParseBody(r, &m); err != nil {
		render.Respond(w, r, err)
		return
	}

	err := c.manager.CreateDefinition(r.Context(), &m)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, m)
}

func (c DefinitionsController) upsertMergeDefinition(w http.ResponseWriter, r *http.Request) {
	var m eve.Definition
	if err := json.ParseBody(r, &m); err != nil {
		render.Respond(w, r, err)
		return
	}

	err := c.manager.UpsertMergeDefinition(r.Context(), &m)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.Respond(w, r, m)
}

func (c DefinitionsController) deleteDefinitionKey(w http.ResponseWriter, r *http.Request) {
	definitionID := chi.URLParam(r, "definition")
	intID, err := strconv.Atoi(definitionID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid definition route parameter, required int value"))
		return
	}

	key := chi.URLParam(r, "key")
	if key == "" {
		render.Respond(w, r, errors.BadRequest("invalid key route parameter"))
		return
	}

	definition, err := c.manager.DeleteDefinitionKey(r.Context(), intID, key)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, definition)
}

func (c DefinitionsController) deleteDefinition(w http.ResponseWriter, r *http.Request) {
	definitionID := chi.URLParam(r, "definition")
	intID, err := strconv.Atoi(definitionID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid definition route parameter, required int value"))
		return
	}
	err = c.manager.DeleteDefinition(r.Context(), intID)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (c DefinitionsController) getDefinition(w http.ResponseWriter, r *http.Request) {
	definitionID := chi.URLParam(r, "definition")
	definition, err := c.manager.GetDefinition(r.Context(), definitionID)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, definition)
}

func (c DefinitionsController) upsertDefinitionServiceMap(w http.ResponseWriter, r *http.Request) {
	definitionID := chi.URLParam(r, "definition")
	intID, err := strconv.Atoi(definitionID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid definition route parameter, required int value"))
		return
	}

	var m eve.DefinitionServiceMap
	if err = json.ParseBody(r, &m); err != nil {
		render.Respond(w, r, err)
		return
	}

	m.DefinitionID = intID

	err = c.manager.UpsertDefinitionServiceMap(r.Context(), &m)
	if err != nil {
		render.Respond(w, r, err)
		return
	}
	render.Respond(w, r, m)
}

func (c DefinitionsController) deleteServiceDefinitionMap(w http.ResponseWriter, r *http.Request) {
	definitionID := chi.URLParam(r, "definition")
	intID, err := strconv.Atoi(definitionID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid definition route parameter, required int value"))
		return
	}

	description := chi.URLParam(r, "description")
	if description == "" {
		render.Respond(w, r, errors.BadRequest("invalid description route parameter"))
		return
	}

	err = c.manager.DeleteDefinitionServiceMap(r.Context(), intID, description)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (c DefinitionsController) getServiceDefinitionMapsByDefinitionID(w http.ResponseWriter, r *http.Request) {
	definitionID := chi.URLParam(r, "definition")
	intID, err := strconv.Atoi(definitionID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid definition route parameter, required int value"))
		return
	}
	result, err := c.manager.ServiceDefinitionMapsByDefinitionID(r.Context(), intID)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, result)
}

func (c DefinitionsController) upsertDefinitionJobMap(w http.ResponseWriter, r *http.Request) {
	definitionID := chi.URLParam(r, "definition")
	intID, err := strconv.Atoi(definitionID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid definition route parameter, required int value"))
		return
	}

	var m eve.DefinitionJobMap
	if err = json.ParseBody(r, &m); err != nil {
		render.Respond(w, r, err)
		return
	}

	m.DefinitionID = intID

	err = c.manager.UpsertDefinitionJobMap(r.Context(), &m)
	if err != nil {
		render.Respond(w, r, err)
		return
	}
	render.Respond(w, r, m)
}

func (c DefinitionsController) deleteJobDefinitionMap(w http.ResponseWriter, r *http.Request) {
	definitionID := chi.URLParam(r, "definition")
	intID, err := strconv.Atoi(definitionID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid definition route parameter, required int value"))
		return
	}

	description := chi.URLParam(r, "description")
	if description == "" {
		render.Respond(w, r, errors.BadRequest("invalid description route parameter"))
		return
	}

	err = c.manager.DeleteDefinitionJobMap(r.Context(), intID, description)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (c DefinitionsController) getJobDefinitionMapsByDefinitionID(w http.ResponseWriter, r *http.Request) {
	definitionID := chi.URLParam(r, "definition")
	intID, err := strconv.Atoi(definitionID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid definition route parameter, required int value"))
		return
	}
	result, err := c.manager.JobDefinitionMapsByDefinitionID(r.Context(), intID)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, result)
}

func (c DefinitionsController) definitionTypes(w http.ResponseWriter, r *http.Request) {

	results, err := c.manager.DefinitionTypes(r.Context())

	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, results)
}

func (c DefinitionsController) createDefinitionType(w http.ResponseWriter, r *http.Request) {

	var m eve.DefinitionType
	if err := json.ParseBody(r, &m); err != nil {
		render.Respond(w, r, err)
		return
	}

	err := c.manager.CreateDefinitionType(r.Context(), &m)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, m)
}

func (c DefinitionsController) updateDefinitionType(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "definitionType")
	intID, err := strconv.Atoi(id)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid definitionType in route"))
		return
	}

	var m eve.DefinitionType
	if err := json.ParseBody(r, &m); err != nil {
		render.Respond(w, r, err)
		return
	}

	m.ID = intID

	err = c.manager.UpdateDefinitionType(r.Context(), &m)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, m)
}

func (c DefinitionsController) deleteDefinitionType(w http.ResponseWriter, r *http.Request) {
	// TODO conversation is needed about if this is needed or do we do a soft delete
	render.Status(r, http.StatusNotImplemented)
	return

	clusterID := chi.URLParam(r, "definitionType")
	intID, err := strconv.Atoi(clusterID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid definitionType in route"))
		return
	}

	if err = c.manager.DeleteDefinitionType(r.Context(), intID); err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (c DefinitionsController) definitionJobMaps(w http.ResponseWriter, r *http.Request) {

	results, err := c.manager.DefinitionJobMaps(r.Context())

	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, results)
}

func (c DefinitionsController) definitionServiceMaps(w http.ResponseWriter, r *http.Request) {

	results, err := c.manager.DefinitionServiceMaps(r.Context())

	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, results)
}
