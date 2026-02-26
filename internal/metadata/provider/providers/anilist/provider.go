package anilist

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/keithics/devops-dashboard/api/internal/metadata/provider"
	showmodel "github.com/keithics/devops-dashboard/api/internal/show"
)

const (
	defaultEndpoint = "https://graphql.anilist.co"
	defaultPage     = 1
	defaultPageSize = 10
	maxPageSize     = 50
	idPrefix        = "anilist:"
	searchQuery     = `query ($query: String!, $page: Int!, $perPage: Int!) {
  Page(page: $page, perPage: $perPage) {
    media(search: $query, type: ANIME, sort: SEARCH_MATCH) {
      id
      type
      status
      averageScore
      description(asHtml: false)
      bannerImage
      synonyms
      title {
        romaji
        english
        native
      }
    }
  }
}`
	showQuery = `query ($id: Int!) {
  Media(id: $id, type: ANIME) {
    id
    status
    description(asHtml: false)
    bannerImage
    synonyms
    coverImage {
      large
    }
    startDate {
      year
      month
      day
    }
    endDate {
      year
      month
      day
    }
    title {
      romaji
      english
      native
    }
  }
}`
	episodesQuery = `query ($mediaId: Int!, $page: Int!, $perPage: Int!) {
  Page(page: $page, perPage: $perPage) {
    airingSchedules(mediaId: $mediaId, sort: EPISODE) {
      episode
      airingAt
    }
  }
}`
)

func New() *Provider {
	return &Provider{
		endpoint: defaultEndpoint,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (p *Provider) Search(ctx context.Context, query string, opts metadata.SearchOpts) ([]metadata.SearchHit, error) {
	if strings.TrimSpace(query) == "" {
		return []metadata.SearchHit{}, nil
	}

	request := graphQLRequest{
		Query: searchQuery,
		Variables: map[string]any{
			"query":   query,
			"page":    normalizePage(opts.Page),
			"perPage": normalizeLimit(opts.Limit),
		},
	}

	var response graphQLSearchResponse
	if err := p.execute(ctx, request, &response); err != nil {
		return nil, err
	}

	hits := make([]metadata.SearchHit, 0, len(response.Data.Page.Media))
	for _, media := range response.Data.Page.Media {
		titlePreferred, titleOriginal := pickTitles(media.Title)
		typeValue := mapAniListType(media.Type)
		bannerURL := normalizeStringPtr(media.BannerImage)
		synopsis := normalizeStringPtr(media.Description)
		_ = normalizeAverageScore(media.AverageScore)

		hits = append(hits, metadata.SearchHit{
			Show: showmodel.Show{
				ExternalID:     formatExternalID(strconv.FormatInt(media.ID, 10)),
				TitlePreferred: titlePreferred,
				TitleOriginal:  titleOriginal,
				AltTitles:      buildAniListAltTitles(titlePreferred, titleOriginal, media.Title, media.Synonyms),
				Type:           typeValue,
				Status:         showmodel.NormalizeStatusOrDefault(media.Status, showmodel.StatusOngoing),
				Synopsis:       synopsis,
				BannerUrl:      bannerURL,
			},
		})
	}

	return hits, nil
}

func (p *Provider) GetShow(ctx context.Context, externalID string) (metadata.Show, error) {
	mediaID, err := parseExternalID(externalID)
	if err != nil {
		return metadata.Show{}, err
	}

	request := graphQLRequest{
		Query: showQuery,
		Variables: map[string]any{
			"id": mediaID,
		},
	}

	var response graphQLShowResponse
	if err := p.execute(ctx, request, &response); err != nil {
		return metadata.Show{}, err
	}

	if response.Data.Media.ID == 0 {
		return metadata.Show{}, fmt.Errorf("anilist show %s not found", externalID)
	}

	titlePreferred, titleOriginal := pickTitles(response.Data.Media.Title)
	startDate := normalizeAniListDate(response.Data.Media.StartDate)
	endDate := normalizeAniListDate(response.Data.Media.EndDate)
	synopsis := normalizeStringPtr(response.Data.Media.Description)
	posterURL := normalizeStringPtr(response.Data.Media.CoverImage.Large)
	bannerURL := normalizeStringPtr(response.Data.Media.BannerImage)

	return metadata.Show{
		Show: showmodel.Show{
			ExternalID:     formatExternalID(strconv.FormatInt(response.Data.Media.ID, 10)),
			TitlePreferred: titlePreferred,
			TitleOriginal:  titleOriginal,
			AltTitles:      buildAniListAltTitles(titlePreferred, titleOriginal, response.Data.Media.Title, response.Data.Media.Synonyms),
			Synopsis:       synopsis,
			StartDate:      startDate,
			EndDate:        endDate,
			PosterUrl:      posterURL,
			BannerUrl:      bannerURL,
			Type:           "anime",
			Status:         showmodel.NormalizeStatusOrDefault(response.Data.Media.Status, showmodel.StatusOngoing),
		},
	}, nil
}

func (p *Provider) ListEpisodes(ctx context.Context, externalID string, opts metadata.ListEpisodesOpts) ([]metadata.Episode, error) {
	mediaID, err := parseExternalID(externalID)
	if err != nil {
		return nil, err
	}

	request := graphQLRequest{
		Query: episodesQuery,
		Variables: map[string]any{
			"mediaId": mediaID,
			"page":    normalizePage(opts.Page),
			"perPage": normalizeLimit(opts.Limit),
		},
	}

	var response graphQLEpisodesResponse
	if err := p.execute(ctx, request, &response); err != nil {
		return nil, err
	}

	episodes := make([]metadata.Episode, 0, len(response.Data.Page.AiringSchedules))
	for _, item := range response.Data.Page.AiringSchedules {
		seasonNumber := int64(1)
		if opts.SeasonNumber != nil {
			seasonNumber = *opts.SeasonNumber
		}
		title := fmt.Sprintf("Episode %d", item.Episode)
		airDate := normalizeUnixDate(item.AiringAt)

		episodes = append(episodes, metadata.Episode{
			Provider:      metadata.ProviderAniList,
			ExternalID:    formatExternalID(fmt.Sprintf("%d:%d", mediaID, item.Episode)),
			SeasonNumber:  seasonNumber,
			EpisodeNumber: item.Episode,
			Title:         title,
			AirDate:       airDate,
		})
	}

	return episodes, nil
}

func (p *Provider) execute(ctx context.Context, request graphQLRequest, target any) error {
	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.endpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	httpRes, err := p.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()

	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return err
	}

	if httpRes.StatusCode < 200 || httpRes.StatusCode >= 300 {
		return fmt.Errorf("anilist request failed with status %d", httpRes.StatusCode)
	}

	if err := json.Unmarshal(body, target); err != nil {
		return err
	}

	switch value := target.(type) {
	case *graphQLSearchResponse:
		return firstGraphQLError(value.Errors)
	case *graphQLShowResponse:
		return firstGraphQLError(value.Errors)
	case *graphQLEpisodesResponse:
		return firstGraphQLError(value.Errors)
	default:
		return nil
	}
}

func firstGraphQLError(items []graphQLError) error {
	if len(items) == 0 {
		return nil
	}
	return fmt.Errorf("anilist graphql error: %s", items[0].Message)
}

func parseExternalID(value string) (int64, error) {
	normalized := strings.TrimSpace(value)
	normalized = strings.TrimPrefix(strings.ToLower(normalized), idPrefix)
	id, err := strconv.ParseInt(strings.TrimSpace(normalized), 10, 64)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("anilist external id must be a positive integer")
	}
	return id, nil
}

