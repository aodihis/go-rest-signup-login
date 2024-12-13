package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aodihis/go-rest-signup-login/database"
	"github.com/aodihis/go-rest-signup-login/internal/utils"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
)

func dbSetup() (sqlmock.Sqlmock, error) {
	mockDb, mock, err := sqlmock.New()

	if err != nil {
		return nil, err
	}
	// defer database.DB.Close()
	database.DB = mockDb
	return mock, nil
}

func TestSignUp(t *testing.T) {

	mock, err := dbSetup()

	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}

	mockHashPassword := "$2a$10$lu7Yqwu6qpLyDsUXP4fXXOK.tiu3xY9igPvsxNK3nTdGm7bFrYbHi"
	fixedTime := time.Now()
	utils.HashPassword = func(password string) ([]byte, error) {
		return []byte(mockHashPassword), nil
	}

	utils.Now = func() time.Time {
		return fixedTime
	}

	tests := []struct {
		name            string
		email           string
		password        string
		confirmPassword string
		expectedStatus  int
		expectedError   string
		expectedQuery   string
	}{
		{
			name:            "Valid SignUp",
			email:           "test@example.com",
			password:        "password123",
			confirmPassword: "password123",
			expectedStatus:  http.StatusCreated,
			expectedError:   "",
			expectedQuery:   "",
		},
		{
			name:            "Empty email",
			email:           "",
			password:        "password123",
			confirmPassword: "password123",
			expectedStatus:  http.StatusBadRequest,
			expectedError:   "email is required",
			expectedQuery:   "no_query",
		},
		{
			name:            "Empty password",
			email:           "les@test.com",
			password:        "",
			confirmPassword: "",
			expectedStatus:  http.StatusBadRequest,
			expectedError:   "password must be at least 8 characters long",
			expectedQuery:   "no_query",
		},
		{
			name:            "Invalid Email Format",
			email:           "test@",
			password:        "password123",
			confirmPassword: "password123",
			expectedStatus:  http.StatusBadRequest,
			expectedError:   "invalid email",
			expectedQuery:   "no_query",
		},
		{
			name:            "Using Registered Email",
			email:           "test@example.com",
			password:        "password123",
			confirmPassword: "password123",
			expectedStatus:  http.StatusBadRequest,
			expectedError:   "email is not available",
			expectedQuery:   "duplicate",
		},
		{
			name:            "Password less than 8",
			email:           "test123@example.com",
			password:        "ppp",
			confirmPassword: "ppp",
			expectedStatus:  http.StatusBadRequest,
			expectedError:   "password must be at least 8 characters long",
			expectedQuery:   "no_query",
		},
		{
			name:            "Password and confirm not match",
			email:           "test123@example.com",
			password:        "password123",
			confirmPassword: "password1234",
			expectedStatus:  http.StatusBadRequest,
			expectedError:   "password and confirm password do not match",
			expectedQuery:   "no_query",
		},
	}

	authController := NewAuthHandler()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := SignUpRequest{Email: tt.email, Password: tt.password, ConfirmPassword: tt.confirmPassword}

			if err != nil {
				t.Errorf("Cannot hash password")
			}

			if tt.expectedQuery == "duplicate" {
				mock.ExpectQuery("INSERT INTO users").
					WithArgs(user.Email, mockHashPassword, true, fixedTime).
					WillReturnError(&pq.Error{
						Code:    "23505",
						Message: "duplicate key value violates unique constraint",
						Detail:  "Key (email)=() already exists.",
					})
			} else if tt.expectedQuery == "" {
				mock.ExpectQuery("INSERT INTO users").
					WithArgs(user.Email, mockHashPassword, true, fixedTime).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1)) // Simulating returning ID 1
			}
			body, _ := json.Marshal(user)

			req := httptest.NewRequest("POST", "/signup", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			authController.SignUp(w, req)
			res := w.Result()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("SignUp() status = %v, expected %v", res.StatusCode, tt.expectedStatus)
			}

			if tt.expectedError != "" {
				var ersp map[string]string
				json.NewDecoder(res.Body).Decode(&ersp)
				if ersp["message"] != tt.expectedError {
					t.Errorf("SignUp() error = %v, expected %v", ersp["message"], tt.expectedError)
				}
			}

		})
	}
}

func TestLogin(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatalf("Failed to load .env: %v", err)
	}
	mock, err := dbSetup()

	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}

	mockHashPassword := "$2a$10$lu7Yqwu6qpLyDsUXP4fXXOK.tiu3xY9igPvsxNK3nTdGm7bFrYbHi"
	fixedTime := time.Now()
	utils.HashPassword = func(password string) ([]byte, error) {
		return []byte(mockHashPassword), nil
	}

	utils.Now = func() time.Time {
		return fixedTime
	}

	authController := NewAuthHandler()

	tests := []struct {
		name           string
		email          string
		password       string
		expectedStatus int
		expectedError  string
		expectedQuery  string
	}{
		{
			name:           "Valid Login",
			email:          "test@example.com",
			password:       "password123",
			expectedStatus: http.StatusOK,
			expectedError:  "",
			expectedQuery:  "",
		},
		{
			name:           "Missing Email",
			email:          "",
			password:       "password123",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "invalid credentials",
			expectedQuery:  "not_found",
		},
		{
			name:           "Missing Password",
			email:          "test@example.com",
			password:       "",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "invalid credentials",
			expectedQuery:  "not_found",
		},
		{
			name:           "Invalid Email",
			email:          "invalidemail",
			password:       "password123",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "invalid credentials",
			expectedQuery:  "not_found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := LoginRequest{Email: tt.email, Password: tt.password}
			body, _ := json.Marshal(user)

			if tt.expectedQuery == "not_found" {
				mock.ExpectQuery("SELECT id, email, password, is_active, last_login, created_at FROM users").
					WillReturnError(sql.ErrNoRows)
			} else if tt.expectedQuery == "" {
				mock.ExpectQuery("SELECT id, email, password, is_active, last_login, created_at FROM users").
					WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "is_active", "last_login", "created_at"}).
						AddRow(1, strings.ToLower(tt.email), mockHashPassword, true, nil, fixedTime))
			}

			req := httptest.NewRequest("POST", "/signup", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			authController.Login(w, req)
			res := w.Result()
			if res.StatusCode != tt.expectedStatus {
				t.Errorf("Login() status = %v, expected %v", res.StatusCode, tt.expectedStatus)
			}

			if tt.expectedError != "" {
				var errResp map[string]string
				json.NewDecoder(res.Body).Decode(&errResp)
				if errResp["message"] != tt.expectedError {
					t.Errorf("Login() error = %v, expected %v", errResp["error"], tt.expectedError)
				}
			}

		})
	}
}
