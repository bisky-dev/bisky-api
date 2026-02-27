package show

import (
	normalizeutil "github.com/keithics/devops-dashboard/api/internal/utils/normalize"
)

const (
	StatusOngoing  = "ongoing"
	StatusFinished = "finished"
)

func NormalizeStatus(value string) string {
	normalized := normalizeutil.LowerString(value)

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
