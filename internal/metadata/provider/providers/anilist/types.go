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
			Media []anilistMediaSummary `json:"media"`
		} `json:"Page"`
	} `json:"data"`
	Errors []graphQLError `json:"errors"`
}

type graphQLDiscoverResponse struct {
	Data struct {
		Trending struct {
			Media []anilistMediaSummary `json:"media"`
		} `json:"trending"`
		Popular struct {
			Media []anilistMediaSummary `json:"media"`
		} `json:"popular"`
		TopRated struct {
			Media []anilistMediaSummary `json:"media"`
		} `json:"topRated"`
		Upcoming struct {
			Media []anilistMediaSummary `json:"media"`
		} `json:"upcoming"`
		CurrentlyAiring struct {
			Media []anilistMediaSummary `json:"media"`
		} `json:"currentlyAiring"`
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

type anilistMediaSummary struct {
	ID          int64             `json:"id"`
	Type        string            `json:"type"`
	Status      string            `json:"status"`
	Description *string           `json:"description"`
	BannerImage *string           `json:"bannerImage"`
	Synonyms    []string          `json:"synonyms"`
	Episodes    *int64            `json:"episodes"`
	CoverImage  anilistCoverImage `json:"coverImage"`
	StartDate   anilistDate       `json:"startDate"`
	EndDate     anilistDate       `json:"endDate"`
	Title       anilistMediaTitle `json:"title"`
}

type anilistCoverImage struct {
	Large *string `json:"large"`
}

type anilistMediaDetails struct {
	ID          int64             `json:"id"`
	Title       anilistMediaTitle `json:"title"`
	Status      string            `json:"status"`
	Synonyms    []string          `json:"synonyms"`
	Description *string           `json:"description"`
	BannerImage *string           `json:"bannerImage"`
	CoverImage  anilistCoverImage `json:"coverImage"`
	StartDate   anilistDate       `json:"startDate"`
	EndDate     anilistDate       `json:"endDate"`
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
