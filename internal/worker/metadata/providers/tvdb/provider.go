package tvdb

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	showmodel "github.com/keithics/devops-dashboard/api/internal/show"
	"github.com/keithics/devops-dashboard/api/internal/worker/metadata"
)

const (
	defaultBaseURL  = "https://api4.thetvdb.com/v4"
	defaultTimeout  = 20 * time.Second
	defaultPage     = 1
	defaultPageSize = 10
	maxPageSize     = 100
)

func New() *Provider {
	return &Provider{
		baseURL: strings.TrimRight(getEnv("TVDB_BASE_URL", defaultBaseURL), "/"),
		apiKey:  strings.TrimSpace(os.Getenv("TVDB_API_KEY")),
		pin:     strings.TrimSpace(os.Getenv("TVDB_PIN")),
		client: &http.Client{
			Timeout: defaultTimeout,
		},
		debug: strings.EqualFold(strings.TrimSpace(os.Getenv("TVDB_DEBUG")), "1") ||
			strings.EqualFold(strings.TrimSpace(os.Getenv("TVDB_DEBUG")), "true"),
	}
}

func (p *Provider) Search(ctx context.Context, query string, opts metadata.SearchOpts) ([]metadata.SearchHit, error) {
	if strings.TrimSpace(query) == "" {
		return []metadata.SearchHit{}, nil
	}

	items, err := p.searchItems(ctx, query, opts, true)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		// TVDB's search behavior can vary by account/index; retry without type filter.
		items, err = p.searchItems(ctx, query, opts, false)
		if err != nil {
			return nil, err
		}
	}

	hits := make([]metadata.SearchHit, 0, len(items))
	for _, item := range items {
		titlePreferred := firstString(item, "name", "name_translated", "seriesName", "title", "slug")
		if titlePreferred == "" {
			titlePreferred = "Untitled"
		}

		titleOriginal := stringPtr(firstString(item, "name"))
		bannerURL := stringPtr(firstString(item, "image_url", "thumbnail", "image", "banner"))
		synopsis := stringPtr(firstString(item, "overview", "overview_translated", "summary", "plot"))
		_ = floatPtr(firstFloat64(item, "score"))

		hits = append(hits, metadata.SearchHit{
			Show: showmodel.Show{
				TitlePreferred: titlePreferred,
				TitleOriginal:  titleOriginal,
				Type:           mapShowType(firstString(item, "primary_type", "type")),
				Synopsis:       synopsis,
				BannerUrl:      bannerURL,
			},
		})
	}

	return hits, nil
}

func (p *Provider) searchItems(ctx context.Context, query string, opts metadata.SearchOpts, includeType bool) ([]map[string]any, error) {
	// TODO cleanup: remove temporary debug logging and search fallback once TVDB response handling is fully stabilized.
	page := normalizePage(opts.Page)
	tvdbPage := toTVDBPage(page)

	params := url.Values{}
	params.Set("query", query)
	if includeType {
		params.Set("type", "series")
	}
	params.Set("page", strconv.Itoa(tvdbPage))
	params.Set("limit", strconv.Itoa(normalizeLimit(opts.Limit)))

	body, err := p.doRequest(ctx, http.MethodGet, "/search?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	if p.debug {
		log.Printf("[tvdb] search includeType=%v query=%q page=%d tvdbPage=%d limit=%d raw=%s", includeType, query, page, tvdbPage, normalizeLimit(opts.Limit), truncateForLog(body, 1200))
	}
	items, err := decodeArrayData(body)
	if err != nil {
		return nil, err
	}
	if p.debug {
		log.Printf("[tvdb] search parsed includeType=%v count=%d", includeType, len(items))
	}
	return items, nil
}

func (p *Provider) GetShow(ctx context.Context, externalID string) (metadata.Show, error) {
	id, err := parseExternalID(externalID)
	if err != nil {
		return metadata.Show{}, err
	}

	body, err := p.doRequest(ctx, http.MethodGet, "/series/"+id+"/extended", nil)
	if err != nil {
		return metadata.Show{}, err
	}

	item, err := decodeObjectData(body)
	if err != nil {
		return metadata.Show{}, err
	}

	titlePreferred := firstString(item, "name", "slug")
	if titlePreferred == "" {
		titlePreferred = "Untitled"
	}

	posterURL := stringPtr(firstString(item, "image", "image_url"))
	bannerURL := posterURL
	if artworks := mapSlice(item, "artworks"); len(artworks) > 0 {
		for _, artwork := range artworks {
			typ := strings.ToLower(firstString(artwork, "type", "typeName"))
			if strings.Contains(typ, "banner") {
				bannerURL = stringPtr(firstString(artwork, "image", "image_url"))
				break
			}
		}
	}

	return metadata.Show{
		Show: showmodel.Show{
			TitlePreferred: titlePreferred,
			TitleOriginal:  stringPtr(firstString(item, "name")),
			Type:           mapShowType(firstString(item, "type", "primary_type")),
			Synopsis:       stringPtr(firstString(item, "overview")),
			StartDate:      stringPtr(firstString(item, "firstAired", "first_air_time")),
			EndDate:        stringPtr(firstString(item, "lastAired")),
			PosterUrl:      posterURL,
			BannerUrl:      bannerURL,
		},
	}, nil
}

