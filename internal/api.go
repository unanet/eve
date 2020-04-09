package internal

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"

	eveErrors "gitlab.unanet.io/devops/eve/internal/errors"
	"gitlab.unanet.io/devops/eve/internal/handlers/ping"
	"gitlab.unanet.io/devops/eve/pkg/log"
)

func init() {
	render.Respond = func(w http.ResponseWriter, r *http.Request, v interface{}) {
		if err, ok := v.(error); ok {
			var restError *eveErrors.RestError
			if errors.As(err, &restError) {
				render.Status(r, restError.Code)
				render.DefaultResponder(w, r, restError)
				return
			}

			render.Status(r, 500)
			internalServerError := eveErrors.RestError{Code: http.StatusInternalServerError, Message: "Internal Server Error", OriginalError: err}
			log.GetHttpLogger(r).WithField("error", err).Error("Internal Server Error")
			render.DefaultResponder(w, r, internalServerError)
			return
		}
		render.DefaultResponder(w, r, v)
	}
}

func StartApi() {
	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(log.NewMiddlewareLogger())
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Timeout(60 * time.Second))

	r.Mount("/", ping.Routes())

	http.ListenAndServe(":8080", r)
}
