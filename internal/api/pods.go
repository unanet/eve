package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"gitlab.unanet.io/devops/eve/internal/service/crud"
)

type PodController struct {
	manager *crud.Manager
}

func NewPodController(manager *crud.Manager) *PodController {
	return &PodController{
		manager: manager,
	}
}

func (s PodController) Setup(r chi.Router) {
	r.Get("/pod-resources", s.podResources)
	r.Get("/pod-autoscale", s.podAutoscale)
}

func (s PodController) podAutoscale(w http.ResponseWriter, r *http.Request) {
	serviceID := r.URL.Query().Get("service")
	environmentID := r.URL.Query().Get("environment")
	namespaceID := r.URL.Query().Get("namespace")

	result, err := s.manager.PodAutoscale(r.Context(), serviceID, environmentID, namespaceID)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, result)
}

func (s PodController) podResources(w http.ResponseWriter, r *http.Request) {
	serviceID := r.URL.Query().Get("service")
	environmentID := r.URL.Query().Get("environment")
	namespaceID := r.URL.Query().Get("namespace")
	artifactID := r.URL.Query().Get("artifact")

	result, err := s.manager.PodResources(r.Context(), serviceID, environmentID, namespaceID, artifactID)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, result)
}
