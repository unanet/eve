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
	r.Post("/release", c.release)
}

func (c ReleaseController) release(w http.ResponseWriter, r *http.Request) {
	var release eve.Release
	if err := json.ParseBody(r, &release); err != nil {
		render.Respond(w, r, err)
		return
	}
	msg, err := c.svc.Release(r.Context(), release)
	if err != nil {
		render.Respond(w, r, err)
		return
	}
	render.Respond(w, r, msg)
}
