package httpx

import (
	"strconv"
	"strings"
)

func ParsePositiveInt(raw string, fallback int) int {
	if strings.TrimSpace(raw) == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(raw)
	if err != nil || parsed < 1 {
		return fallback
	}

	return parsed
}
