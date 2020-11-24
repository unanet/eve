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

func (c JobController) updateJob(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "job")
	intID, err := strconv.Atoi(jobID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid job in route"))
		return
	}

	var job eve.Job
	if iErr := json.ParseBody(r, &job); iErr != nil {
		render.Respond(w, r, iErr)
		return
	}

	job.ID = intID
	rs, err := c.manager.UpdateJob(r.Context(), &job)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, rs)
}

func (c JobController) getJobMetadata(w http.ResponseWriter, r *http.Request) {
	job := chi.URLParam(r, "job")
	jobID, err := strconv.Atoi(job)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid job route parameter, required int value"))
		return
	}
	result, err := c.manager.JobMetadata(r.Context(), jobID)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, result)
}

func (c JobController) getJobMetadataMaps(w http.ResponseWriter, r *http.Request) {
	job := chi.URLParam(r, "job")
	jobID, err := strconv.Atoi(job)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid job route parameter, required int value"))
		return
	}
	result, err := c.manager.JobMetadataMaps(r.Context(), jobID)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, result)
}
