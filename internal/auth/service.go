package auth

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/keithics/devops-dashboard/api/internal/db/sqlc"
	"golang.org/x/crypto/bcrypt"
)

func NewHandler(q *sqlc.Queries, tokenKey string) *Handler {
	return &Handler{
		svc: &Service{
			q:          q,
			signingKey: []byte(tokenKey),
			tokenTTL:   time.Hour,
		},
	}
}

func (s *Service) Register(ctx context.Context, req registerRequest) (sqlc.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return sqlc.User{}, err
	}

	return s.q.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        req.Email,
		PasswordHash: string(hash),
	})
}

func (s *Service) Login(ctx context.Context, req loginRequest) (sqlc.User, string, int64, error) {
	user, err := s.q.GetUserByEmail(ctx, req.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		return sqlc.User{}, "", 0, errInvalidCredentials
	}
	if err != nil {
		return sqlc.User{}, "", 0, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return sqlc.User{}, "", 0, errInvalidCredentials
	}

	accessToken, expiresIn, err := s.generateAccessToken(user.ID)
	if err != nil {
		return sqlc.User{}, "", 0, err
	}

	return user, accessToken, expiresIn, nil
}

func (s *Service) ForgotPassword(ctx context.Context, req forgotPasswordRequest) error {
	user, err := s.q.GetUserByEmail(ctx, req.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}

	resetToken, err := generateResetToken()
	if err != nil {
		return err
	}

	log.Printf("password reset requested for user_id=%s email=%s token=%s", user.ID, user.Email, resetToken)
	return nil
}

func (s *Service) generateAccessToken(userID string) (string, int64, error) {
	expiresAt := time.Now().Add(s.tokenTTL).Unix()
	payloadBytes, err := json.Marshal(map[string]string{
		"sub": userID,
		"exp": strconv.FormatInt(expiresAt, 10),
	})
	if err != nil {
		return "", 0, err
	}

	encodedPayload := base64.RawURLEncoding.EncodeToString(payloadBytes)
	mac := hmac.New(sha256.New, s.signingKey)
	if _, err := mac.Write([]byte(encodedPayload)); err != nil {
		return "", 0, err
	}
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return encodedPayload + "." + signature, int64(s.tokenTTL.Seconds()), nil
}

func generateResetToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
