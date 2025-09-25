package main

import (
	"benthos/config"
	"benthos/db"
	"benthos/db/migrations"
	"benthos/server"
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

var Settings config.Config

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting benthos server.")
	slog.Info("Loading environment variables...")
	var err error
	err = godotenv.Load()
	if err != nil {
		slog.Warn("Error loading .env file")
	}
	if os.Getenv("ENCRYPTION_KEY") == "" || os.Getenv("DB_USER") == "" {
		log.Fatal("Server inizialization error: no environment variables found!")
	} else {
		slog.Info("Environment variables loaded OK!")
	}

	slog.Info("Loading service configuration...")
	config.InitConfiguration()
	slog.Info("Service configuration loaded OK!")

	slog.Info("Connecting to the database...")
	context := context.Background()
	err = db.Connect(context)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("Database connection OK!")

	slog.Info("Running migrations...")
	err = migrations.RunMigrations(&context)
	if err != nil {
		slog.Warn(err.Error())
	}
	slog.Info("Migrations OK!")

	srv := server.New(&context)
	slog.Info("Server running at " + srv.Addr + "!")
	slog.Info(config.Settings.App.Name + " " + config.Settings.App.Version + " ready to go!")
	err = srv.ListenAndServe()
	if err != nil {
		db.Pool.Close()
		log.Fatal(err)
	}
}
