package apikey

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{
		svc: &Service{pool: pool},
	}
}

func (s *Service) Create(ctx context.Context, name string) (createAPIKeyResponse, error) {
	rawKey, err := generateAPIKey()
	if err != nil {
		return createAPIKeyResponse{}, err
	}

	keyHash := hashAPIKey(rawKey)
	last4 := last4(rawKey)

	var result createAPIKeyResponse
	err = s.pool.QueryRow(ctx, `
INSERT INTO api_keys (name, key_hash, key_last4)
VALUES ($1, $2, $3)
RETURNING id::text, name, key_last4, created_at
`, strings.TrimSpace(name), keyHash, last4).Scan(&result.ID, &result.Name, &result.Last4, &result.CreatedAt)
	if err != nil {
		return createAPIKeyResponse{}, err
	}

	result.Key = rawKey
	return result, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	_, err := s.pool.Exec(ctx, `
DELETE FROM api_keys
WHERE id = $1::uuid
`, id)
	return err
}

func (s *Service) Validate(ctx context.Context, rawKey string) (bool, error) {
	keyHash := hashAPIKey(strings.TrimSpace(rawKey))

	var id string
	err := s.pool.QueryRow(ctx, `
SELECT id::text
FROM api_keys
WHERE key_hash = $1
LIMIT 1
`, keyHash).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func generateAPIKey() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return "bisky_" + base64.RawURLEncoding.EncodeToString(buf), nil
}

func hashAPIKey(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func last4(value string) string {
	if len(value) <= 4 {
		return value
	}
	return value[len(value)-4:]
}
