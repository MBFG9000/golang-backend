package userservice

import "time"

type UpdateUserInput struct {
	FirstName *string    `json:"first_name"`
	LastName  *string    `json:"last_name"`
	Email     *string    `json:"email"`
	Phone     *string    `json:"phone"`
	City      *string    `json:"city"`
	Country   *string    `json:"country"`
	Zip       *string    `json:"zip"`
	Gender    *string    `json:"gender"`
	BirthDate *time.Time `json:"birth_date"`
}
