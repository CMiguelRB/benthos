package infra

import (
	app "benthos/user/app"
	"log/slog"
)

type Module struct {
	Repo    *UserRepo
	Service *app.UserService
	Routes  *UserRoutes
}

func NewModule() Module {
	slog.Info("Loading Users module...")
	
	repo := NewUserRepo()
	service := app.NewUserService(repo)
	routes := NewUserRoutes(service)

	return Module{
		Repo:    repo,
		Service: service,
		Routes:  routes,
	}
}
