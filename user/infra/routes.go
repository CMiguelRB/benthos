package infra

import (
	"benthos_go/user/app"
	"context"

	"github.com/go-chi/chi/v5"
)

var userPath string = "/user"
var userIdPath = userPath + "/{id}"

type UserRoutes struct {
    handler *Handler
}

func NewUserRoutes(service *app.Service) *UserRoutes {
    return &UserRoutes{
        handler: NewHandler(service),
    }
}

func (r *UserRoutes) Configure(mux *chi.Mux, ctx *context.Context) {
    mux.Get(userPath, r.handler.list)
    mux.Get(userIdPath, r.handler.getById)
    mux.Post(userPath, r.handler.create)
    mux.Post(userIdPath, r.handler.update)
    mux.Delete(userIdPath, r.handler.delete)
}