package infra

import (
	common "benthos/common/infra"
	"benthos/user/app"
)

func NewModule() common.Module[*UserRepo, *app.UserService, *UserRoutes] {
	return common.Module[*UserRepo, *app.UserService, *UserRoutes]{
		NewRepo: NewUserRepo,
		NewService: func(r *UserRepo) *app.UserService {
			return app.NewUserService(r) 
		},
		NewRoutes: NewRoutes, 
	}
}