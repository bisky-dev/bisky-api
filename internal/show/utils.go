package show

import (
	"encoding/json"
	"strings"

	"github.com/keithics/devops-dashboard/api/internal/db/sqlc"
	"github.com/keithics/devops-dashboard/api/internal/httpx"
	normalizeutil "github.com/keithics/devops-dashboard/api/internal/utils/normalize"
)

func normalizeCreateShowRequest(req *createShowRequest) {
	normalizeShowFields(
		&req.ExternalID,
		&req.TitlePreferred,
		&req.TitleOriginal,
		&req.Synopsis,
		&req.StartDate,
		&req.EndDate,
		&req.PosterUrl,
		&req.BannerUrl,
		&req.AltTitles,
	)
}

func normalizeUpdateShowRequest(req *updateShowRequest) {
	normalizeShowFields(
		&req.ExternalID,
		&req.TitlePreferred,
		&req.TitleOriginal,
		&req.Synopsis,
		&req.StartDate,
		&req.EndDate,
		&req.PosterUrl,
		&req.BannerUrl,
		&req.AltTitles,
	)
}

func normalizeShowFields(
	externalID *string,
	titlePreferred *string,
	titleOriginal **string,
	synopsis **string,
	startDate **string,
	endDate **string,
	posterURL **string,
	bannerURL **string,
	altTitles *[]string,
) {
	*externalID = normalizeutil.String(*externalID)
	*titlePreferred = normalizeutil.String(*titlePreferred)
	*titleOriginal = normalizeutil.StringPtr(*titleOriginal)
	*synopsis = normalizeutil.StringPtr(*synopsis)
	*startDate = normalizeutil.StringPtr(*startDate)
	*endDate = normalizeutil.StringPtr(*endDate)
	*posterURL = normalizeutil.StringPtr(*posterURL)
	*bannerURL = normalizeutil.StringPtr(*bannerURL)
	*altTitles = normalizeutil.Strings(*altTitles)
}

func validateCreateShowRequest(req createShowRequest) error {
	return validateShowPayload(
		req.ExternalID,
		req.TitlePreferred,
		req.Type,
		req.Status,
		req.StartDate,
		req.EndDate,
		req.SeasonCount,
		req.EpisodeCount,
	)
}

func validateUpdateShowRequest(req updateShowRequest) error {
	return validateShowPayload(
		req.ExternalID,
		req.TitlePreferred,
		req.Type,
		req.Status,
		req.StartDate,
		req.EndDate,
		req.SeasonCount,
		req.EpisodeCount,
	)
}

func validateShowPayload(
	externalID string,
	titlePreferred string,
	showType string,
	status string,
	startDate *string,
	endDate *string,
	seasonCount *int64,
	episodeCount *int64,
) error {
	if externalID != "" {
		if err := httpx.ValidateVar(externalID, "max=128", "externalId is invalid"); err != nil {
			return err
		}
	}
	if err := httpx.ValidateVar(titlePreferred, "required,max=500", "titlePreferred is invalid"); err != nil {
		return err
	}
	if err := httpx.ValidateVar(showType, "required,oneof=anime tv movie ova special", "type is invalid"); err != nil {
		return err
	}
	if err := httpx.ValidateVar(status, "required,oneof=ongoing finished", "status is invalid"); err != nil {
		return err
	}
	if err := httpx.ValidateOptionalDate(startDate, "startDate is invalid"); err != nil {
		return err
	}
	if err := httpx.ValidateOptionalDate(endDate, "endDate is invalid"); err != nil {
		return err
	}
	if err := validateOptionalInt64(seasonCount, "gte=0", "seasonCount is invalid"); err != nil {
		return err
	}
	if err := validateOptionalInt64(episodeCount, "gte=0", "episodeCount is invalid"); err != nil {
		return err
	}
	return nil
}

func validateShowID(id string) error {
	return httpx.ValidateVar(id, "required,uuid4", "internalShowId is invalid")
}

func validateOptionalInt64(value *int64, rule string, message string) error {
	if value == nil {
		return nil
	}
	return httpx.ValidateVar(*value, rule, message)
}

type externalIDPayload struct {
	ExternalID string `json:"externalId,omitempty"`
}

func marshalExternalID(externalID string) ([]byte, error) {
	externalID = strings.TrimSpace(externalID)
	if externalID == "" {
		return []byte("{}"), nil
	}
	return json.Marshal(externalIDPayload{ExternalID: externalID})
}

func unmarshalExternalID(raw []byte) (string, error) {
	if len(raw) == 0 {
		return "", nil
	}
	var payload externalIDPayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return "", err
	}
	return strings.TrimSpace(payload.ExternalID), nil
}

func toShowResponse(show sqlc.Show) (showResponse, error) {
	externalID, err := unmarshalExternalID(show.ExternalIds)
	if err != nil {
		return showResponse{}, err
	}
	return showResponse{
		InternalShowID: show.InternalShowID,
		Show: Show{
			ExternalID:     externalID,
			TitlePreferred: show.TitlePreferred,
			TitleOriginal:  show.TitleOriginal,
			AltTitles:      normalizeutil.Strings(show.AltTitles),
			Type:           show.Type,
			Status:         show.Status,
			Synopsis:       show.Synopsis,
			StartDate:      show.StartDate,
			EndDate:        show.EndDate,
			PosterUrl:      show.PosterUrl,
			BannerUrl:      show.BannerUrl,
			SeasonCount:    show.SeasonCount,
			EpisodeCount:   show.EpisodeCount,
		},
		CreatedAt: show.CreatedAt,
		UpdatedAt: show.UpdatedAt,
	}, nil
}
