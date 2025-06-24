package app

import (
	common "benthos/common/res"
	"benthos/user/dom"
	"context"
)

type Service struct {
	repo UserRepo
}

func NewService(repo UserRepo) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetUsers(ctx context.Context) (result common.QResult[dom.User]) {

	users, err := s.repo.GetUsers(ctx)

	result = common.QResult[dom.User]{
		Success: false,
	}

	if err != nil {
		result.Error = err.Error()
	} else {
		result.Data = users
		result.Success = true
	}

	return result
}

func (s *Service) GetUserById(ctx context.Context, id string) (result common.QResult[dom.User]) {

	user, err := s.repo.GetUserById(ctx, id)

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

func (s *Service) CreateUser(ctx context.Context, user dom.User) (result common.WResult) {

	id, err := s.repo.CreateUser(ctx, user)

	return createWResult(&id, nil, err)
}

func (s *Service) UpdateUser(ctx context.Context, id string, user dom.User) (result common.WResult) {
	rows, err := s.repo.UpdateUser(ctx, id, user)

	return createWResult(nil, &rows, err)
}

func (s *Service) DeleteUser(ctx context.Context, id string) (result common.WResult) {
	rows, err := s.repo.DeleteUser(ctx, id)

	return createWResult(nil, &rows, err)
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
