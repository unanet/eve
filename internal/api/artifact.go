package api

import (
	"gitlab.unanet.io/devops/eve/internal/service/crud"
	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/go/pkg/errors"
	"gitlab.unanet.io/devops/go/pkg/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type ArtifactController struct {
	manager *crud.Manager
}

func NewArtifactController(manager *crud.Manager) *ArtifactController {
	return &ArtifactController{
		manager: manager,
	}
}

func (c ArtifactController) Setup(r chi.Router) {
	r.Get("/artifacts", c.artifacts)
	r.Post("/artifacts", c.createArtifact)
	r.Put("/artifacts/{artifactID}", c.updateArtifact)
	//r.Delete("/artifacts/{artifact}", c.deleteArtifact)
}

func (c ArtifactController) artifacts(w http.ResponseWriter, r *http.Request) {

	results, err := c.manager.Artifacts(r.Context())

	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, results)
}

func (c ArtifactController) createArtifact(w http.ResponseWriter, r *http.Request) {

	var m eve.Artifact
	if err := json.ParseBody(r, &m); err != nil {
		render.Respond(w, r, err)
		return
	}

	err := c.manager.CreateArtifact(r.Context(), &m)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, m)
}

func (c ArtifactController) updateArtifact(w http.ResponseWriter, r *http.Request) {

	artifactID := chi.URLParam(r, "artifactID")

	intID, err := strconv.Atoi(artifactID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid id route parameter, required int value"))
		return
	}

	var m eve.Artifact
	if err := json.ParseBody(r, &m); err != nil {
		render.Respond(w, r, err)
		return
	}

	m.ID = intID

	if err = c.manager.UpdateArtifact(r.Context(), &m); err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, m)
}


func (c ArtifactController) deleteArtifact(w http.ResponseWriter, r *http.Request) {
	artifactID := chi.URLParam(r, "artifactID")
	intID, err := strconv.Atoi(artifactID)
	if err != nil {
		render.Respond(w, r, errors.BadRequest("invalid id in route"))
		return
	}

	if err = c.manager.DeleteArtifact(r.Context(), intID); err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}


