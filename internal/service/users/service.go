package userservice

import (
	"fmt"
	"strings"
	"time"

	"github.com/MBFG9000/golang-backend/internal/domain"
	"github.com/MBFG9000/golang-backend/internal/utils"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
)

type UserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (u *UserService) CreateUser(user domain.User) (uuid.UUID, error) {
	if err := ValidateUser(user); err != nil {
		return uuid.Nil, err
	}

	return u.repo.CreateUser(user)
}

func (u *UserService) GetAllUsers() ([]domain.User, error) {
	return u.repo.GetUsers()
}

func (u *UserService) GetUserByID(id uuid.UUID) (domain.User, error) {
	if err := ValidateUUID(id); err != nil {
		return domain.User{}, err
	}

	return u.repo.GetUserByID(id)
}

func (u *UserService) GetPaginatedUsers(page, pageSize int, filter domain.UserFilter) (domain.PaginatedResponse[domain.User], error) {
	if err := ValidatePagination(page, pageSize); err != nil {
		return domain.PaginatedResponse[domain.User]{}, err
	}

	normalizedFilter, err := NormalizeUserFilter(filter)
	if err != nil {
		return domain.PaginatedResponse[domain.User]{}, err
	}

	return u.repo.GetPaginatedUsers(page, pageSize, normalizedFilter)
}

func (u *UserService) GetCommonFriends(userID1, userID2 uuid.UUID) ([]domain.User, error) {
	if err := ValidateUUID(userID1); err != nil {
		return nil, err
	}
	if err := ValidateUUID(userID2); err != nil {
		return nil, err
	}
	if userID1 == userID2 {
		return nil, fmt.Errorf("%w: user ids must be different", utils.ErrInvalidData)
	}

	return u.repo.GetCommonFriends(userID1, userID2)
}

func (u *UserService) UpdateUserByID(id uuid.UUID, input UpdateUserInput) error {
	if err := ValidateUUID(id); err != nil {
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
	if input.Gender != nil {
		user.Gender = input.Gender
	}
	if input.BirthDate != nil {
		user.BirthDate = input.BirthDate
	}

	if err := ValidateUser(user); err != nil {
		return err
	}

	return u.repo.UpdateUserByID(id, user)
}

func (u *UserService) DeleteUserByID(id uuid.UUID) error {
	if err := ValidateUUID(id); err != nil {
		return err
	}

	return u.repo.DeleteUserByID(id)
}

func (u *UserService) SoftDeleteUserByID(id uuid.UUID) error {
	if err := ValidateUUID(id); err != nil {
		return err
	}

	return u.repo.SoftDeleteUserByID(id)
}

func ValidateUUID(id uuid.UUID) error {
	err := validation.Validate(id, validation.Required, is.UUID)

	return err
}

func ValidatePagination(page, pageSize int) error {
	if err := validation.Validate(page, validation.Required, validation.Min(1)); err != nil {
		return fmt.Errorf("%w: invalid page: %w", utils.ErrInvalidData, err)
	}
	if err := validation.Validate(pageSize, validation.Required, validation.Min(1), validation.Max(100)); err != nil {
		return fmt.Errorf("%w: invalid page_size: %w", utils.ErrInvalidData, err)
	}

	return nil
}

func NormalizeUserFilter(filter domain.UserFilter) (domain.UserFilter, error) {
	if filter.ID != nil {
		if err := ValidateUUID(*filter.ID); err != nil {
			return domain.UserFilter{}, err
		}
	}

	if filter.Email != nil {
		email := strings.TrimSpace(*filter.Email)
		if email != "" {
			if err := validation.Validate(email, is.Email); err != nil {
				return domain.UserFilter{}, fmt.Errorf("%w: invalid email filter: %w", utils.ErrInvalidData, err)
			}
			filter.Email = &email
		}
	}

	if filter.FirstName != nil {
		firstName := strings.TrimSpace(*filter.FirstName)
		filter.FirstName = &firstName
	}

	if filter.Gender != nil {
		gender := strings.ToLower(strings.TrimSpace(*filter.Gender))
		if gender != "" {
			if err := validation.Validate(gender, validation.In("male", "female")); err != nil {
				return domain.UserFilter{}, fmt.Errorf("%w: invalid gender filter: %w", utils.ErrInvalidData, err)
			}
		}
		filter.Gender = &gender
	}

	if filter.BirthDate != nil && filter.BirthDate.After(time.Now()) {
		return domain.UserFilter{}, fmt.Errorf("%w: birth_date cannot be in the future", utils.ErrInvalidData)
	}

	sortBy := strings.ToLower(strings.TrimSpace(filter.SortBy))
	if sortBy == "" {
		sortBy = "id"
	}

	if err := validation.Validate(sortBy, validation.In("id", "first_name", "email", "gender", "birth_date")); err != nil {
		return domain.UserFilter{}, fmt.Errorf("%w: invalid sort_by: %w", utils.ErrInvalidData, err)
	}

	sortOrder := strings.ToLower(strings.TrimSpace(filter.SortOrder))
	if sortOrder == "" {
		sortOrder = "asc"
	}

	if err := validation.Validate(sortOrder, validation.In("asc", "desc")); err != nil {
		return domain.UserFilter{}, fmt.Errorf("%w: invalid sort_order: %w", utils.ErrInvalidData, err)
	}

	filter.SortBy = sortBy
	filter.SortOrder = sortOrder
	return filter, nil
}

func ValidateUser(user domain.User) error {
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

	if user.Gender != nil {
		gender := strings.ToLower(strings.TrimSpace(*user.Gender))
		if err := validation.Validate(gender, validation.In("male", "female")); err != nil {
			return fmt.Errorf("%w: invalid gender: %w", utils.ErrInvalidData, err)
		}
		user.Gender = &gender
	}

	if user.BirthDate != nil && user.BirthDate.After(time.Now()) {
		return fmt.Errorf("%w: birth_date cannot be in the future", utils.ErrInvalidData)
	}

	return nil
}
