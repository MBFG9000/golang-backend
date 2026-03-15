package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	FirstName string     `db:"first_name" json:"first_name"`
	LastName  string     `db:"last_name" json:"last_name"`
	Email     string     `db:"email" json:"email"`
	Phone     *string    `db:"phone" json:"phone"`
	City      *string    `db:"city" json:"city"`
	Country   *string    `db:"country" json:"country"`
	Zip       *string    `db:"zip" json:"zip"`
	Gender    *string    `db:"gender" json:"gender"`
	BirthDate *time.Time `db:"birth_date" json:"birth_date"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}

type UserRepository interface {
	CreateUser(User) (uuid.UUID, error)
	UpdateUserByID(uuid.UUID, User) error
	GetUsers() ([]User, error)
	GetUserByID(uuid.UUID) (User, error)
	GetPaginatedUsers(page int, pageSize int, filter UserFilter) (PaginatedResponse[User], error)
	DeleteUserByID(uuid.UUID) error
	SoftDeleteUserByID(uuid.UUID) error
	GetCommonFriends(userID1, userID2 uuid.UUID) ([]User, error)
}
