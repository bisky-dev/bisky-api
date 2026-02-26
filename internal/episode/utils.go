package episode

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/keithics/devops-dashboard/api/internal/db/sqlc"
	"github.com/keithics/devops-dashboard/api/internal/httpx"
)

func normalizeCreateEpisodeRequest(req *createEpisodeRequest) {
	normalizeEpisodeFields(
		&req.ShowID,
		&req.Title,
		&req.AirDate,
	)
}

func normalizeUpdateEpisodeRequest(req *updateEpisodeRequest) {
	normalizeEpisodeFields(
		&req.ShowID,
		&req.Title,
		&req.AirDate,
	)
}

func normalizeEpisodeFields(showID *string, title *string, airDate **string) {
	*showID = strings.TrimSpace(*showID)
	*title = strings.TrimSpace(*title)
	*airDate = httpx.TrimmedOrNil(*airDate)
}

func validateCreateEpisodeRequest(req createEpisodeRequest) error {
	return validateEpisodePayload(req.ShowID, req.SeasonNumber, req.EpisodeNumber, req.Title, req.AirDate, req.RuntimeMinutes, req.ExternalIDs)
}

func validateUpdateEpisodeRequest(req updateEpisodeRequest) error {
	return validateEpisodePayload(req.ShowID, req.SeasonNumber, req.EpisodeNumber, req.Title, req.AirDate, req.RuntimeMinutes, req.ExternalIDs)
}

func validateEpisodePayload(showID string, seasonNumber, episodeNumber int64, title string, airDate *string, runtimeMinutes *int64, externalIDs externalIDs) error {
	if err := httpx.ValidateVar(showID, "required,uuid4", "showId is invalid"); err != nil {
		return err
	}
	if err := httpx.ValidateVar(seasonNumber, "gte=0", "seasonNumber is invalid"); err != nil {
		return err
	}
	if err := httpx.ValidateVar(episodeNumber, "gte=0", "episodeNumber is invalid"); err != nil {
		return err
	}
	if err := httpx.ValidateVar(title, "required,max=500", "title is invalid"); err != nil {
		return err
	}
	if err := httpx.ValidateOptionalDate(airDate, "airDate is invalid"); err != nil {
		return err
	}
	if runtimeMinutes != nil {
		if err := httpx.ValidateVar(*runtimeMinutes, "gte=0", "runtimeMinutes is invalid"); err != nil {
			return err
		}
	}
	if err := validateExternalIDs(externalIDs); err != nil {
		return err
	}
	return nil
}

func validateEpisodeID(id string) error {
	return httpx.ValidateVar(id, "required,uuid4", "internalEpisodeId is invalid")
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

func toEpisodeResponse(item sqlc.Episode) (episodeResponse, error) {
	externalIDs, err := unmarshalExternalIDs(item.ExternalIds)
	if err != nil {
		return episodeResponse{}, err
	}
	return episodeResponse{
		InternalEpisodeID: item.InternalEpisodeID,
		ShowID:            item.ShowID,
		SeasonNumber:      item.SeasonNumber,
		EpisodeNumber:     item.EpisodeNumber,
		Title:             item.Title,
		AirDate:           item.AirDate,
		RuntimeMinutes:    item.RuntimeMinutes,
		ExternalIDs:       externalIDs,
		CreatedAt:         item.CreatedAt,
		UpdatedAt:         item.UpdatedAt,
	}, nil
}
