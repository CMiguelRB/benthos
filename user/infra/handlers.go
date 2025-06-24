package infra

import (
	"benthos/common/res"
	"benthos/user/app"
	"benthos/user/dom"
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *app.Service
	validator *UserValidator
}

var contentType string = "Content-Type"
var applicationJson string = "application/json"

func NewHandler(service *app.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentType, applicationJson)

	result := h.service.GetUsers(context.Background())

	if !result.Success {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(result)
}

func (h *Handler) getById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentType, applicationJson)
	id := chi.URLParam(r, "id")

	result := h.service.GetUserById(context.Background(), id)

	if !result.Success {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(result)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentType, applicationJson)

	user := dom.User{}

	if validation := validateBody(h, r, &user); len(validation.Error) > 0{
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validation)
		return
	}

	result := h.service.CreateUser(context.Background(), user)

	if !result.Success {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(result)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentType, applicationJson)
	id := chi.URLParam(r, "id")

	user := dom.User{}

	if validation := validateBody(h, r, &user); len(validation.Error) > 0{
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validation)
		return
	}

	result := h.service.UpdateUser(context.Background(), id, user)

	if !result.Success {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(result)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentType, applicationJson)
	id := chi.URLParam(r, "id")

	result := h.service.DeleteUser(context.Background(), id)

	if !result.Success {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(result)
}

func validateBody(h *Handler, r *http.Request, user *dom.User) (result common.ErrorResponse){

	
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		result = common.ErrorResponse{
			Error:   []common.Error{{Code: "EUSRG1", Message: "Invalid JSON"}},
			Success: false,
		}
		return result
	}

	if validationErrors := h.validator.ValidateUser(user); len(validationErrors) > 0 {
		result = common.ErrorResponse{
			Error:   validationErrors,
			Success: false,
		}
		return result
    }

	return
}
