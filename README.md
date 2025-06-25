# benthos
From Ancient Greek βένθος (bénthos) 'the depths (of the sea)', is the community of organisms that live on, in, or near the bottom of a sea, river, lake, or stream, also known as the benthic zone.

benthos is a backend service written in go, meant to be used as a starting point for any REST API web server project.

## Environment variables

Environment variables must either exists in the os environment variables or in a .env file. This is an .env file example. 

```
.env

DB_USER="TEST"
DB_PASSWORD="TEST"
DB_HOSTNAME="localhost"
DB_PORT="5432"
DB_NAME="TEST"
ENCRYPTION_KEY="32 bit HEX encryption key" ## Used for credentials encryption
```

## Database

benthos uses PostreSQL database. The database must be manually created and the name must be configured in the DB_NAME env variable.

#### Migrations
Migrations are defined by file under `db/migrations/*`. Migration files are executed automatically at the start of the application. They must follow the following structure:
```
examplemigration.go

package migrations

import (
	"benthos/db"
	"context"
)

//init function is executed when instantiated
func init() {
	RegisterMigration("migrationName", up, down) // migrationName must be unique!! ie: 20250625102030migrationReason
}

func up(ctx *context.Context) (err error) {
	//Database operations here
}

func down(ctx *context.Context) (err error) {
  //Database operations rollback here
}
```

## Run

```
go run main.go
```

## Test
```
go text ./...
```

## Build

```
go build
```
