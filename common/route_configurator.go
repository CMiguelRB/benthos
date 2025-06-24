package common

import "github.com/go-chi/chi/v5"

type RouteSetup interface {
	Configure(mux *chi.Mux)
}
