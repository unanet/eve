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

type JobController struct {
	manager *crud.Manager
}

func NewJobController(manager *crud.Manager) *JobController {
	return &JobController{
		manager: manager,
	}
}

func (c JobController) Setup(r chi.Router) {
	r.Get("/jobs/{job}", c.job)
	r.Post("/jobs/{job}", c.updateJob)
	r.Get("/jobs/{job}/metadata", c.getJobMetadata)
	r.Get("/jobs/{job}/metadata-maps", c.getJobMetadataMaps)
}

func (c JobController) job(w http.ResponseWriter, r *http.Request) {
	namespaceID := r.URL.Query().Get("namespace")
	if namespaceID == "" {
		namespaceID = chi.URLParam(r, "namespace")
	}

	if jobID := chi.URLParam(r, "job"); jobID != "" {
		service, err := c.manager.Service(r.Context(), jobID, namespaceID)
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

func (c JobController) updateJob(w http.ResponseWriter, r *http.Request) {
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

func (c JobController) getJobMetadata(w http.ResponseWriter, r *http.Request) {
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

func (c JobController) getJobMetadataMaps(w http.ResponseWriter, r *http.Request) {
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
