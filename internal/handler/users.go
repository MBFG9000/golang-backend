package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/MBFG9000/golang-backend/internal/domain"
	userservice "github.com/MBFG9000/golang-backend/internal/service/users"
	"github.com/MBFG9000/golang-backend/internal/utils"

	"github.com/google/uuid"
)

type UserHandler struct {
	service userservice.UserServiceInterface
}

func NewUserHandler(service userservice.UserServiceInterface) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.URL.RawQuery != "" {
		h.GetPaginatedUsers(w, r)
		return
	}

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
	id, err := uuid.Parse(idStr)

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
	user, err := decode[domain.User](r)

	if err != nil {
		handleError(w, err)
		return
	}

	id, err := h.service.CreateUser(user)

	if err != nil {
		handleError(w, err)
		return
	}

	err = encode(w, r, http.StatusOK, map[string]string{"id": id.String()})

	if err != nil {
		handleError(w, err)
	}

}

func (h *UserHandler) UpdateUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)

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
	id, err := uuid.Parse(idStr)

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

func (h *UserHandler) GetPaginatedUsers(w http.ResponseWriter, r *http.Request) {
	page, err := parsePositiveIntParam(r, "page", 1)
	if err != nil {
		handleError(w, err)
		return
	}

	pageSize, err := parsePositiveIntParam(r, "page_size", 10)
	if err != nil {
		handleError(w, err)
		return
	}

	filter, err := parseUserFilter(r)
	if err != nil {
		handleError(w, err)
		return
	}

	response, err := h.service.GetPaginatedUsers(page, pageSize, filter)
	if err != nil {
		handleError(w, err)
		return
	}

	if err := encode(w, r, http.StatusOK, response); err != nil {
		handleError(w, err)
	}
}

func (h *UserHandler) GetCommonFriends(w http.ResponseWriter, r *http.Request) {
	userID1, err := parseUUIDQueryParam(r, "user_id_1")
	if err != nil {
		handleError(w, err)
		return
	}

	userID2, err := parseUUIDQueryParam(r, "user_id_2")
	if err != nil {
		handleError(w, err)
		return
	}

	users, err := h.service.GetCommonFriends(userID1, userID2)
	if err != nil {
		handleError(w, err)
		return
	}

	if err := encode(w, r, http.StatusOK, users); err != nil {
		handleError(w, err)
	}
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

func parsePositiveIntParam(r *http.Request, key string, defaultValue int) (int, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue, nil
	}

	num, err := strconv.Atoi(value)
	if err != nil || num <= 0 {
		return 0, fmt.Errorf("invalid %s: %w", key, utils.ErrInvalidData)
	}

	return num, nil
}

func parseUUIDQueryParam(r *http.Request, key string) (uuid.UUID, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return uuid.Nil, fmt.Errorf("%s is required: %w", key, utils.ErrInvalidData)
	}

	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid %s: %w", key, utils.ErrInvalidData)
	}

	return id, nil
}

func parseUserFilter(r *http.Request) (domain.UserFilter, error) {
	query := r.URL.Query()
	filter := domain.UserFilter{
		SortBy:    query.Get("sort_by"),
		SortOrder: query.Get("sort_order"),
	}

	if idRaw := query.Get("id"); idRaw != "" {
		id, err := uuid.Parse(idRaw)
		if err != nil {
			return domain.UserFilter{}, fmt.Errorf("invalid filter id: %w", utils.ErrInvalidData)
		}
		filter.ID = &id
	}

	if firstName := query.Get("first_name"); firstName != "" {
		filter.FirstName = &firstName
	}

	if email := query.Get("email"); email != "" {
		filter.Email = &email
	}

	if gender := query.Get("gender"); gender != "" {
		filter.Gender = &gender
	}

	if birthDateRaw := query.Get("birth_date"); birthDateRaw != "" {
		birthDate, err := time.Parse("2006-01-02", birthDateRaw)
		if err != nil {
			return domain.UserFilter{}, fmt.Errorf("invalid birth_date, expected YYYY-MM-DD: %w", utils.ErrInvalidData)
		}
		filter.BirthDate = &birthDate
	}

	return filter, nil
}
