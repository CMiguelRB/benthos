package server

import (
	"benthos/config"
	user "benthos/user/infra"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v3"
	"github.com/go-chi/httprate"
)

func New(ctx *context.Context) *http.Server {

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
		config.Settings.Server.RateLimit.Period*time.Millisecond,
		httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint),
	))

	slog.Info("Loading domain modules...")
	user.NewModule().Routes.Configure(mux)

	if config.Settings.Server.WebStaticEnabled {
		slog.Info("Serving statics...")

		mux.Handle("/assets/*", cacheControlMiddleware(http.StripPrefix("/assets/", http.FileServer(http.Dir(filepath.Join(config.Settings.Server.WebStaticDir, "assets"))))))

		mux.Handle("/favicon.svg", http.FileServer(http.Dir(config.Settings.Server.WebStaticDir)))

		mux.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, filepath.Join(config.Settings.Server.WebStaticDir, "index.html"))
		}))
	}

	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "3800")
	}

	if os.Getenv("HOSTNAME") == "" {
		os.Setenv("HOSTNAME", "0.0.0.0")
	}

	return &http.Server{
		Addr:         fmt.Sprintf("%s:%s", os.Getenv("HOSTNAME"), os.Getenv("PORT")),
		Handler:      mux,
		ReadTimeout:  config.Settings.Server.ReadTimeout * time.Millisecond,
		WriteTimeout: config.Settings.Server.WriteTimeout * time.Millisecond,
		IdleTimeout:  config.Settings.Server.IdleTimeout * time.Millisecond,
	}
}

func cacheControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		next.ServeHTTP(w, r)
	})
}
