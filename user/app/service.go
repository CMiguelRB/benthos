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

func (s *UserService) GetUsers(ctx context.Context) (result shared.QueryResult[users.User]) {

	users, err := s.repo.GetUsers(ctx)

	return createQueryResult(users, err)
}

func (s *UserService) GetUserById(ctx context.Context, id string) (result shared.QueryResult[users.User]) {

	user, err := s.repo.GetUserById(ctx, id)

	return createQueryResult(user, err)
}

func (s *UserService) CreateUser(ctx context.Context, user users.User) (result shared.PersistenceResult) {

	id, err := s.repo.CreateUser(ctx, user)

	return createPersistenceResult(&id, nil, err)
}

func (s *UserService) UpdateUser(ctx context.Context, id string, user users.User) (result shared.PersistenceResult) {
	rows, err := s.repo.UpdateUser(ctx, id, user)

	return createPersistenceResult(nil, &rows, err)
}

func (s *UserService) DeleteUser(ctx context.Context, id string) (result shared.PersistenceResult) {
	rows, err := s.repo.DeleteUser(ctx, id)

	return createPersistenceResult(nil, &rows, err)
}

func createQueryResult(user []users.User, err error) (result shared.QueryResult[users.User]) {

	result = shared.QueryResult[users.User]{
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

func createPersistenceResult(id *string, rows *int64, err error) (result shared.PersistenceResult) {

	result = shared.PersistenceResult{
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
