package models

import (
	"regexp"
	"time"

	"github.com/aodihis/go-rest-signup-login/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

// User struct represents the user model
type User struct {
	ID        int        `json:"id"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	IsActive  bool       `json:"is_active"`
	LastLogin *time.Time `json:"last_login"`
	CreatedAt time.Time  `json:"created_at"`
}

func (u *User) HashPassword() error {
	hashed, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	u.Password = string(hashed)
	return nil
}

func (u *User) IsValidEmail() bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(u.Email)
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
