# benthos
From Ancient Greek βένθος (bénthos) 'the depths (of the sea)', is the community of organisms that live on, in, or near the bottom of a sea, river, lake, or stream, also known as the benthic zone.

benthos is a backend service written in go, meant to be used as a starting point for any REST API web server project.

## Environment variables

Environment variables must either exists in the os environment variables or in a .env file. This is an .env file example. 

```
.env

ENV=DEV
NAME=benthos
VERSION=0.0.1   ## Only for local execution, the GH Action sets the version based on the pushed tag.
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
## Domain modules

benthos was architectured following clean architecture and DDD paradigms. Domain modules keep the model, services and infra logic of each domain:

`domname is the domain name: user, order, animal, pirola...`

```
/domname
  /app
    service.go
  /dom
    ports.go        // port definition for the interface between app (service.go) and infra (repo.go)
    domname.go      // model of the current domain
  /infra
    handlers.go     // API controller
    module.go       // domain module loader
    repo.go         // repository
    routes.go       // route definition
    validation.go   // data validation
```

For the auto module loading to work, some files must implement some common logic:
```
service.go

type DomnameService struct {
	repo userDom.DomnameRepo
}

func NewDomnamerService(repo DomnameRepo) *DomnameService {
	return &DomnameService{repo: repo}
}
```
```
module.go

func NewModule() common.Module[*DomnameRepo, *app.DomnameService, *DomnameRoutes] { // Check common Module abstract implementation
	return common.Module[*DomnameRepo, *app.DomnameService, *DomnameRoutes]{
		NewRepo: NewDomnameRepo,
		NewService: func(r *DomnameRepo) *app.DomnameService {
			return app.NewDomnameService(r)
		},
		NewRoutes: NewRoutes,
	}
}
```

```
routes.go

type Routes struct {
	handler *Handler
	ctx     *context.Context
}

func NewRoutes(service *app.DomnameService) *DomnameRoutes {
	return &DomnameRoutes{
		handler: NewHandler(service),
	}
}

func (r *DomnameRoutes) Configure(mux *chi.Mux) {
	// routes configuration
}
```

Once the new module is set up, the only required step is to add it to he server modules:

```
modules := []commonInfra.ModuleInitializer{
	existingDomnameInfra.NewModule(), // -> Existing module
	domnameInfra.NewModule(),         // -> your new module
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
