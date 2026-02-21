package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"taskmanager/internal/service/userservice"
	"taskmanager/internal/utils"
	"taskmanager/pkg/modules"
)

type UserHandler struct {
	service userservice.UserUseCaseInterface
}

func NewUserHandler(service *userservice.UserUseCase) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAllUsers()

	if err != nil {
		handleError(w, err)
		return
	}

	if err := encode(w, r, http.StatusOK, users); err != nil {
		handleError(w, err)
		return
	}
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		handleError(w, fmt.Errorf("invalid id: %w", utils.ErrInvalidData))
		return
	}

	user, err := h.service.GetUserByID(id)

	if err != nil {
		handleError(w, err)
		return
	}

	if err := encode(w, r, http.StatusOK, user); err != nil {
		handleError(w, err)
		return
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	user, err := decode[modules.User](r)

	if err != nil {
		handleError(w, err)
		return
	}

	id, err := h.service.CreateUser(user)

	if err != nil {
		handleError(w, err)
		return
	}

	err = encode(w, r, http.StatusOK, map[string]int{"id": id})

	if err != nil {
		handleError(w, err)
	}

}

func (h *UserHandler) UpdateUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		handleError(w, fmt.Errorf("invalid id: %w", utils.ErrInvalidData))
		return
	}

	user, err := decode[userservice.UpdateUserInput](r)

	if err != nil {
		handleError(w, err)
		return
	}

	if err := h.service.UpdateUserByID(id, user); err != nil {
		handleError(w, err)
		return
	}

	encode(w, r, http.StatusOK, map[string]string{"message": "user updated"})
}

func (h *UserHandler) DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		handleError(w, err)
		return
	}

	if err := h.service.DeleteUserByID(id); err != nil {
		handleError(w, err)
		return
	}

	encode(w, r, http.StatusOK, map[string]string{"message": "user deleted"})
}

//**********************
// Additional functions

func encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}

	return nil
}

func decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}

	return v, nil
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, utils.ErrObjectNotFound):
		encode(w, nil, http.StatusNotFound, map[string]string{"error": err.Error()})
	case errors.Is(err, utils.ErrInvalidData):
		encode(w, nil, http.StatusBadRequest, map[string]string{"error": err.Error()})
	default:
		encode(w, nil, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

}
