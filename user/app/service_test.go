package app

import (	
	users "benthos/user/dom"
	"context"
	"errors"
	"testing"
)

type UserRepoStub struct {
	GetUsersFunc     func(ctx context.Context) ([]users.User, error)
	GetUserByIdFunc  func(ctx context.Context, id string) ([]users.User, error)
	CreateUserFunc   func(ctx context.Context, user users.User) (string, error)
	UpdateUserFunc   func(ctx context.Context, id string, user users.User) (int64, error)
	DeleteUserFunc   func(ctx context.Context, id string) (int64, error)
}

func (r *UserRepoStub) GetUsers(ctx context.Context) ([]users.User, error) {
	return r.GetUsersFunc(ctx)
}

func (r *UserRepoStub) GetUserById(ctx context.Context, id string) ([]users.User, error) {
	return r.GetUserByIdFunc(ctx, id)
}

func (r *UserRepoStub) CreateUser(ctx context.Context, user users.User) (string, error) {
	return r.CreateUserFunc(ctx, user)
}

func (r *UserRepoStub) UpdateUser(ctx context.Context, id string, user users.User) (int64, error) {
	return r.UpdateUserFunc(ctx, id, user)
}

func (r *UserRepoStub) DeleteUser(ctx context.Context, id string) (int64, error) {
	return r.DeleteUserFunc(ctx, id)
}

func TestGetUsersSuccess(t *testing.T) {
	repo := &UserRepoStub{
		GetUsersFunc: func(ctx context.Context) ([]users.User, error) {
			return []users.User{{Id: "123", Username: "Alice"}}, nil
		},
	}

	service := NewUserService(repo)
	result := service.GetUsers(context.Background())

	if !result.Success {
		t.Errorf("expected Success=true, got false: %+v", result)
	}
	if len(result.Data) != 1 || result.Data[0].Id != "123" {
		t.Errorf("unexpected data: %+v", result.Data)
	}
}

func TestGetUserByIdError(t *testing.T) {
	repo := &UserRepoStub{
		GetUserByIdFunc: func(ctx context.Context, id string) ([]users.User, error) {
			return nil, errors.New("user not found")
		},
	}

	service := NewUserService(repo)
	result := service.GetUserById(context.Background(), "invalid-id")

	if result.Success {
		t.Errorf("expected Success=false, got true")
	}
	if result.Error != "user not found" {
		t.Errorf("unexpected error: %s", result.Error)
	}
}

func TestCreateUserSuccess(t *testing.T) {
	repo := &UserRepoStub{
		CreateUserFunc: func(ctx context.Context, user users.User) (string, error) {
			return "new-id", nil
		},
	}

	service := NewUserService(repo)
	result := service.CreateUser(context.Background(), users.User{Username: "Bob"})

	if !result.Success || result.Id == nil || *result.Id != "new-id" {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestUpdateUserSuccess(t *testing.T) {
	repo := &UserRepoStub{
		UpdateUserFunc: func(ctx context.Context, id string, user users.User) (int64, error) {
			return 1, nil
		},
	}

	service := NewUserService(repo)
	result := service.UpdateUser(context.Background(), "123", users.User{Username: "Bob"})

	if !result.Success || result.Rows == nil || *result.Rows != 1 {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestDeleteUserFailure(t *testing.T) {
	repo := &UserRepoStub{
		DeleteUserFunc: func(ctx context.Context, id string) (int64, error) {
			return 0, errors.New("delete failed")
		},
	}

	service := NewUserService(repo)
	result := service.DeleteUser(context.Background(), "123")

	if result.Success {
		t.Errorf("expected Success=false, got true")
	}
	if result.Error != "delete failed" {
		t.Errorf("unexpected error: %s", result.Error)
	}
}