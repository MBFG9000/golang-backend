package modules

import "time"

type User struct {
	ID        int        `db:"id" json:"id"`
	FirstName string     `db:"first_name" json:"first_name"`
	LastName  string     `db:"last_name" json:"last_name"`
	Email     string     `db:"email" json:"email"`
	Phone     *string    `db:"phone" json:"phone"`
	City      *string    `db:"city" json:"city"`
	Country   *string    `db:"country" json:"country"`
	Zip       *string    `db:"zip" json:"zip"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}
