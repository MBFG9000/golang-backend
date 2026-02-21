package repository

import (
	postgresrepo "taskmanager/internal/repository/postgresql"
	"taskmanager/internal/repository/postgresql/users"
)

type Repositories struct {
	UserRepository
}

func NewRepositories(db *postgresrepo.Dialect) *Repositories {
	return &Repositories{
		UserRepository: users.NewUserRepository(db),
	}
}
