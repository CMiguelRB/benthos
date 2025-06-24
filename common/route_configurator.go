package common

import "github.com/go-chi/chi/v5"

type RouteConfigurator interface {
	Configure(mux *chi.Mux)
}
