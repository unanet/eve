package api

import "github.com/go-chi/chi"

type Controller interface {
	Setup(chi.Router)
}
