package show

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/keithics/devops-dashboard/api/internal/db/sqlc"
	"github.com/keithics/devops-dashboard/api/internal/httpx"
)

func normalizeCreateShowRequest(req *createShowRequest) {
	normalizeShowFields(
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
	titlePreferred *string,
	titleOriginal **string,
	synopsis **string,
	startDate **string,
	endDate **string,
	posterURL **string,
	bannerURL **string,
	altTitles *[]string,
) {
	*titlePreferred = strings.TrimSpace(*titlePreferred)
	*titleOriginal = httpx.TrimmedOrNil(*titleOriginal)
	*synopsis = httpx.TrimmedOrNil(*synopsis)
	*startDate = httpx.TrimmedOrNil(*startDate)
	*endDate = httpx.TrimmedOrNil(*endDate)
	*posterURL = httpx.TrimmedOrNil(*posterURL)
	*bannerURL = httpx.TrimmedOrNil(*bannerURL)
	*altTitles = normalizeAltTitles(*altTitles)
}

func normalizeAltTitles(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}

	normalized := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		normalized = append(normalized, trimmed)
	}
	return normalized
}

func validateCreateShowRequest(req createShowRequest) error {
	return validateShowPayload(
		req.TitlePreferred,
		req.Type,
		req.Status,
		req.StartDate,
		req.EndDate,
		req.SeasonCount,
		req.EpisodeCount,
		req.ExternalIDs,
	)
}

func validateUpdateShowRequest(req updateShowRequest) error {
	return validateShowPayload(
		req.TitlePreferred,
		req.Type,
		req.Status,
		req.StartDate,
		req.EndDate,
		req.SeasonCount,
		req.EpisodeCount,
		req.ExternalIDs,
	)
}

func validateShowPayload(
	titlePreferred string,
	showType string,
	status string,
	startDate *string,
	endDate *string,
	seasonCount *int64,
	episodeCount *int64,
	externalIDs externalIDs,
) error {
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
	if err := validateExternalIDs(externalIDs); err != nil {
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

func validateExternalIDs(ids externalIDs) error {
	if ids.Anilist == nil && ids.Tvdb == nil {
		return errors.New("externalIds must include at least one provider id")
	}
	if ids.Anilist != nil {
		if err := httpx.ValidateVar(*ids.Anilist, "gt=0", "externalIds.anilist is invalid"); err != nil {
			return err
		}
	}
	if ids.Tvdb != nil {
		if err := httpx.ValidateVar(*ids.Tvdb, "gt=0", "externalIds.tvdb is invalid"); err != nil {
			return err
		}
	}
	return nil
}

func marshalExternalIDs(ids externalIDs) ([]byte, error) {
	return json.Marshal(ids)
}

func unmarshalExternalIDs(raw []byte) (externalIDs, error) {
	if len(raw) == 0 {
		return externalIDs{}, nil
	}
	var ids externalIDs
	if err := json.Unmarshal(raw, &ids); err != nil {
		return externalIDs{}, err
	}
	return ids, nil
}

func toShowResponse(show sqlc.Show) (showResponse, error) {
	externalIDs, err := unmarshalExternalIDs(show.ExternalIds)
	if err != nil {
		return showResponse{}, err
	}
	return showResponse{
		InternalShowID: show.InternalShowID,
		Show: Show{
			TitlePreferred: show.TitlePreferred,
			TitleOriginal:  show.TitleOriginal,
			AltTitles:      show.AltTitles,
			Type:           show.Type,
			Status:         show.Status,
			Synopsis:       show.Synopsis,
			StartDate:      show.StartDate,
			EndDate:        show.EndDate,
			PosterUrl:      show.PosterUrl,
			BannerUrl:      show.BannerUrl,
			SeasonCount:    show.SeasonCount,
			EpisodeCount:   show.EpisodeCount,
			ExternalIDs:    externalIDs,
		},
		CreatedAt: show.CreatedAt,
		UpdatedAt: show.UpdatedAt,
	}, nil
}
