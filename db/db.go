package db

import (
	"context"
	"os"
	"fmt"
	"log/slog"
	"sync"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
    Pool *pgxpool.Pool
	err error
    once sync.Once
)

func Connect(context context.Context) (error){

	once.Do(func () {
		connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

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

