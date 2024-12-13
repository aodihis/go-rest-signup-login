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

	if !user.IsValidEmail() {
		return nil, errors.New("invalid email")
	}

	err := s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(email string, password string) (*models.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	user, err := s.userRepo.FindUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := user.CheckPassword(password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
