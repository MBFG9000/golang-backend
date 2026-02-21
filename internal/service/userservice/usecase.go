package userservice

import (
	"fmt"
	"taskmanager/internal/repository"
	"taskmanager/internal/utils"
	"taskmanager/pkg/modules"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type UserUseCase struct {
	repo repository.UserRepository
}

func NewUserUseCase(repo repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		repo: repo,
	}
}

func (u *UserUseCase) CreateUser(user modules.User) (int, error) {

	err := ValidateUser(user)

	if err != nil {
		return -1, err
	}

	id, err := u.repo.CreateUser(user)

	if err != nil {
		return -1, err
	}

	return id, nil
}

func (u *UserUseCase) GetAllUsers() ([]modules.User, error) {
	return u.repo.GetUsers()
}

func (u *UserUseCase) GetUserByID(id int) (modules.User, error) {
	err := ValidateID(id)

	if err != nil {
		return modules.User{}, err
	}

	return u.repo.GetUserByID(id)
}

func (u *UserUseCase) UpdateUserByID(id int, input UpdateUserInput) error {
	if err := ValidateID(id); err != nil {
		return err
	}

	user, err := u.repo.GetUserByID(id)

	if err != nil {
		return err
	}

	if input.FirstName != nil {
		user.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		user.LastName = *input.LastName
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Phone != nil {
		user.Phone = input.Phone
	}
	if input.City != nil {
		user.City = input.City
	}
	if input.Country != nil {
		user.Country = input.Country
	}
	if input.Zip != nil {
		user.Zip = input.Zip
	}

	if err := ValidateUser(user); err != nil {
		return err
	}

	return u.repo.UpdateUserByID(id, user)
}

func (u *UserUseCase) DeleteUserByID(id int) error {
	if err := ValidateID(id); err != nil {
		return err
	}

	return u.repo.DeleteUserByID(id)
}

func (u *UserUseCase) SoftDeleteUserByID(id int) error {
	if err := ValidateID(id); err != nil {
		return err
	}

	return u.repo.SoftDeleteUserByID(id)
}

// *********************
// Validation functions
// *********************
func ValidateID(id int) error {
	err := validation.Validate(
		id,
		validation.Required,
		validation.Min(1),
	)

	if err != nil {
		return fmt.Errorf("%w: %w", utils.ErrInvalidData, err)
	}

	return nil
}

func ValidateUser(user modules.User) error {
	err := validation.ValidateStruct(&user,
		validation.Field(
			&user.FirstName,
			validation.Required,
			validation.Length(1, 100),
			is.Alpha,
		),
		validation.Field(
			&user.LastName,
			validation.Required,
			validation.Length(1, 100),
			is.Alpha,
		),
		validation.Field(
			&user.Email,
			validation.Required,
			validation.Length(1, 255),
			is.Email,
		),
		validation.Field(
			&user.Phone,
			validation.NilOrNotEmpty,
			validation.Length(1, 20),
		),
		validation.Field(
			&user.City,
			validation.NilOrNotEmpty,
			validation.Length(1, 100),
		),

		validation.Field(
			&user.Country,
			validation.NilOrNotEmpty,
			validation.Length(1, 100),
		),

		validation.Field(
			&user.Zip,
			validation.NilOrNotEmpty,
			validation.Length(1, 20),
		),
	)

	if err != nil {
		return fmt.Errorf("%w: %w", utils.ErrInvalidData, err)
	}

	return nil
}
