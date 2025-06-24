package server

import (
	"benthos_go/common"
	userApp "benthos_go/user/app"
	userInfra "benthos_go/user/infra"
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// New crea y configura el servidor
func New(ctx *context.Context) *http.Server {

    mux := chi.NewMux()
    userRepo := userInfra.NewRepo()
    userService := userApp.NewService(userRepo)
    userRoutes := userInfra.NewUserRoutes(userService)

    // Aquí usamos la interfaz RouteConfigurator para asegurarnos de que Configure está implementado
    var configurators = []common.RouteConfigurator{
        userRoutes,
    }

    // Configurar las rutas de los dominios
    for _, configurator := range configurators {
        configurator.Configure(mux, ctx)
    }

    return &http.Server{
        Addr:    "127.0.0.1:3120",
        Handler: mux,
    }
}