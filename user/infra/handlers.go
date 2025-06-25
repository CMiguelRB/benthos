package infra

import (
	commonDom "benthos/common/dom"
	"benthos/user/app"
	"benthos/user/dom"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service   *app.UserService
	validator *UserValidator
}

var contentType string = "Content-Type"
var applicationJson string = "application/json"

func NewHandler(service *app.UserService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentType, applicationJson)

	result := h.service.GetUsers(r.Context())

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

	result := h.service.GetUserById(r.Context(), id)

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

	if validation := validateBody(h, r, &user); len(validation.Error) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validation)
		return
	}

	result := h.service.CreateUser(r.Context(), user)

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

	if validation := validateBody(h, r, &user); len(validation.Error) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validation)
		return
	}

	result := h.service.UpdateUser(r.Context(), id, user)

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

	result := h.service.DeleteUser(r.Context(), id)

	if !result.Success {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(result)
}

func validateBody(h *Handler, r *http.Request, user *dom.User) (result commonDom.ErrorResponse) {

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		result = commonDom.ErrorResponse{
			Error:   []commonDom.Error{{Code: "EUSRG1", Message: "Invalid JSON"}},
			Success: false,
		}
		return result
	}

	if validationErrors := h.validator.ValidateUser(user); len(validationErrors) > 0 {
		result = commonDom.ErrorResponse{
			Error:   validationErrors,
			Success: false,
		}
		return result
	}

	return
}
