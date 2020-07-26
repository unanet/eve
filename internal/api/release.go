package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"gitlab.unanet.io/devops/eve/internal/service/releases"
	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/eve/pkg/json"
)

type ReleaseController struct {
	svc *releases.ReleaseSvc
}

func NewReleaseController(s *releases.ReleaseSvc) *ReleaseController {
	return &ReleaseController{
		svc: s,
	}
}

func (c ReleaseController) Setup(r chi.Router) {
	r.Post("/promote", c.promoteArtifact)
	r.Post("/demote", c.demoteArtifact)
}

func (c ReleaseController) promoteArtifact(w http.ResponseWriter, r *http.Request) {
	var release eve.Release
	if err := json.ParseBody(r, &release); err != nil {
		render.Respond(w, r, err)
		return
	}
	err := c.svc.PromoteRelease(r.Context(), release)
	if err != nil {
		render.Respond(w, r, err)
		return
	}
	render.Respond(w, r, "success")
}

func (c ReleaseController) demoteArtifact(w http.ResponseWriter, r *http.Request) {
	var release eve.Release
	if err := json.ParseBody(r, &release); err != nil {
		render.Respond(w, r, err)
		return
	}
	err := c.svc.DemoteRelease(r.Context(), release)
	if err != nil {
		render.Respond(w, r, err)
		return
	}
	render.Respond(w, r, "success")
}
