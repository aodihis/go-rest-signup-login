package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aodihis/go-rest-signup-login/config"
	"github.com/aodihis/go-rest-signup-login/internal/services"
	"github.com/golang-jwt/jwt/v4"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: services.NewAuthService(),
	}
}

type SignUpRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type JWTClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func response(w http.ResponseWriter, res map[string]interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(res)
}
func (c *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		res := map[string]interface{}{
			"message": "method not allowed",
			"status":  "failed",
		}
		response(w, res, http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		res := map[string]interface{}{
			"message": "invalid content type, expected application/json",
			"status":  "failed",
		}
		response(w, res, http.StatusUnsupportedMediaType)
		return

	}

	var req SignUpRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res := map[string]interface{}{
			"message": "invalid input",
			"status":  "failed",
		}
		response(w, res, http.StatusBadRequest)
		return
	}

	user, err := c.authService.SignUp(req.Email, req.Password, req.ConfirmPassword)

	if err != nil {
		res := map[string]interface{}{
			"message": err.Error(),
			"status":  "failed",
		}
		response(w, res, http.StatusBadRequest)
		return
	}

	res := map[string]interface{}{
		"message": "user created",
		"status":  "success",
		"data": map[string]interface{}{
			"user_id": user.ID,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (c *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		res := map[string]interface{}{
			"message": "method not allowed",
			"status":  "failed",
		}
		response(w, res, http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		res := map[string]interface{}{
			"message": "invalid content type, expected application/json",
			"status":  "failed",
		}
		response(w, res, http.StatusUnsupportedMediaType)
		return
	}

	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res := map[string]interface{}{
			"message": "invalid input",
			"status":  "failed",
		}
		response(w, res, http.StatusBadRequest)
		return
	}

	user, err := c.authService.Login(req.Email, req.Password)
	if err != nil {
		res := map[string]interface{}{
			"message": err.Error(),
			"status":  "failed",
		}
		response(w, res, http.StatusUnauthorized)
		return
	}

	secretKey := config.GetEnv("JWT_SECRET_KEY")

	if secretKey == "" {
		res := map[string]interface{}{
			"message": "internal server error",
			"status":  "failed",
		}
		response(w, res, http.StatusInternalServerError)
		return
	}

	expirationTime := time.Now().Add(168 * time.Hour)
	claims := JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))

	if err != nil {
		res := map[string]interface{}{
			"message": "internal server error",
			"status":  "failed",
		}
		response(w, res, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Login successful",
		"status":  "failed",
		"data": map[string]interface{}{
			"user_id": user.ID,
			"token":   tokenString,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
