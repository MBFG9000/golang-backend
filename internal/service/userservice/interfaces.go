package userservice

import (
	"taskmanager/pkg/modules"
)

type UserUseCaseInterface interface {
	CreateUser(user modules.User) (int, error)
	GetAllUsers() ([]modules.User, error)
	GetUserByID(id int) (modules.User, error)
	UpdateUserByID(id int, user UpdateUserInput) error
	DeleteUserByID(id int) error
	SoftDeleteUserByID(id int) error
}
