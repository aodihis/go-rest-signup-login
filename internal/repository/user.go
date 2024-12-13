package repository

import (
	"errors"
	"time"

	"github.com/aodihis/go-rest-signup-login/database"
	"github.com/aodihis/go-rest-signup-login/internal/models"
	"github.com/lib/pq"
)

type UserRepository struct {
}

func (r *UserRepository) CreateUser(user *models.User) error {

	if err := user.HashPassword(); err != nil {
		return err
	}
	user.CreatedAt = time.Now()
	err := database.DB.QueryRow("INSERT INTO users (email, password, is_active, created_at) VALUES ($1, $2, $3, $4) RETURNING id", user.Email, user.Password, user.IsActive, user.CreatedAt).Scan(&user.ID)
	if err != nil {
		if sqlErr, ok := err.(*pq.Error); ok && sqlErr.Code == "23505" {
			return errors.New("user with this email already exists")
		}
		return err
	}
	return nil
}
