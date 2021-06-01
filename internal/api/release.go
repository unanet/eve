package api

import (
	"net/http"

	"gitlab.unanet.io/devops/eve/internal/service/releases"
	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/go/pkg/json"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
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
	resp, err := c.svc.Release(r.Context(), release)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, resp)
}
