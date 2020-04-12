package mux

import "github.com/go-chi/chi"

type EveController interface {
	Setup(chi.Router)
}
