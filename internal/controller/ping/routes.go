package ping

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"gitlab.unanet.io/devops/eve/internal/controller"
	"gitlab.unanet.io/devops/eve/pkg/eveerrs"
)

type Controller struct {
	controller.Base
}

func New() *Controller {
	return &Controller{}
}

func (c Controller) Setup(r chi.Router) {
	r.Get("/internal-error", c.internalError)
	r.Get("/rest-error", c.restError)
	r.Get("/ping", c.ping)
}

func (c Controller) restError(w http.ResponseWriter, r *http.Request) {
	render.Respond(w, r, &eveerrs.RestError{
		Code:          400,
		Message:       "Bad Request",
		OriginalError: nil,
	})
}

func (c Controller) internalError(w http.ResponseWriter, r *http.Request) {
	render.Respond(w, r, fmt.Errorf("Some Error"))
}

func (c Controller) ping(w http.ResponseWriter, r *http.Request) {
	render.Respond(w, r, render.M{
		"message": "pong",
	})
}
