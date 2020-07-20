package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/eve/pkg/json"
)

type ReleaseController struct {
}

func NewReleaseController() *ReleaseController {
	return &ReleaseController{}
}

func (c ReleaseController) Setup(r chi.Router) {
	r.Post("/release", c.releaseArtifact)
}

func (c ReleaseController) releaseArtifact(w http.ResponseWriter, r *http.Request) {
	var release eve.Release
	if err := json.ParseBody(r, &release); err != nil {
		render.Respond(w, r, err)
		return
	}
}
