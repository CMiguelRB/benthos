package server

import (
	commonInfra "benthos/common/infra"
	userInfra "benthos/user/infra"
	"context"
	"net/http"
	"time"
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v3"
)

const Version string = "version"

func New(ctx *context.Context) *http.Server {

	os.Setenv("VERSION", Version)

	mux := chi.NewMux()


	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		slog.String("app", os.Getenv("NAME")),
		slog.String("version", Version),
		slog.String("env", os.Getenv("ENV")),
	)

	mux.Use(httplog.RequestLogger(logger, &httplog.Options{
		Level: slog.LevelInfo,

		Schema: httplog.SchemaECS,
		RecoverPanics: true,
	}))
	
	modules := []commonInfra.ModuleInitializer{
        userInfra.NewModule(),
    }

	configurators := make([]commonInfra.RouteSetup, 0, len(modules))
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
	}
}
