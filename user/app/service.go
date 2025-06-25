package app

import (
	commonDom "benthos/common/dom"
	userDom "benthos/user/dom"
	"context"
)

type UserService struct {
	repo UserRepo
}

func NewUserService(repo UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUsers(ctx context.Context) (result commonDom.QResult[userDom.User]) {

	users, err := s.repo.GetUsers(ctx)

	return createQResult(users, err)
}

func (s *UserService) GetUserById(ctx context.Context, id string) (result commonDom.QResult[userDom.User]) {

	user, err := s.repo.GetUserById(ctx, id)

	return createQResult(user, err)
}

func (s *UserService) CreateUser(ctx context.Context, user userDom.User) (result commonDom.WResult) {

	id, err := s.repo.CreateUser(ctx, user)

	return createWResult(&id, nil, err)
}

func (s *UserService) UpdateUser(ctx context.Context, id string, user userDom.User) (result commonDom.WResult) {
	rows, err := s.repo.UpdateUser(ctx, id, user)

	return createWResult(nil, &rows, err)
}

func (s *UserService) DeleteUser(ctx context.Context, id string) (result commonDom.WResult) {
	rows, err := s.repo.DeleteUser(ctx, id)

	return createWResult(nil, &rows, err)
}

func createQResult(user []userDom.User, err error) (result commonDom.QResult[userDom.User]) {

	result = commonDom.QResult[userDom.User]{
		Success: false,
	}

	if err != nil {
		result.Error = err.Error()
	} else {
		if user == nil {
			result.Data = []userDom.User{}
		} else {
			result.Data = user
		}
		result.Success = true
	}

	return result
}

func createWResult(id *string, rows *int64, err error) (result commonDom.WResult) {

	result = commonDom.WResult{
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
