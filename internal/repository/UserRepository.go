package repository

import (
	"taskmanager/pkg/modules"
)

type UserRepository interface {
	CreateUser(modules.User) (int, error)
	UpdateUserByID(int, modules.User) error
	GetUsers() ([]modules.User, error)
	GetUserByID(int) (modules.User, error)
	DeleteUserByID(int) error
	SoftDeleteUserByID(int) error
}
