package server

import (
	"benthos/common"
	userApp "benthos/user/app"
	userInfra "benthos/user/infra"
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func New(ctx *context.Context) *http.Server {

	mux := chi.NewMux()
	userRepo := userInfra.NewRepo()
	userService := userApp.NewService(userRepo)
	userRoutes := userInfra.NewUserRoutes(userService)

	var routesSetup = []common.RouteSetup{
		userRoutes,
	}

	for _, routeSetup := range routesSetup {
		routeSetup.Configure(mux)
	}

	return &http.Server{
		Addr:    "127.0.0.1:3120",
		Handler: mux,
	}
}
