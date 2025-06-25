package migrations

import (
	"benthos/db"
	"context"
)

func init() {
	RegisterMigration("20250625init", up, down)
}

func up(ctx *context.Context) (err error) {

	query := `CREATE TABLE IF NOT EXISTS users (
		id uuid DEFAULT uuid_generate_v4() NOT NULL,
		username varchar NOT NULL,
		"password" varchar NOT NULL,
		created_on timestamptz DEFAULT now() NOT NULL,
		updated_on timestamptz NULL,
		last_access timestamptz NULL,
		CONSTRAINT users_pk PRIMARY KEY (id),
		CONSTRAINT users_username_uq UNIQUE (username)
	);`

	_, err = db.Pool.Exec(*ctx, query)

	return err
}

func down(ctx *context.Context) (err error) {

	query := `DROP TABLE IF EXISTS public.users ;`

	_, err = db.Pool.Exec(*ctx, query)

	return err
}
