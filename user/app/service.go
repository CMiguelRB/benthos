package app

import (
	"benthos/common/res"
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

	res, err := s.repo.CreateUser(ctx, user)

	return createWResult(res, err)
}

func (s *Service) UpdateUser(ctx context.Context, id string, user dom.User) (result common.WResult) {
	res, err := s.repo.UpdateUser(ctx, id, user)

	return createWResult(res, err)
}

func (s *Service) DeleteUser(ctx context.Context, id string) (result common.WResult) {
	res, err := s.repo.DeleteUser(ctx, id)

	return createWResult(res, err)
}

func createWResult(res int64, err error) (result common.WResult) {

	result = common.WResult{
		Success: false,
	}

	if err != nil {
		result.Error = err.Error()
	} else {
		result.AffectedRows = res
		if res == 0 {
			result.Success = false
		}
		result.Success = true
	}

	return result
}
