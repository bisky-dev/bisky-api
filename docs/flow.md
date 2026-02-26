# Metadata Add Show Flow

## 1) Search metadata

- Endpoint: `GET /metadata/search?query={q}&type={provider}&page=1&limit=10`
- Provider is selected by `type` (`anidb`, `anilist`, `tvdb`).
- Response items include `externalId` (provider-prefixed), plus show fields used by add-show.

## 2) Add show from search result

- Endpoint: `POST /metadata/show`
- Request body uses the same shape as a search item (show payload), including `externalId`.
- `externalId` must be prefixed (`anidb:...`, `anilist:...`, `tvdb:...`).

## 3) Service orchestration

- `metadata.AddShow` derives provider from `externalId` prefix.
- It calls metadata worker `ListEpisodes(...)` for that provider/external id.
- It maps fetched episodes into enqueue input.

## 4) Enqueue transaction

- Runs in one DB transaction:
1. Checks for existing `pending` row in `job_shows` for the same `externalId` (via `shows.external_ids->>'externalId'`).
2. If found, returns the existing row (no new insert).
3. If not found:
   - Inserts into `shows`.
   - Inserts fetched episodes into `episodes` linked by `show_id` (`ON CONFLICT (show_id, season_number, episode_number) DO NOTHING`).
   - Inserts one `pending` row into `job_shows`.

## 5) Response

- Returns:
  - `internalShowId`
  - `internalJobShowId`
  - `status`
  - `retryCount`
