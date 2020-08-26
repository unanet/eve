package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve/internal/service/crud"
	"gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/eve/pkg/json"
	"gitlab.unanet.io/devops/eve/pkg/log"
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
	r.Get("/environments/{environment}", s.environment)
	r.Post("/environments/{environment}", s.updateEnvironment)

	r.Get("/namespaces", s.namespaces)
	r.Get("/namespaces/{namespace}", s.namespace)
	r.Post("/namespaces/{namespace}", s.updateNamespace)
	r.Get("/namespaces/{namespace}/services", s.namespaceServices)
	r.Get("/namespaces/{namespace}/services/{service}", s.service)

	r.Get("/services/{service}", s.service)
	r.Post("/services/{service}", s.updateService)
	r.Patch("/services/{service}/metadata", s.updateMetadata)
	r.Delete("/services/{service}/metadata/{key}", s.deleteMetadata)
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
	if environmentID := chi.URLParam(r, "environment"); environmentID != "" {
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
	if environmentID := r.URL.Query().Get("environment"); environmentID != "" {
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
	if namespaceID := chi.URLParam(r, "namespace"); namespaceID != "" {
		namespace, err := s.manager.Namespace(r.Context(), namespaceID)
		if err != nil {
			render.Respond(w, r, err)
			return
		}
		render.Respond(w, r, namespace)
	} else {
		render.Respond(w, r, errors.NotFoundf("namespace not specified"))
		return
	}
}

func (s CrudController) namespaceServices(w http.ResponseWriter, r *http.Request) {
	if namespaceID := chi.URLParam(r, "namespace"); namespaceID != "" {
		services, err := s.manager.ServicesByNamespace(r.Context(), namespaceID)
		if err != nil {
			render.Respond(w, r, err)
			return
		}
		render.Respond(w, r, services)
	} else {
		render.Respond(w, r, errors.NotFoundf("namespace not specified"))
		return
	}
}

func (s CrudController) service(w http.ResponseWriter, r *http.Request) {
	namespaceID := r.URL.Query().Get("namespace")
	if namespaceID == "" {
		namespaceID = chi.URLParam(r, "namespace")
	}

	if serviceID := chi.URLParam(r, "service"); serviceID != "" {
		service, err := s.manager.Service(r.Context(), serviceID, namespaceID)
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

func (s CrudController) updateService(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "service")
	intID, err := strconv.Atoi(serviceID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid service in route"))
		return
	}

	var service eve.Service
	if err := json.ParseBody(r, &service); err != nil {
		render.Respond(w, r, err)
		return
	}

	log.Logger.Warn("Update Service", zap.Any("service", service), zap.Any("service.metadata", service.Metadata))

	if service.Metadata == nil {
		service.Metadata = make(map[string]interface{})
		log.Logger.Warn("Update Service Metedata nil", zap.Any("service", service), zap.Any("service.metadata", service.Metadata))
	}

	service.ID = intID
	rs, err := s.manager.UpdateService(r.Context(), &service)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, rs)
}

func (s CrudController) updateNamespace(w http.ResponseWriter, r *http.Request) {
	namespaceID := chi.URLParam(r, "namespace")
	intID, err := strconv.Atoi(namespaceID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid namespace in route"))
		return
	}

	var namespace eve.Namespace
	if err := json.ParseBody(r, &namespace); err != nil {
		render.Respond(w, r, err)
		return
	}

	namespace.ID = intID
	rs, err := s.manager.UpdateNamespace(r.Context(), &namespace)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, rs)
}

func (s CrudController) updateEnvironment(w http.ResponseWriter, r *http.Request) {
	environmentID := chi.URLParam(r, "environment")
	intID, err := strconv.Atoi(environmentID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid environment in route"))
		return
	}

	var environment eve.Environment
	if err = json.ParseBody(r, &environment); err != nil {
		render.Respond(w, r, err)
		return
	}

	environment.ID = intID
	rs, err := s.manager.UpdateEnvironment(r.Context(), &environment)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, rs)
}

func (s CrudController) updateMetadata(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "service")
	intID, err := strconv.Atoi(serviceID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid service in route"))
		return
	}

	var metadata map[string]interface{}
	if err := json.ParseBody(r, &metadata); err != nil {
		render.Respond(w, r, err)
		return
	}

	service, err := s.manager.UpdateServiceMetadata(r.Context(), intID, metadata)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, service.Metadata)
}

func (s CrudController) deleteMetadata(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "service")
	intID, err := strconv.Atoi(serviceID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid service route parameter"))
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
