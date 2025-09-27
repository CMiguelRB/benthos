# benthos
From Ancient Greek βένθος (bénthos) 'the depths (of the sea)', is the community of organisms that live on, in, or near the bottom of a sea, river, lake, or stream, also known as the benthic zone.

benthos is a backend service written in go, meant to be used as a starting point for any REST API web server project.

## Environment variables

Environment variables must either exists in the os environment variables or in a .env file. This is an .env file example. 

```
.env

ENV=DEV
DB_USER="TEST"
DB_PASSWORD="TEST"
DB_HOSTNAME="localhost"
DB_PORT="5432"
DB_NAME="TEST"
ENCRYPTION_KEY="32 bit HEX encryption key" ## Used for credentials encryption
```

## Configuration

Service configuration can be found inside the `config` package. A `Config` struct is defined and globally exporeted as `Settings`. The values are set in the `InitConfiguration` method. This way of setting the global service configuration has two goals: to avoid depending on external files for variables that rarely change, other than env variables, and to provide a structured way to access these variables from the global context.

```
config/config.go

package config

import (
	"os"
	"time"
	"sync"
)

type Config struct {
	App    App
	Server Server
	...
}

type App struct {
	Name    string
	Version string
	...
}

type Server struct {
	ReadTimeoutMs  time.Duration
	WriteTimeoutMs time.Duration
	IdleTimeoutMs  time.Duration
	RateLimit      RateLimit
	...
}

type RateLimit struct {
	Requests int
	PeriodMs int
	...
}

type Database struct {
	Hostname string
	...
}

type Security struct {
	EncryptionKey string
	...
}

var (
	once     sync.Once
	Settings Config
)

//This global variable is updated with the provided git tag at release.yml Github Action runtime
var Version = "version"

func InitConfiguration() {
	once.Do(func() {
		//App
		Settings.App.Name = "benthos"
		Settings.App.Version = "v0.0.4"
		...
		//Server
		Settings.Server.ReadTimeoutMs = 15000
		Settings.Server.WriteTimeoutMs = 15000
		Settings.Server.IdleTimeoutMs = 60000
		...
		//Server RateLimit
		Settings.Server.RateLimit.Requests = 10
		Settings.Server.RateLimit.PeriodMs = 10000
		...
		//DB
		Settings.Database.Hostname = os.Getenv("DB_HOSTNAME")
		...
		//Security
		Settings.Security.EncryptionKey = os.Getenv("ENCRYPTION_KEY")
		...
	})
}
```

Also, in order to get secrets, they can be either stored in a .env file, but for a production environment they may have to be loaded from the secrets folder (Docker):

```
...
func InitConfiguration() {
	...
	if os.Getenv("ENV") == "PROD" {
		Settings.App.Version = Version
		...
		//DB
		dbPassword, err := loadSecret("db_password")
		if err != nil {
			log.Fatal("DB Password secret not found")
		}
		Settings.Database.Password = dbPassword;
		//Security
		encryptionKey, err := loadSecret("benthos_encryption_key")
		if err != nil {
			log.Fatal("Encryption key secret not found")
		}
		Settings.Security.EncryptionKey = encryptionKey
	}
}

func loadSecret(name string) (string, error) {
	secretPath := filepath.Join("/run/secrets", name) //Load the secret from the default path Docker loads the secrets configured in the docker-compose.yml file.

	if data, err := os.ReadFile(secretPath); err == nil {
		return strings.TrimSpace(string(data)), nil
	}

	return "", os.ErrNotExist
}
```

To access the configuration from other packages:

```
server/server.go

package server

import (
	"benthos/config"
	"fmt"
	...
)

func New(ctx *context.Context) *http.Server {
	...
	fmt.Println("App version: " + config.Settings.App.Version)
	//Output: App version: v0.0.3
	...
}

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

func NewDomnameService(repo DomnameRepo) *DomnameService {
	return &DomnameService{repo: repo}
}
```
```
module.go

type Module struct {
	Repo    *DomnameRepo
	Service *app.DomnameService
	Routes  *DomnameRoutes
}

func NewModule() Module {
	slog.Info("Loading Domname module...")
	
	repo := NewDomnameRepo()
	service := app.NewDomnameService(repo)
	routes := NewDomnameRoutes(service)

	return Module{
		Repo:    repo,
		Service: service,
		Routes:  routes,
	}
}
```

```
routes.go

type Routes struct {
	handler *Handler
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
slog.Info("Loading domain modules...")
domname.NewModule().Routes.Configure(mux) // -> Existing module
newmodule.NewModule.Routes.Configure(mux) // -> your new module
```

## Run

```
go run main.go
```

## Test

```
go test ./...
```

## Build

```
go build -v -ldflags "-X benthos/config.Version=v0.0.0"

v0.0.0 -> your app version. The GH action does this for you at build stage based on the pushed tag.
```
