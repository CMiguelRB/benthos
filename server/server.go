package server

import (
	"benthos/common/infra"
	userInfra "benthos/user/infra"
	"context"
	"log/slog"
	"net/http"
	"time"
	"os"

	"github.com/go-chi/chi/v5"
)

func New(ctx *context.Context) *http.Server {

	handler := slog.NewJSONHandler(os.Stdout, nil)

    logger := slog.NewLogLogger(handler, slog.LevelError)

	mux := chi.NewMux()
	
	modules := []common.ModuleInitializer{
        userInfra.NewModule(),
    }

	configurators := make([]common.RouteSetup, 0, len(modules))
    for _, module := range modules {
        configurators = append(configurators, module.Initialize())
    }
    
    for _, configurator := range configurators {
        configurator.Configure(mux)
    }

	return &http.Server{
		Addr:         "127.0.0.1:3120",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		ErrorLog: logger,
	}
}
