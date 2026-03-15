package userservice

import (
	"github.com/MBFG9000/golang-backend/internal/domain"

	"github.com/google/uuid"
)

type UserServiceInterface interface {
	CreateUser(user domain.User) (uuid.UUID, error)
	GetAllUsers() ([]domain.User, error)
	GetUserByID(id uuid.UUID) (domain.User, error)
	GetPaginatedUsers(page, pageSize int, filter domain.UserFilter) (domain.PaginatedResponse[domain.User], error)
	GetCommonFriends(userID1, userID2 uuid.UUID) ([]domain.User, error)
	UpdateUserByID(id uuid.UUID, input UpdateUserInput) error
	DeleteUserByID(id uuid.UUID) error
	SoftDeleteUserByID(id uuid.UUID) error
}
