package users

import (
	"database/sql"
	"fmt"
	"strings"
	"github.com/MBFG9000/golang-backend/internal/domain"
	postgres "github.com/MBFG9000/golang-backend/internal/repository/postgresql"
	"github.com/MBFG9000/golang-backend/internal/utils"
	"time"

	"github.com/google/uuid"
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

func (r *Repository) GetUsers() ([]domain.User, error) {
	var users []domain.User
	err := r.db.DB.Select(&users, "SELECT * FROM users WHERE deleted_at IS NULL")
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) GetPaginatedUsers(page, pageSize int, filter domain.UserFilter) (domain.PaginatedResponse[domain.User], error) {
	offset := (page - 1) * pageSize
	args := []any{}
	i := 1

	// WHERE conditions
	conditions := []string{"deleted_at IS NULL"}

	if filter.ID != nil {
		conditions = append(conditions, fmt.Sprintf("id = $%d", i))
		args = append(args, *filter.ID)
		i++
	}
	if filter.FirstName != nil {
		conditions = append(conditions, fmt.Sprintf("first_name ILIKE $%d", i))
		args = append(args, "%"+*filter.FirstName+"%")
		i++
	}
	if filter.Email != nil {
		conditions = append(conditions, fmt.Sprintf("email = $%d", i))
		args = append(args, *filter.Email)
		i++
	}
	if filter.Gender != nil {
		conditions = append(conditions, fmt.Sprintf("gender = $%d", i))
		args = append(args, *filter.Gender)
		i++
	}
	if filter.BirthDate != nil {
		conditions = append(conditions, fmt.Sprintf("birth_date = $%d", i))
		args = append(args, *filter.BirthDate)
		i++
	}

	where := "WHERE " + strings.Join(conditions, " AND ")
	orderBy := fmt.Sprintf("ORDER BY %s %s", filter.SortBy, strings.ToUpper(filter.SortOrder))

	// count
	var totalCount int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users %s", where)
	err := r.db.DB.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return domain.PaginatedResponse[domain.User]{}, err
	}

	// data
	dataArgs := append(args, pageSize, offset)
	query := fmt.Sprintf(
		"SELECT * FROM users %s %s LIMIT $%d OFFSET $%d",
		where, orderBy, i, i+1,
	)

	var users []domain.User
	err = r.db.DB.Select(&users, query, dataArgs...)
	if err != nil {
		return domain.PaginatedResponse[domain.User]{}, err
	}

	return domain.PaginatedResponse[domain.User]{
		Data:       users,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (r *Repository) CreateUser(user domain.User) (uuid.UUID, error) {
	var id uuid.UUID

	rows, err := r.db.DB.NamedQuery(`INSERT INTO users
		(first_name, last_name, email, phone, city, country, zip)
	VALUES (:first_name, :last_name, :email, :phone, :city, :country, :zip)
	RETURNING id`, &user)

	if err != nil {
		return uuid.Nil, fmt.Errorf("CreateUser: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return uuid.Nil, fmt.Errorf("CreateUser scan: %w", err)
		}
	}

	return id, nil
}

func (r *Repository) GetCommonFriends(userID1, userID2 uuid.UUID) ([]domain.User, error) {
	var users []domain.User

	query := `
        SELECT u.*
        FROM users u
        JOIN user_friends uf1 ON u.id = uf1.friend_id
        JOIN user_friends uf2 ON u.id = uf2.friend_id
        WHERE uf1.user_id = $1
          AND uf2.user_id = $2
          AND u.deleted_at IS NULL
    `

	err := r.db.DB.Select(&users, query, userID1, userID2)
	if err != nil {
		return nil, fmt.Errorf("GetCommonFriends: %w", err)
	}

	return users, nil
}

func (r *Repository) UpdateUserByID(id uuid.UUID, user domain.User) error {
	user.ID = id
	result, err := r.db.DB.NamedExec(`UPDATE users SET
    first_name=:first_name, last_name=:last_name, email=:email,
    phone=:phone, city=:city, country=:country, zip=:zip,
    gender=:gender, birth_date=:birth_date,
    updated_at=now()
    WHERE id=:id`, &user)

	if err != nil {
		return fmt.Errorf("UpdateUserByID: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("UpdateUserByID RowsAffected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("User with ID %s not found ErrCode: %w", id, utils.ErrObjectNotFound)
	}

	return nil
}

func (r *Repository) GetUserByID(id uuid.UUID) (domain.User, error) {
	var user domain.User

	err := r.db.DB.Get(&user, "SELECT * FROM users WHERE id=$1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, fmt.Errorf("User with ID %s not found ErrCode: %w", id, utils.ErrObjectNotFound)
		}
		return domain.User{}, fmt.Errorf("GetUserByID: %w", err)
	}

	return user, nil
}

func (r *Repository) DeleteUserByID(id uuid.UUID) error {
	result, err := r.db.DB.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("DeleteUserByID: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteUserByID RowsAffected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("User with ID %s not found ErrCode: %w", id, utils.ErrObjectNotFound)
	}

	return nil
}

func (r *Repository) SoftDeleteUserByID(id uuid.UUID) error {
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
		return fmt.Errorf("User with ID %s not found ErrCode: %w", id, utils.ErrObjectNotFound)
	}

	return nil
}
