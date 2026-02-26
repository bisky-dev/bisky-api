package anilist

import "net/http"

type Provider struct {
	endpoint string
	client   *http.Client
}

type graphQLRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

type graphQLError struct {
	Message string `json:"message"`
}

type graphQLSearchResponse struct {
	Data struct {
		Page struct {
			Media []anilistMediaSearch `json:"media"`
		} `json:"Page"`
	} `json:"data"`
	Errors []graphQLError `json:"errors"`
}

type graphQLShowResponse struct {
	Data struct {
		Media anilistMediaDetails `json:"Media"`
	} `json:"data"`
	Errors []graphQLError `json:"errors"`
}

type graphQLEpisodesResponse struct {
	Data struct {
		Page struct {
			AiringSchedules []anilistAiringSchedule `json:"airingSchedules"`
		} `json:"Page"`
	} `json:"data"`
	Errors []graphQLError `json:"errors"`
}

type anilistMediaSearch struct {
	ID           int64             `json:"id"`
	Type         string            `json:"type"`
	AverageScore *float64          `json:"averageScore"`
	Title        anilistMediaTitle `json:"title"`
}

type anilistMediaDetails struct {
	ID          int64             `json:"id"`
	Title       anilistMediaTitle `json:"title"`
	Description *string           `json:"description"`
	BannerImage *string           `json:"bannerImage"`
	CoverImage  struct {
		Large *string `json:"large"`
	} `json:"coverImage"`
	StartDate anilistDate `json:"startDate"`
	EndDate   anilistDate `json:"endDate"`
}

type anilistMediaTitle struct {
	Romaji  *string `json:"romaji"`
	English *string `json:"english"`
	Native  *string `json:"native"`
}

type anilistDate struct {
	Year  *int `json:"year"`
	Month *int `json:"month"`
	Day   *int `json:"day"`
}

type anilistAiringSchedule struct {
	Episode  int64  `json:"episode"`
	AiringAt *int64 `json:"airingAt"`
}
