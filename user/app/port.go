package app

import (
	"benthos/user/dom"
	"context"
)

type UserRepo interface {
	GetUsers(ctx context.Context) ([]dom.User, error)
	GetUserById(ctx context.Context, id string) ([]dom.User, error)
	CreateUser(ctx context.Context, user dom.User) (string, error)
	UpdateUser(ctx context.Context, id string, user dom.User) (int64, error)
	DeleteUser(ctx context.Context, id string) (int64, error)
}
