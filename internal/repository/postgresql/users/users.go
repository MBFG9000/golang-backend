package users

import (
	"database/sql"
	"fmt"
	postgres "taskmanager/internal/repository/postgresql"
	"taskmanager/internal/utils"
	"taskmanager/pkg/modules"
	"time"
)

type Repository struct {
	db               *postgres.Dialect
	executionTimeout time.Duration
}

func NewUserRepository(db *postgres.Dialect) *Repository {
	return &Repository{
		db:               db,
		executionTimeout: time.Second * 5,
	}
}

func (r *Repository) GetUsers() ([]modules.User, error) {
	var users []modules.User
	err := r.db.DB.Select(&users, "SELECT * FROM users")

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Repository) CreateUser(user modules.User) (int, error) {
	var id int

	rows, err := r.db.DB.NamedQuery(`INSERT INTO users
		(first_name, last_name, email, phone, city, country, zip) 
	VALUES (:first_name, :last_name, :email, :phone, :city, :country, :zip)
	RETURNING id`, &user)

	if err != nil {
		return 0, fmt.Errorf("CreateUser: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		rows.Scan(&id)
	}

	return id, nil
}

func (r *Repository) UpdateUserByID(id int, user modules.User) error {
	user.ID = id
	result, err := r.db.DB.NamedExec(`UPDATE users SET
	first_name=:first_name,last_name=:last_name,email=:email,
	phone=:phone,city=:city,country=:country,zip=:zip,
	updated_at=now() 
	WHERE id=:id`, &user)

	if err != nil {
		return fmt.Errorf("UpdatedUserByID: %w", err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf("UpdateUserByID RowsAffected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("User with ID %d not found ErrCode: %w", id, utils.ErrObjectNotFound)
	}

	return nil

}

func (r *Repository) GetUserByID(id int) (modules.User, error) {
	var user modules.User

	err := r.db.DB.Get(&user, "SELECT * FROM users WHERE id=$1", id)

	if err != nil {
		if err == sql.ErrNoRows {
			return modules.User{}, fmt.Errorf("User with ID %d not found ErrCode: %w", id, utils.ErrObjectNotFound)
		}

		return modules.User{}, fmt.Errorf("GetUserByID: %w", err)
	}

	return user, nil
}

func (r *Repository) DeleteUserByID(id int) error {
	result, err := r.db.DB.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("DeleteUserByID: %w", err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf("DeleteUserByID RowsAffected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("User with ID %d not found ErrCode: %w", id, utils.ErrObjectNotFound)
	}

	return nil
}

func (r *Repository) SoftDeleteUserByID(id int) error {
	result, err := r.db.DB.Exec(
		"UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL", id)

	if err != nil {
		return fmt.Errorf("SoftDeleteUserByID: %w", err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf("SoftDeleteUserByID RowsAffected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("User with ID %d not found ErrCode: %w", id, utils.ErrObjectNotFound)
	}

	return nil
}
