package services

import (
	"errors"
	"strings"

	"github.com/aodihis/go-rest-signup-login/internal/models"
	"github.com/aodihis/go-rest-signup-login/internal/repository"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: &repository.UserRepository{},
	}
}

func (s *AuthService) SignUp(email string, password string, confirm string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters long")
	}

	if password != confirm {
		return nil, errors.New("password and confirm password do not match")
	}

	user := &models.User{
		Email:    strings.ToLower(strings.TrimSpace(email)),
		Password: password,
		IsActive: true,
	}

	err := s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
