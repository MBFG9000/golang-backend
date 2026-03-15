package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserFilter struct {
	ID        *uuid.UUID
	FirstName *string
	Email     *string
	Gender    *string
	BirthDate *time.Time
	SortBy    string // "id", "first_name", "email", "gender", "birth_date"
	SortOrder string // "asc", "desc"
}
