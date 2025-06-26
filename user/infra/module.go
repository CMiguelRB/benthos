package infra

import (
	shared "benthos/shared/infra"
	"benthos/user/app"
	"log/slog"
)

func NewModule() shared.Module[*UserRepo, *app.UserService, *UserRoutes] {
	slog.Info("Loading User module...")
	return shared.Module[*UserRepo, *app.UserService, *UserRoutes]{
		NewRepo: NewUserRepo,
		NewService: func(r *UserRepo) *app.UserService {
			return app.NewUserService(r)
		},
		NewRoutes: NewRoutes,
	}
}
