package main

import (
	"benthos/db"
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

	var err error
	err = godotenv.Load()
	if err != nil {
		slog.Warn("Error loading .env file")
	}

	context := context.Background()

	err = db.Connect(context)

	if err != nil {
		log.Fatal(db.Connect(context))
	}

	srv := server.New(&context)
	slog.Info("Server running in port " + strings.Split(srv.Addr, ":")[1])
	log.Fatal(srv.ListenAndServe())
}