func (p *Provider) ListEpisodes(ctx context.Context, externalID string, opts metadata.ListEpisodesOpts) ([]metadata.Episode, error) {
	id, err := parseExternalID(externalID)
	if err != nil {
		return nil, err
	}

	page := normalizePage(opts.Page)
	tvdbPage := toTVDBPage(page)

	params := url.Values{}
	params.Set("page", strconv.Itoa(tvdbPage))
	params.Set("limit", strconv.Itoa(normalizeLimit(opts.Limit)))
	body, err := p.doRequest(ctx, http.MethodGet, "/series/"+id+"/episodes/default?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	items, err := decodeArrayData(body)
	if err != nil {
		return nil, err
	}

	episodes := make([]metadata.Episode, 0, len(items))
	for _, item := range items {
		seasonNumber := int64(firstInt(item, "seasonNumber", "season"))
		episodeNumber := int64(firstInt(item, "number", "episodeNumber"))
		if seasonNumber < 0 || episodeNumber < 0 {
			continue
		}
		if opts.SeasonNumber != nil && seasonNumber != *opts.SeasonNumber {
			continue
		}

		title := firstString(item, "name")
		if title == "" {
			title = fmt.Sprintf("Episode %d", episodeNumber)
		}

		epID := firstString(item, "id")
		if epID == "" {
			epID = fmt.Sprintf("%s:%d:%d", id, seasonNumber, episodeNumber)
		}

		runtime := int64Ptr(firstInt(item, "runtime", "runtimeMinutes"))

		episodes = append(episodes, metadata.Episode{
			Provider:       metadata.ProviderTVDB,
			ExternalID:     epID,
			SeasonNumber:   seasonNumber,
			EpisodeNumber:  episodeNumber,
			Title:          title,
			AirDate:        stringPtr(firstString(item, "aired", "firstAired")),
			RuntimeMinutes: runtime,
		})
	}

	sort.Slice(episodes, func(i, j int) bool {
		if episodes[i].SeasonNumber != episodes[j].SeasonNumber {
			return episodes[i].SeasonNumber < episodes[j].SeasonNumber
		}
		return episodes[i].EpisodeNumber < episodes[j].EpisodeNumber
	})

	return episodes, nil
}

func (p *Provider) doRequest(ctx context.Context, method, path string, payload any) ([]byte, error) {
	token, err := p.ensureToken(ctx)
	if err != nil {
		return nil, err
	}

	var bodyReader io.Reader
	if payload != nil {
		raw, marshalErr := json.Marshal(payload)
		if marshalErr != nil {
			return nil, marshalErr
		}
		bodyReader = bytes.NewReader(raw)
	}

	req, err := http.NewRequestWithContext(ctx, method, p.baseURL+path, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("tvdb request failed with status %d", res.StatusCode)
	}
	return body, nil
}

func (p *Provider) ensureToken(ctx context.Context) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.token != "" && time.Now().Before(p.tokenUntil) {
		return p.token, nil
	}
	if p.apiKey == "" {
		return "", errors.New("tvdb provider is not configured: missing TVDB_API_KEY")
	}

	loginPayload := map[string]string{"apikey": p.apiKey}
	if p.pin != "" {
		loginPayload["pin"] = p.pin
	}

	raw, err := json.Marshal(loginPayload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/login", bytes.NewReader(raw))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	res, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", fmt.Errorf("tvdb login failed with status %d", res.StatusCode)
	}

	token, err := extractToken(body)
	if err != nil {
		return "", err
	}

	p.token = token
	p.tokenUntil = time.Now().Add(24 * time.Hour)
	return p.token, nil
}

func extractToken(body []byte) (string, error) {
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", err
	}

	if token := firstString(payload, "token"); token != "" {
		return token, nil
	}
	if data, ok := payload["data"].(map[string]any); ok {
		if token := firstString(data, "token"); token != "" {
			return token, nil
		}
	}

	return "", errors.New("tvdb login response missing token")
}

