package app

import (
	shared "benthos/shared/dom"
	users "benthos/user/dom"
	"context"
)

type UserService struct {
	repo users.UserRepo
}

func NewUserService(repo users.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUsers(ctx context.Context) (result shared.QResult[users.User]) {

	users, err := s.repo.GetUsers(ctx)

	return createQResult(users, err)
}

func (s *UserService) GetUserById(ctx context.Context, id string) (result shared.QResult[users.User]) {

	user, err := s.repo.GetUserById(ctx, id)

	return createQResult(user, err)
}

func (s *UserService) CreateUser(ctx context.Context, user users.User) (result shared.WResult) {

	id, err := s.repo.CreateUser(ctx, user)

	return createWResult(&id, nil, err)
}

func (s *UserService) UpdateUser(ctx context.Context, id string, user users.User) (result shared.WResult) {
	rows, err := s.repo.UpdateUser(ctx, id, user)

	return createWResult(nil, &rows, err)
}

func (s *UserService) DeleteUser(ctx context.Context, id string) (result shared.WResult) {
	rows, err := s.repo.DeleteUser(ctx, id)

	return createWResult(nil, &rows, err)
}

func createQResult(user []users.User, err error) (result shared.QResult[users.User]) {

	result = shared.QResult[users.User]{
		Success: false,
	}

	if err != nil {
		result.Error = err.Error()
	} else {
		if user == nil {
			result.Data = []users.User{}
		} else {
			result.Data = user
		}
		result.Success = true
	}

	return result
}

func createWResult(id *string, rows *int64, err error) (result shared.WResult) {

	result = shared.WResult{
		Success: false,
	}

	if err != nil {
		result.Error = err.Error()
	} else {
		if result.Id = id; id != nil {
			result.Success = true
		}
		if result.Rows = rows; rows != nil {
			result.Success = true
		}
	}

	return result
}
