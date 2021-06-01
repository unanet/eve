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

type ServiceController struct {
	manager *crud.Manager
}

func NewServiceController(manager *crud.Manager) *ServiceController {
	return &ServiceController{
		manager: manager,
	}
}

func (c ServiceController) Setup(r chi.Router) {
	r.Get("/services", c.services)
	r.Post("/services", c.create)
	r.Get("/services/{service}", c.service)
	r.Post("/services/{service}", c.updateService)
	r.Delete("/services/{service}", c.delete)
	r.Get("/services/{service}/metadata", c.getServiceMetadata)
	r.Get("/services/{service}/metadata-maps", c.getServiceMetadataMaps)
	r.Get("/services/{service}/definitions", c.getServiceDefinitionResult)
	r.Get("/services/{service}/definition-maps", c.getServiceDefinitions)

}

func (c ServiceController) service(w http.ResponseWriter, r *http.Request) {
	namespaceID := r.URL.Query().Get("namespace")
	if namespaceID == "" {
		namespaceID = chi.URLParam(r, "namespace")
	}

	if serviceID := chi.URLParam(r, "service"); serviceID != "" {
		service, err := c.manager.Service(r.Context(), serviceID, namespaceID)
		if err != nil {
			render.Respond(w, r, err)
			return
		}
		render.Respond(w, r, service)
	} else {
		render.Respond(w, r, errors.NotFoundf("service not specified"))
		return
	}
}

func (c ServiceController) updateService(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "service")
	intID, err := strconv.Atoi(serviceID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid service in route"))
		return
	}

	var service eve.Service
	if iErr := json.ParseBody(r, &service); iErr != nil {
		render.Respond(w, r, iErr)
		return
	}

	service.ID = intID
	rs, err := c.manager.UpdateService(r.Context(), &service)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, rs)
}

func (c ServiceController) getServiceMetadata(w http.ResponseWriter, r *http.Request) {
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

func (c ServiceController) getServiceMetadataMaps(w http.ResponseWriter, r *http.Request) {
	service := chi.URLParam(r, "service")
	serviceID, err := strconv.Atoi(service)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid service route parameter, required int value"))
		return
	}
	result, err := c.manager.ServiceMetadataMaps(r.Context(), serviceID)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, result)
}

func (c ServiceController) getServiceDefinitionResult(w http.ResponseWriter, r *http.Request) {
	service := chi.URLParam(r, "service")
	serviceID, err := strconv.Atoi(service)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid service route parameter, required int value"))
		return
	}
	result, err := c.manager.ServiceDefinitionResults(r.Context(), serviceID)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, result)
}

func (c ServiceController) getServiceDefinitions(w http.ResponseWriter, r *http.Request) {
	service := chi.URLParam(r, "service")
	serviceID, err := strconv.Atoi(service)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid service route parameter, required int value"))
		return
	}
	result, err := c.manager.ServiceDefinitions(r.Context(), serviceID)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, result)
}


func (c ServiceController) services(w http.ResponseWriter, r *http.Request) {

	results, err := c.manager.Services(r.Context())

	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, results)
}

func (c ServiceController) create(w http.ResponseWriter, r *http.Request) {
	var m eve.Service
	if err := json.ParseBody(r, &m); err != nil {
		render.Respond(w, r, err)
		return
	}

	err := c.manager.CreateService(r.Context(), &m)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, m)
}

func (c ServiceController) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "service")
	intID, err := strconv.Atoi(id)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid service in route"))
		return
	}

	if err = c.manager.DeleteService(r.Context(), intID); err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}
