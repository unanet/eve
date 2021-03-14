package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"gitlab.unanet.io/devops/eve/internal/service/crud"
	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/go/pkg/errors"
	"gitlab.unanet.io/devops/go/pkg/json"
	"net/http"
	"strconv"
)

type DefinitionController struct {
	manager *crud.Manager
}

func NewDefinitionController(manager *crud.Manager) *DefinitionController {
	return &DefinitionController{
		manager: manager,
	}
}

func (c DefinitionController) Setup(r chi.Router) {
	r.Get("/definitions", c.definitions)
	r.Put("/definitions", c.upsertDefinition)
	r.Patch("/definitions", c.upsertMergeDefinition)
	r.Delete("/definitions/{definition}/{key}", c.deleteDefinitionKey)
	r.Delete("/definitions/{definition}", c.deleteDefinition)
	r.Get("/definitions/{definition}", c.getDefinition)

	r.Put("/definitions/{definition}/service-maps", c.upsertDefinitionServiceMap)
	r.Delete("/definitions/{definition}/service-maps/{description}", c.deleteServiceDefinitionMap)
	r.Get("/definitions/{definition}/service-maps", c.getServiceDefinitionMapsByDefinitionID)

	r.Put("/definitions/{definition}/job-maps", c.upsertDefinitionJobMap)
	r.Delete("/definitions/{definition}/job-maps/{description}", c.deleteJobDefinitionMap)
	r.Get("/definitions/{definition}/job-maps", c.getJobDefinitionMapsByDefinitionID)
}

func (c DefinitionController) definitions(w http.ResponseWriter, r *http.Request) {
	result, err := c.manager.Definitions(r.Context())
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, result)
}

func (c DefinitionController) upsertDefinition(w http.ResponseWriter, r *http.Request) {
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

func (c DefinitionController) upsertMergeDefinition(w http.ResponseWriter, r *http.Request) {
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

func (c DefinitionController) deleteDefinitionKey(w http.ResponseWriter, r *http.Request) {
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

func (c DefinitionController) deleteDefinition(w http.ResponseWriter, r *http.Request) {
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

func (c DefinitionController) getDefinition(w http.ResponseWriter, r *http.Request) {
	definitionID := chi.URLParam(r, "definition")
	definition, err := c.manager.GetDefinition(r.Context(), definitionID)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, definition)
}

func (c DefinitionController) upsertDefinitionServiceMap(w http.ResponseWriter, r *http.Request) {
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

func (c DefinitionController) deleteServiceDefinitionMap(w http.ResponseWriter, r *http.Request) {
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

func (c DefinitionController) getServiceDefinitionMapsByDefinitionID(w http.ResponseWriter, r *http.Request) {
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

func (c DefinitionController) upsertDefinitionJobMap(w http.ResponseWriter, r *http.Request) {
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

func (c DefinitionController) deleteJobDefinitionMap(w http.ResponseWriter, r *http.Request) {
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

func (c DefinitionController) getJobDefinitionMapsByDefinitionID(w http.ResponseWriter, r *http.Request) {
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
