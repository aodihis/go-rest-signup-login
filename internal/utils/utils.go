package utils

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

var HashPassword = func(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

var Now = func() time.Time {
	return time.Now()
}
