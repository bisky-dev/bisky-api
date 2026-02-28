package normalize

import (
	"strconv"
	"strings"
)

func String(value string) string {
	return strings.TrimSpace(value)
}

func LowerString(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func StringPtr(value *string) *string {
	if value == nil {
		return nil
	}

	normalized := String(*value)
	if normalized == "" {
		return nil
	}

	return &normalized
}

func StringValuePtr(value string) *string {
	return StringPtr(&value)
}

func Strings(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}

	normalized := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := String(value)
		if trimmed == "" {
			continue
		}
		normalized = append(normalized, trimmed)
	}

	return normalized
}

func Page(value int, fallback int) int {
	if value < 1 {
		return fallback
	}
	return value
}

func Limit(value int, fallback int, max int) int {
	if value < 1 {
		return fallback
	}
	if max > 0 && value > max {
		return max
	}
	return value
}

func ParseBool(raw string, fallback bool) bool {
	if strings.TrimSpace(raw) == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(raw)
	if err != nil {
		return fallback
	}

	return parsed
}
