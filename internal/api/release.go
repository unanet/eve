package api

import (
	"net/http"

	"github.com/unanet/eve/internal/service/releases"
	"github.com/unanet/eve/pkg/eve"
	"github.com/unanet/go/pkg/json"

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

func (c ReleaseController) Setup(r *Routers) {
	r.Auth.Post("/release", c.release)
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
