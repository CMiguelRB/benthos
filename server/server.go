package server

import (
	shared "benthos/shared/infra"
	user "benthos/user/infra"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v3"
)

var Version = "version"

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
	
	slog.Info("Loading domain modules...")
	modules := []shared.ModuleInitializer{
        user.NewModule(),
    }

	configurators := make([]shared.RouteSetup, 0, len(modules))
    for _, module := range modules {
        configurators = append(configurators, module.Initialize())
    }
    
    for _, configurator := range configurators {
        configurator.Configure(mux)
    }

	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "3120")
	}

	return &http.Server{
		Addr:         fmt.Sprintf("127.0.0.1:%s", os.Getenv("PORT")),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}
