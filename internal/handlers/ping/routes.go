package ping

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"gitlab.unanet.io/devops/eve/internal/errors"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/internal-error", internalError)
	router.Get("/rest-error", restError)
	router.Get("/ping", ping)
	return router
}

func restError(w http.ResponseWriter, r *http.Request) {
	render.Respond(w, r, &errors.RestError{
		Code:          400,
		Message:       "Bad Request",
		OriginalError: nil,
	})
}

func internalError(w http.ResponseWriter, r *http.Request) {
	render.Respond(w, r, fmt.Errorf("Some Error"))
}

func ping(w http.ResponseWriter, r *http.Request) {
	render.Respond(w, r, render.M{
		"message": "pong",
	})
}
