package repository

import (
	"github.com/MBFG9000/golang-backend/internal/domain"
	postgresrepo "github.com/MBFG9000/golang-backend/internal/repository/postgresql"
	"github.com/MBFG9000/golang-backend/internal/repository/postgresql/users"
)

type Repositories struct {
	domain.UserRepository
}

func NewRepositories(db *postgresrepo.Dialect) *Repositories {
	return &Repositories{
		UserRepository: users.NewUserRepository(db),
	}
}
