package main

import (
	"benthos/db"
	"benthos/db/migrations"
	"benthos/server"
	"context"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

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
	}else{
		slog.Info("Environment variables loaded OK!")
	}

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
	slog.Info("Server running at port " + strings.Split(srv.Addr, ":")[1] + "!")
	err = srv.ListenAndServe()
	if err != nil {
		db.Pool.Close();
		log.Fatal(err)
	}
}