func formatExternalID(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if strings.HasPrefix(strings.ToLower(value), idPrefix) {
		return value
	}
	return idPrefix + value
}

func normalizePage(value int) int {
	if value < 1 {
		return defaultPage
	}
	return value
}

func normalizeLimit(value int) int {
	if value < 1 {
		return defaultPageSize
	}
	if value > maxPageSize {
		return maxPageSize
	}
	return value
}

func pickTitles(title anilistMediaTitle) (string, *string) {
	preferred := firstNonEmptyPtr(title.English, title.Romaji, title.Native)
	original := firstNonEmptyPtr(title.Native, title.Romaji, title.English)

	preferredValue := ""
	if preferred != nil {
		preferredValue = *preferred
	}
	if preferredValue == "" {
		preferredValue = "Untitled"
	}
	return preferredValue, normalizeStringPtr(original)
}

func firstNonEmptyPtr(values ...*string) *string {
	for _, value := range values {
		normalized := normalizeStringPtr(value)
		if normalized != nil {
			return normalized
		}
	}
	return nil
}

func normalizeStringPtr(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func normalizeAniListDate(value anilistDate) *string {
	if value.Year == nil || value.Month == nil || value.Day == nil {
		return nil
	}
	return ptrString(fmt.Sprintf("%04d-%02d-%02d", *value.Year, *value.Month, *value.Day))
}

func normalizeUnixDate(value *int64) *string {
	if value == nil || *value <= 0 {
		return nil
	}
	t := time.Unix(*value, 0).UTC()
	formatted := t.Format("2006-01-02")
	return &formatted
}

func mapAniListType(value string) string {
	if strings.EqualFold(value, "ANIME") {
		return "anime"
	}
	return strings.ToLower(strings.TrimSpace(value))
}

func normalizeAverageScore(value *float64) *float64 {
	if value == nil {
		return nil
	}
	normalized := *value / 100
	return &normalized
}

func ptrString(value string) *string {
	return &value
}

func buildAniListAltTitles(preferred string, original *string, title anilistMediaTitle, synonyms []string) []string {
	seen := map[string]struct{}{}
	add := func(value string, out *[]string) {
		normalized := strings.TrimSpace(value)
		if normalized == "" {
			return
		}
		key := strings.ToLower(normalized)
		if strings.EqualFold(normalized, preferred) {
			return
		}
		if original != nil && strings.EqualFold(normalized, *original) {
			return
		}
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		*out = append(*out, normalized)
	}

	out := make([]string, 0)
	if title.English != nil {
		add(*title.English, &out)
	}
	if title.Romaji != nil {
		add(*title.Romaji, &out)
	}
	if title.Native != nil {
		add(*title.Native, &out)
	}
	for _, item := range synonyms {
		add(item, &out)
	}
	return out
}
