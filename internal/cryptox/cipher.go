package cryptox

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

func NewCipher(secret string) (*Cipher, error) {
	if secret == "" {
		return nil, errors.New("token encryption key is required")
	}

	derived := sha256.Sum256([]byte(secret))
	block, err := aes.NewCipher(derived[:])
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &Cipher{aead: aead}, nil
}

func (c *Cipher) Encrypt(plain string) (string, error) {
	if c == nil || c.aead == nil {
		return "", errors.New("cipher is not initialized")
	}

	nonce := make([]byte, c.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	sealed := c.aead.Seal(nil, nonce, []byte(plain), nil)
	payload := append(nonce, sealed...)
	return base64.StdEncoding.EncodeToString(payload), nil
}

func (c *Cipher) Decrypt(payload string) (string, error) {
	if c == nil || c.aead == nil {
		return "", errors.New("cipher is not initialized")
	}

	raw, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return "", err
	}

	nonceSize := c.aead.NonceSize()
	if len(raw) < nonceSize {
		return "", errors.New("invalid encrypted payload")
	}

	nonce := raw[:nonceSize]
	ciphertext := raw[nonceSize:]
	plain, err := c.aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}
