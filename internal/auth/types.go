package auth

import (
	"errors"
	"time"

	"github.com/keithics/devops-dashboard/api/internal/db/sqlc"
)

const (
	ctxRegisterRequestKey       = "auth.register.request"
	ctxLoginRequestKey          = "auth.login.request"
	ctxForgotPasswordRequestKey = "auth.forgot-password.request"
	ctxUserIDKey                = "auth.user.id"
)

var errInvalidCredentials = errors.New("invalid credentials")

type Handler struct {
	svc *Service
}

type Service struct {
	q          *sqlc.Queries
	signingKey []byte
	tokenTTL   time.Duration
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type forgotPasswordRequest struct {
	Email string `json:"email"`
}

type userResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type registerResponse struct {
	User userResponse `json:"user"`
}

type loginResponse struct {
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	ExpiresIn   int64        `json:"expires_in"`
	User        userResponse `json:"user"`
}

type forgotPasswordResponse struct {
	Message string `json:"message"`
}
