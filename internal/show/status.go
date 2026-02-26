package show

import "strings"

const (
	StatusOngoing  = "ongoing"
	StatusFinished = "finished"
)

func NormalizeStatus(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))

	switch normalized {
	case StatusOngoing:
		return StatusOngoing
	case StatusFinished:
		return StatusFinished
	case "releasing", "continuing", "upcoming", "in production", "returning series", "current series":
		return StatusOngoing
	case "ended", "completed", "cancelled", "canceled":
		return StatusFinished
	default:
		return ""
	}
}

func NormalizeStatusOrDefault(value string, fallback string) string {
	normalized := NormalizeStatus(value)
	if normalized == "" {
		return fallback
	}
	return normalized
}
