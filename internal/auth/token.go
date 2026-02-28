package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var errInvalidAccessToken = errors.New("invalid access token")

func (s *Service) verifyAccessToken(token string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithLeeway(time.Minute),
	)

	parsedToken, err := parser.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return s.signingKey, nil
	})
	if err != nil || !parsedToken.Valid {
		return "", errInvalidAccessToken
	}
	if claims.ExpiresAt == nil {
		return "", errInvalidAccessToken
	}
	subject := strings.TrimSpace(claims.Subject)
	if subject == "" {
		return "", errInvalidAccessToken
	}
	return subject, nil
}
