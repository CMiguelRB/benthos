package app

import (
	common "benthos/common/res"
	"benthos/user/dom"
	"context"
)

type UserService struct {
	repo UserRepo
}

func NewUserService(repo UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUsers(ctx context.Context) (result common.QResult[dom.User]) {

	users, err := s.repo.GetUsers(ctx)

	return createQResult(users, err)
}

func (s *UserService) GetUserById(ctx context.Context, id string) (result common.QResult[dom.User]) {

	user, err := s.repo.GetUserById(ctx, id)

	return createQResult(user, err)
}

func (s *UserService) CreateUser(ctx context.Context, user dom.User) (result common.WResult) {

	id, err := s.repo.CreateUser(ctx, user)

	return createWResult(&id, nil, err)
}

func (s *UserService) UpdateUser(ctx context.Context, id string, user dom.User) (result common.WResult) {
	rows, err := s.repo.UpdateUser(ctx, id, user)

	return createWResult(nil, &rows, err)
}

func (s *UserService) DeleteUser(ctx context.Context, id string) (result common.WResult) {
	rows, err := s.repo.DeleteUser(ctx, id)

	return createWResult(nil, &rows, err)
}

func createQResult(user []dom.User, err error) (result common.QResult[dom.User]) {

	result = common.QResult[dom.User]{
		Success: false,
	}

	if err != nil {
		result.Error = err.Error()
	} else {
		if user == nil {
			result.Data = []dom.User{}
		} else {
			result.Data = user
		}
		result.Success = true
	}

	return result
}

func createWResult(id *string, rows *int64, err error) (result common.WResult) {

	result = common.WResult{
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
