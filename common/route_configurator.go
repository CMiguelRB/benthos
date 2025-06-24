package common

import (
	"context"

	"github.com/go-chi/chi/v5"
)

// RouteConfigurator es la interfaz que obliga a cada dominio a implementar Configure
type RouteConfigurator interface {
    Configure(mux *chi.Mux, ctx *context.Context)
}