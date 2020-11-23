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

type NamespaceController struct {
	manager *crud.Manager
}

func NewNamespaceController(manager *crud.Manager) *NamespaceController {
	return &NamespaceController{
		manager: manager,
	}
}

func (c NamespaceController) Setup(r chi.Router) {
	r.Get("/namespaces", c.namespaces)
	r.Get("/namespaces/{namespace}", c.namespace)
	r.Post("/namespaces/{namespace}", c.updateNamespace)
	r.Get("/namespaces/{namespace}/services", c.namespaceServices)
	r.Get("/namespaces/{namespace}/services/{service}", c.service)
	r.Get("/namespaces/{namespace}/jobs", c.namespaceJobs)
	r.Get("/namespaces/{namespace}/jobs/{job}", c.job)
}

func (c NamespaceController) job(w http.ResponseWriter, r *http.Request) {
	namespaceID := r.URL.Query().Get("namespace")
	if namespaceID == "" {
		namespaceID = chi.URLParam(r, "namespace")
	}

	if jobID := chi.URLParam(r, "job"); jobID != "" {
		job, err := c.manager.Job(r.Context(), jobID, namespaceID)
		if err != nil {
			render.Respond(w, r, err)
			return
		}
		render.Respond(w, r, job)
	} else {
		render.Respond(w, r, errors.NotFoundf("job not specified"))
		return
	}
}

func (c NamespaceController) namespaceJobs(w http.ResponseWriter, r *http.Request) {
	if namespaceID := chi.URLParam(r, "namespace"); namespaceID != "" {
		jobs, err := c.manager.JobsByNamespace(r.Context(), namespaceID)
		if err != nil {
			render.Respond(w, r, err)
			return
		}
		render.Respond(w, r, jobs)
	} else {
		render.Respond(w, r, errors.NotFoundf("namespace not specified"))
		return
	}
}

func (c NamespaceController) namespaces(w http.ResponseWriter, r *http.Request) {
	var namespaces []eve.Namespace
	var err error
	if environmentID := r.URL.Query().Get("environment"); environmentID != "" {
		namespaces, err = c.manager.NamespacesByEnvironment(r.Context(), environmentID)
	} else {
		namespaces, err = c.manager.Namespaces(r.Context())
	}

	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, namespaces)
}

func (c NamespaceController) namespace(w http.ResponseWriter, r *http.Request) {
	if namespaceID := chi.URLParam(r, "namespace"); namespaceID != "" {
		namespace, err := c.manager.Namespace(r.Context(), namespaceID)
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

func (c NamespaceController) namespaceServices(w http.ResponseWriter, r *http.Request) {
	if namespaceID := chi.URLParam(r, "namespace"); namespaceID != "" {
		services, err := c.manager.ServicesByNamespace(r.Context(), namespaceID)
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

func (c NamespaceController) updateNamespace(w http.ResponseWriter, r *http.Request) {
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
	rs, err := c.manager.UpdateNamespace(r.Context(), &namespace)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, rs)
}

func (c NamespaceController) service(w http.ResponseWriter, r *http.Request) {
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
