package common

import (
	"github.com/go-chi/chi/v5"
)

type RouteSetup interface {
	Configure(mux *chi.Mux)
}

type ModuleInitializer interface {
    Initialize() RouteSetup
}

type Module[Repo any, Service any, Routes RouteSetup] struct {
	NewRepo    func() Repo
	NewService func(Repo) Service
	NewRoutes  func(Service) Routes
}

func (m Module[Repo, Service, Routes]) Initialize() RouteSetup {
	repo := m.NewRepo()
	service := m.NewService(repo)
	routes := m.NewRoutes(service)
	return routes
}