func decodeArrayData(body []byte) ([]map[string]any, error) {
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	data := payload["data"]
	switch v := data.(type) {
	case []any:
		return anySliceToMapSlice(v), nil
	case map[string]any:
		for _, key := range []string{"series", "episodes", "results", "items", "hits"} {
			if arr, ok := v[key].([]any); ok {
				return anySliceToMapSlice(arr), nil
			}
		}

		// TVDB responses vary by endpoint/version; recursively find first object array.
		if arr := findFirstMapArray(v); len(arr) > 0 {
			return arr, nil
		}
	}

	// Fallback: search entire payload for a nested object array.
	if arr := findFirstMapArray(payload); len(arr) > 0 {
		return arr, nil
	}

	return []map[string]any{}, nil
}

func decodeObjectData(body []byte) (map[string]any, error) {
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	if data, ok := payload["data"].(map[string]any); ok {
		return data, nil
	}

	return map[string]any{}, nil
}

func anySliceToMapSlice(items []any) []map[string]any {
	out := make([]map[string]any, 0, len(items))
	for _, item := range items {
		if m, ok := item.(map[string]any); ok {
			out = append(out, m)
		}
	}
	return out
}

func mapSlice(m map[string]any, key string) []map[string]any {
	raw, ok := m[key].([]any)
	if !ok {
		return nil
	}
	return anySliceToMapSlice(raw)
}

func findFirstMapArray(value any) []map[string]any {
	switch v := value.(type) {
	case []any:
		items := anySliceToMapSlice(v)
		if len(items) > 0 {
			return items
		}
		for _, item := range v {
			if nested := findFirstMapArray(item); len(nested) > 0 {
				return nested
			}
		}
	case map[string]any:
		for _, item := range v {
			if nested := findFirstMapArray(item); len(nested) > 0 {
				return nested
			}
		}
	}
	return nil
}

func firstString(m map[string]any, keys ...string) string {
	for _, key := range keys {
		value, ok := m[key]
		if !ok || value == nil {
			continue
		}

		switch v := value.(type) {
		case string:
			v = strings.TrimSpace(v)
			if v != "" {
				return v
			}
		case float64:
			if v > 0 {
				return strconv.FormatInt(int64(v), 10)
			}
		case int64:
			if v > 0 {
				return strconv.FormatInt(v, 10)
			}
		case map[string]any:
			for _, nestedKey := range []string{"value", "id", "name", "title"} {
				if nested, ok := v[nestedKey]; ok {
					if nestedString, ok := nested.(string); ok {
						nestedString = strings.TrimSpace(nestedString)
						if nestedString != "" {
							return nestedString
						}
					}
				}
			}
		}
	}
	return ""
}

func firstFloat64(m map[string]any, keys ...string) float64 {
	for _, key := range keys {
		value, ok := m[key]
		if !ok || value == nil {
			continue
		}
		switch v := value.(type) {
		case float64:
			return v
		case int64:
			return float64(v)
		case string:
			if parsed, err := strconv.ParseFloat(strings.TrimSpace(v), 64); err == nil {
				return parsed
			}
		}
	}
	return 0
}

func firstInt(m map[string]any, keys ...string) int {
	for _, key := range keys {
		value, ok := m[key]
		if !ok || value == nil {
			continue
		}
		switch v := value.(type) {
		case float64:
			return int(v)
		case int64:
			return int(v)
		case string:
			parsed, err := strconv.Atoi(strings.TrimSpace(v))
			if err == nil {
				return parsed
			}
		}
	}
	return 0
}

func parseExternalID(externalID string) (string, error) {
	id := normalizeSeriesID(strings.TrimSpace(externalID))
	if id == "" {
		return "", errors.New("tvdb external id is required")
	}
	if _, err := strconv.ParseInt(id, 10, 64); err != nil {
		return "", errors.New("tvdb external id must be a positive integer")
	}
	return id, nil
}

func normalizeSeriesID(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	// TVDB search frequently returns prefixed IDs like "series-121361".
	if strings.HasPrefix(strings.ToLower(value), "series-") {
		return strings.TrimSpace(value[len("series-"):])
	}
	return value
}

func mapShowType(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	switch normalized {
	case "anime":
		return "anime"
	case "movie":
		return "movie"
	case "ova":
		return "ova"
	case "special":
		return "special"
	default:
		return "tv"
	}
}

func floatPtr(value float64) *float64 {
	if value <= 0 {
		return nil
	}
	return &value
}

func int64Ptr(value int) *int64 {
	if value <= 0 {
		return nil
	}
	converted := int64(value)
	return &converted
}

func stringPtr(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func normalizePage(value int) int {
	if value < 1 {
		return defaultPage
	}
	return value
}

func toTVDBPage(value int) int {
	page := normalizePage(value)
	if page <= 1 {
		return 0
	}
	return page - 1
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

func getEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func truncateForLog(body []byte, max int) string {
	if len(body) <= max {
		return string(body)
	}
	return string(body[:max]) + "...(truncated)"
}
