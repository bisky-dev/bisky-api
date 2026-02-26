# TODO

- [ ] Remove temporary TVDB debug mode and fallback search path after payload handling is stable.
- [ ] Revisit TVDB response parsing and replace generic map parsing with typed response structs.
- [ ] Decide final API contract for `/metadata/search` and `/metadata/show` after aligning to `show.Show`.
- [ ] Update Swagger docs to match current metadata response shapes.
- [ ] Add integration tests for AniList and TVDB provider adapters.
- [ ] Merge AniList episode schedules with TVDB episode titles so AniList episode responses use real titles instead of fallback `Episode N`.
