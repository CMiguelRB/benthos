package app

import (
	"benthos_go/common"
	"benthos_go/user/dom"
	"context"
)

type Service struct {
	repo UserRepo
}

func NewService(repo UserRepo) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetUsers(ctx context.Context) (result common.Result[dom.User]) {

	users, err := s.repo.GetUsers(ctx)

	result = common.Result[dom.User]{
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

func (s *Service) GetUserById(ctx context.Context, id string) (result common.Result[dom.User]) {

	user, err := s.repo.GetUserById(ctx, id)

	result = common.Result[dom.User]{
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


func (s *Service) CreateUser(ctx context.Context, user dom.User) (result common.Result[common.ActionResult]) {

	res, err := s.repo.CreateUser(ctx, user)

	result = common.Result[common.ActionResult]{
		Success: false,
	}

	if err != nil {
		result.Error = err.Error()
	} else {
		result.Data =  []common.ActionResult{{AffectedRows: res}}
		if res == 0 {
			result.Success = false
		}
		result.Success = true
	}

	return result
}

func (s *Service) UpdateUser(ctx context.Context, id string, user dom.User) (result common.Result[common.ActionResult]){
	res, err := s.repo.UpdateUser(ctx, id, user)

	result = common.Result[common.ActionResult]{
		Success: false,
	}

	if err != nil {
		result.Error = err.Error()
	} else {
		result.Data =  []common.ActionResult{{AffectedRows: res}}
		if res == 0 {
			result.Success = false
		}
		result.Success = true
	}

	return result
}

func (s *Service) DeleteUser(ctx context.Context, id string) (result common.Result[common.ActionResult]){
	res, err := s.repo.DeleteUser(ctx, id)

	result = common.Result[common.ActionResult]{
		Success: false,
	}

	if err != nil {
		result.Error = err.Error()
	} else {
		result.Data =  []common.ActionResult{{AffectedRows: res}}
		if res == 0 {
			result.Success = false
		}
		result.Success = true
	}

	return result
}
