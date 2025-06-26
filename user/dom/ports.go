package dom

import (
	"context"
)

type UserRepo interface {
	GetUsers(ctx context.Context) ([]User, error)
	GetUserById(ctx context.Context, id string) ([]User, error)
	CreateUser(ctx context.Context, user User) (string, error)
	UpdateUser(ctx context.Context, id string, user User) (int64, error)
	DeleteUser(ctx context.Context, id string) (int64, error)
}
