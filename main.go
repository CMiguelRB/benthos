package main

import (
	"benthos_go/server"
	"benthos_go/db"
	"context"
	"log"
	"log/slog"
	"strings"
)

func main() {

	context := context.Background()

	err := db.Connect(context)

	if err != nil {
		log.Fatal(db.Connect(context))
	}

    srv := server.New(&context)
    slog.Info("Server running in port " + strings.Split(srv.Addr, ":")[1])
    log.Fatal(srv.ListenAndServe())
}