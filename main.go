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
	var err error
	err = godotenv.Load()
	if err != nil {
		slog.Warn("Error loading .env file")
	}

	slog.Info("Connecting to the database...")
	context := context.Background()

	err = db.Connect(context)

	if err != nil {
		log.Fatal(db.Connect(context))
	}

	slog.Info("Cheking migrations...")
	err = migrations.RunMigrations(&context)

	if err != nil {
		slog.Warn(err.Error())
	}

	srv := server.New(&context)
	slog.Info("Server running in port " + strings.Split(srv.Addr, ":")[1] + "!")
	log.Fatal(srv.ListenAndServe())
}
