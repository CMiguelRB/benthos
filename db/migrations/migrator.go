package migrations

import (
	"benthos/db"
	"context"
	"database/sql"
	"log"
	"log/slog"
	"sort"

	"github.com/jackc/pgx/v5"
)

type Migration struct {
	Up   func(ctx *context.Context) error
	Down func(ctx *context.Context) error
}

type SchemaVersion struct {
	Id        int
	Timestamp sql.NullTime
	Name      string
	Success   bool
}

var migrationRegistry = make(map[string]Migration)

func RegisterMigration(migration string, up, down func(*context.Context) error) {
	migrationRegistry[migration] = Migration{Up: up, Down: down}
}

func RunMigrations(ctx *context.Context) (error error) {
	migrations := migrationRegistry

	checkSchemaVersion(ctx)

	var keys []string
	for k := range migrations {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, migration := range keys {

		if checkMigration(ctx, migration) == true {
			continue
		}

		slog.Info("Running migration: " + migration)
		success := false
		if error = migrations[migration].Up(ctx); error != nil {
			slog.Warn("Failed migration: " + migration)
			if error = migrations[migration].Down(ctx); error != nil {
				slog.Warn("Failed migration rollback: " + migration)
			}
		} else {
			success = true
		}

		_, error = db.Pool.Exec(*ctx, `INSERT INTO schema_version (name, success) VALUES ($1, $2)`, migration, success)
		if error != nil {
			slog.Warn("Failed schema_version update: " + migration)
		}
	}

	return error
}

func checkSchemaVersion(ctx *context.Context) {
	_, err := db.Pool.Exec(*ctx, `CREATE TABLE IF NOT EXISTS schema_version (
		id serial4 NOT NULL,
		"timestamp" timestamptz DEFAULT now() NOT NULL,
		"name" varchar NOT NULL,
		success bool NOT NULL,
		CONSTRAINT schema_version_pk PRIMARY KEY (id),
		CONSTRAINT schema_version_unique UNIQUE (name)
	);`)

	if err != nil {
		log.Fatalf("Error initialization migrations")
	}
}

func checkMigration(ctx *context.Context, migration string) (bool bool) {

	result := false

	rows, error := db.Pool.Query(*ctx, `SELECT * FROM schema_version WHERE name = $1`, migration)
	if error != nil {
		slog.Error(error.Error())
		return result
	}
	defer rows.Close()
	schema, error := pgx.CollectRows(rows, pgx.RowToStructByName[SchemaVersion])
	if error != nil {
		slog.Error(error.Error())
	} else {
		if len(schema) > 0 && schema[0].Success == true {
			result = true
		}
	}
	return result
}
