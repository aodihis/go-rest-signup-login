package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aodihis/go-rest-signup-login/internal/services"
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

func (c *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Invalid content type, expected application/json", http.StatusUnsupportedMediaType)
		return
	}

	var req SignUpRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("%v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user, err := c.authService.SignUp(req.Email, req.Password, req.ConfirmPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := map[string]interface{}{
		"message": "User created",
		"user_id": user.ID,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}
