package infra

import (
	"benthos/user/app"

	"github.com/go-chi/chi/v5"
)

const (
    userPath   = "/api/users"
    userIdPath = userPath + "/{id}"
)

type UserRoutes struct {
	handler *Handler
}

func NewUserRoutes(service *app.UserService) *UserRoutes {
	return &UserRoutes{
		handler: NewHandler(service),
	}
}

func (r *UserRoutes) Configure(mux *chi.Mux) {
	mux.Get(userPath, r.handler.list)
	mux.Get(userIdPath, r.handler.getById)
	mux.Post(userPath, r.handler.create)
	mux.Put(userIdPath, r.handler.update)
	mux.Delete(userIdPath, r.handler.delete)
}
