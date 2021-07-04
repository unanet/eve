package api

import "github.com/go-chi/chi"

type Routers struct {
	Auth      chi.Router
	Anonymous chi.Router
}

type Controller interface {
	Setup(*Routers)
}
