package repository

import (
	"errors"

	"github.com/aodihis/go-rest-signup-login/database"
	"github.com/aodihis/go-rest-signup-login/internal/models"
	"github.com/aodihis/go-rest-signup-login/internal/utils"
	"github.com/lib/pq"
)

type UserRepository struct {
}

func (r *UserRepository) CreateUser(user *models.User) error {
	if err := user.HashPassword(); err != nil {
		return err
	}
	user.CreatedAt = utils.Now()
	err := database.DB.QueryRow("INSERT INTO users (email, password, is_active, created_at) VALUES ($1, $2, $3, $4) RETURNING id", user.Email, user.Password, user.IsActive, user.CreatedAt).Scan(&user.ID)
	if err != nil {
		if sqlErr, ok := err.(*pq.Error); ok && sqlErr.Code == "23505" {
			return errors.New("email is not available")
		}
		return err
	}
	return nil
}

func (r *UserRepository) FindUserByEmail(email string) (*models.User, error) {
	user := &models.User{}

	query := `SELECT id, email, password, is_active, last_login, created_at FROM users WHERE email = $1`

	err := database.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.IsActive,
		&user.LastLogin,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
