package middleware

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve/pkg/eveerrs"
)

func init() {
	render.Respond = func(w http.ResponseWriter, r *http.Request, v interface{}) {
		if err, ok := v.(error); ok {
			var restError *eveerrs.RestError
			if errors.As(err, &restError) {
				render.Status(r, restError.Code)
				render.DefaultResponder(w, r, restError)
				return
			}

			render.Status(r, 500)
			internalServerError := eveerrs.RestError{Code: http.StatusInternalServerError, Message: "Internal Server Error", OriginalError: err}
			Log(r).Error("Internal Server Error", zap.Error(err))
			render.DefaultResponder(w, r, internalServerError)
			return
		}
		render.DefaultResponder(w, r, v)
	}
}

// This adds the logging automatically for outbound requests
func init() {
	http.DefaultTransport = DefaultTransport
}
