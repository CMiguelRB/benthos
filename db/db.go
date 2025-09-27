package db

import (
	"benthos/config"
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	Pool *pgxpool.Pool
	err  error
	once sync.Once
)

func Connect(context context.Context) error {

	once.Do(func() {
		connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			config.Settings.Database.Username,
			config.Settings.Database.Password,
			config.Settings.Database.Hostname,
			config.Settings.Database.Port,
			config.Settings.Database.Name)
			
		Pool, err = pgxpool.New(context, connectionString)
		if err != nil {
			slog.Error("Unable to connect to the database")
			return
		}

		err = Pool.Ping(context)

		if err != nil {
			slog.Error("Unable to connect to the database")
		}
	})
	return err
}
