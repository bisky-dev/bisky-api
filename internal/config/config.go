package config

import (
	"os"
)

func FromEnv() Config {
	return Config{
		Port:               getEnv("PORT", "8080"),
		DatabaseURL:        mustGetEnv("DATABASE_URL"),
		TokenEncryptionKey: firstNonEmpty(getEnv("TOKEN_ENCRYPTION_KEY", ""), getEnv("PAT_ENCRYPTION_KEY", ""), "dev-token-encryption-key"),
	}
}

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic("missing env: " + key)
	}
	return v
}
