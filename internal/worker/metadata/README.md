# Worker Metadata Module

This package defines provider-agnostic metadata contracts used by background workers.

Provider interface:
- Search(ctx, query, opts) -> []SearchHit
- GetShow(ctx, externalID) -> Show
- ListEpisodes(ctx, externalID, opts) -> []Episode

`externalID` is provider-specific:
- AniList: media ID
- TVDB: series ID

Provider adapters:
- `providers/anilist`
- `providers/tvdb`

Current adapter methods are scaffold stubs and return not implemented errors.
