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

		exists, success := checkMigration(ctx, migration)

		if exists == true && success == true {
			continue
		}

		slog.Info("Running migration: " + migration)
		if error = migrations[migration].Up(ctx); error != nil {
			slog.Warn("Failed migration: " + migration)
			if error = migrations[migration].Down(ctx); error != nil {
				slog.Warn("Failed migration rollback: " + migration)
			}
			updateSchemaVersion(ctx, migration, exists, false)
		} else {
			updateSchemaVersion(ctx, migration, exists, true)
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
		db.Pool.Close()
		log.Fatalf("Error initialization migrations")
	}
}

func checkMigration(ctx *context.Context, migration string) (exists bool, success bool) {

	exists = false
	success = false

	rows, error := db.Pool.Query(*ctx, `SELECT * FROM schema_version WHERE name = $1`, migration)
	if error != nil {
		slog.Error(error.Error())
		return exists, success
	}
	defer rows.Close()
	schema, error := pgx.CollectRows(rows, pgx.RowToStructByName[SchemaVersion])
	if error != nil {
		slog.Error(error.Error())
	} else {
		rNum := len(schema)
		if rNum > 0 {
			exists = true
		}
		if rNum > 0 && schema[0].Success {
			success = true
		}
	}
	return exists, success
}

func updateSchemaVersion(ctx *context.Context, migration string, exists bool, migrationSuccess bool) {

	var error error
	if exists == true {
		_, error = db.Pool.Exec(*ctx, `UPDATE schema_version SET success=$1`, migrationSuccess)
	} else {
		_, error = db.Pool.Exec(*ctx, `INSERT INTO schema_version (name, success) VALUES ($1, $2)`, migration, migrationSuccess)
	}

	if error != nil {
		slog.Warn("Failed schema_version update: " + migration)
	}
}
