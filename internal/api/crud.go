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

type CrudController struct {
	manager *crud.Manager
}

func NewCrudController(manager *crud.Manager) *CrudController {
	return &CrudController{
		manager: manager,
	}
}

func (s CrudController) Setup(r chi.Router) {
	r.Get("/environments", s.environments)
	r.Get("/environments/{environmentID}", s.environment)

	r.Get("/namespaces", s.namespaces)
	r.Get("/namespaces/{namespaceID}", s.namespace)
	r.Get("/namespaces/{namespaceID}/services", s.namespaceServices)
	r.Get("/namespaces/{namespaceID}/services/{serviceID}", s.service)

	r.Get("/services/{serviceID}", s.service)
	r.Post("/services/{serviceID}", s.updateService)
	r.Post("/services/{serviceID}/metadata", s.updateMetadata)
	r.Delete("/services/{serviceID}/metadata/{key}", s.deleteMetadata)

}

func (s CrudController) environments(w http.ResponseWriter, r *http.Request) {
	environments, err := s.manager.Environments(r.Context())
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, environments)
}

func (s CrudController) environment(w http.ResponseWriter, r *http.Request) {
	if environmentID := chi.URLParam(r, "environmentID"); environmentID != "" {
		environment, err := s.manager.Environment(r.Context(), environmentID)
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

func (s CrudController) namespaces(w http.ResponseWriter, r *http.Request) {
	var namespaces []eve.Namespace
	var err error
	if environmentID := r.URL.Query().Get("environmentID"); environmentID != "" {
		namespaces, err = s.manager.NamespacesByEnvironment(r.Context(), environmentID)
	} else {
		namespaces, err = s.manager.Namespaces(r.Context())
	}

	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, namespaces)
}

func (s CrudController) namespace(w http.ResponseWriter, r *http.Request) {
	if namespaceID := chi.URLParam(r, "namespaceID"); namespaceID != "" {
		namespace, err := s.manager.Namespace(r.Context(), namespaceID)
		if err != nil {
			render.Respond(w, r, err)
			return
		}
		render.Respond(w, r, namespace)
	} else {
		render.Respond(w, r, errors.NotFoundf("namespaceID not specified"))
		return
	}
}

func (s CrudController) namespaceServices(w http.ResponseWriter, r *http.Request) {
	if namespaceID := chi.URLParam(r, "namespaceID"); namespaceID != "" {
		services, err := s.manager.ServicesByNamespace(r.Context(), namespaceID)
		if err != nil {
			render.Respond(w, r, err)
			return
		}
		render.Respond(w, r, services)
	} else {
		render.Respond(w, r, errors.NotFoundf("namespaceID not specified"))
		return
	}
}

func (s CrudController) service(w http.ResponseWriter, r *http.Request) {
	namespaceID := r.URL.Query().Get("namespaceID")
	if namespaceID == "" {
		namespaceID = chi.URLParam(r, "namespaceID")
	}

	if serviceID := chi.URLParam(r, "serviceID"); serviceID != "" {
		service, err := s.manager.Service(r.Context(), serviceID, namespaceID)
		if err != nil {
			render.Respond(w, r, err)
			return
		}
		render.Respond(w, r, service)
	} else {
		render.Respond(w, r, errors.NotFoundf("serviceID not specified"))
		return
	}
}

func (s CrudController) updateService(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "serviceID")
	intID, err := strconv.Atoi(serviceID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid service id in route"))
		return
	}

	var service eve.Service
	if err := json.ParseBody(r, &service); err != nil {
		render.Respond(w, r, err)
		return
	}

	service.ID = intID
	rs, err := s.manager.UpdateService(r.Context(), &service)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, rs)
}

func (s CrudController) updateMetadata(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "serviceID")
	intID, err := strconv.Atoi(serviceID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid service id in route"))
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		render.Respond(w, r, errors.BadRequest("you must specify a 'key' query parameter"))
		return
	}

	value := r.URL.Query().Get("value")
	if value == "" {
		render.Respond(w, r, errors.BadRequest("you must specify a 'value' query parameter"))
		return
	}

	service, err := s.manager.UpdateServiceMetadata(r.Context(), intID, key, value)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, service.Metadata)
}

func (s CrudController) deleteMetadata(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "serviceID")
	intID, err := strconv.Atoi(serviceID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid serviceID route parameter"))
		return
	}

	key := chi.URLParam(r, "key")
	if key == "" {
		render.Respond(w, r, errors.BadRequest("invalid key route parameter"))
		return
	}

	service, err := s.manager.DeleteServiceMetadata(r.Context(), intID, key)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, service.Metadata)
}
