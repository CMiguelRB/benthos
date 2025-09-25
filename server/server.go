package server

import (
	"benthos/config"
	user "benthos/user/infra"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v3"
	"github.com/go-chi/httprate"
)

var Version = "version"

func New(ctx *context.Context) *http.Server {

	if os.Getenv("ENV") != "DEV" {
		config.Settings.App.Version = Version
	}

	mux := chi.NewMux()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		slog.String("app", os.Getenv("NAME")),
		slog.String("version", config.Settings.App.Version),
		slog.String("env", os.Getenv("ENV")),
	)

	mux.Use(httplog.RequestLogger(logger, &httplog.Options{
		Level: slog.LevelInfo,

		Schema:        httplog.SchemaECS,
		RecoverPanics: true,
	}))

	mux.Use(httprate.Limit(
		config.Settings.Server.RateLimit.Requests,
		time.Duration(config.Settings.Server.RateLimit.PeriodMs),
		httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint),
	))

	slog.Info("Loading domain modules...")
	user.NewModule().Routes.Configure(mux)

	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "3800")
	}

	return &http.Server{
		Addr:         fmt.Sprintf("%s:%s", os.Getenv("HOSTNAME"), os.Getenv("PORT")),
		Handler:      mux,
		ReadTimeout:  config.Settings.Server.ReadTimeoutMs*time.Millisecond,
		WriteTimeout: config.Settings.Server.WriteTimeoutMs*time.Millisecond,
		IdleTimeout:  config.Settings.Server.IdleTimeoutMs*time.Millisecond,
	}
}
