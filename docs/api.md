# API Documentation

Base URL (local): `http://localhost:3000`

## Conventions

- Content type: `application/json`
- Auth token format from login: `Bearer <access_token>`
- Error response shape:

```json
{
  "error": {
    "code": "SOME_CODE",
    "message": "Human readable message"
  }
}
```

---

## Health

### `GET /health`

Liveness check.

Success response (`200`):

```json
{
  "ok": true
}
```

---

## Auth

### `POST /auth/register`

Create a new user account.

Request body:

```json
{
  "email": "user@example.com",
  "password": "strong-password-123"
}
```

Success response (`201`):

```json
{
  "user": {
    "id": "uuid-or-text-id",
    "email": "user@example.com",
    "created_at": "2026-02-26T16:00:00Z"
  }
}
```

Possible errors:
- `400` invalid request body
- `409` email already exists
- `500` internal error

### `POST /auth/login`

Authenticate a user.

Request body:

```json
{
  "email": "user@example.com",
  "password": "strong-password-123"
}
```

Success response (`200`):

```json
{
  "access_token": "token",
  "token_type": "Bearer",
  "expires_in": 3600,
  "user": {
    "id": "uuid-or-text-id",
    "email": "user@example.com",
    "created_at": "2026-02-26T16:00:00Z"
  }
}
```

Possible errors:
- `400` invalid request body
- `401` invalid email or password
- `500` internal error

### `POST /auth/forgot-password`

Accepts an email and returns a generic success response to avoid account enumeration.

Request body:

```json
{
  "email": "user@example.com"
}
```

Success response (`202`):

```json
{
  "message": "If an account exists for that email, a reset link will be sent."
}
```

Possible errors:
- `400` invalid request body
- `500` internal error

---

## Shows

### `POST /shows`

Create a show.

Request body:

```json
{
  "titlePreferred": "Frieren: Beyond Journey's End",
  "titleOriginal": "Sousou no Frieren",
  "altTitles": ["Frieren"],
  "type": "anime",
  "status": "ongoing",
  "synopsis": "An elf mage reflects on life and time.",
  "startDate": "2023-09-29",
  "endDate": null,
  "posterUrl": "https://example.com/poster.jpg",
  "bannerUrl": "https://example.com/banner.jpg",
  "seasonCount": 1,
  "episodeCount": 28,
  "externalIds": {
    "anilist": 154587,
    "tvdb": 420000
  }
}
```

Success response (`201`): show object.

### `GET /shows`

List shows.

Success response (`200`): array of show objects.

### `GET /shows/{internalShowId}`

Get one show by UUID.

Success response (`200`): show object.

### `PUT /shows/{internalShowId}`

Update one show by UUID (same body shape as create).

Success response (`200`): updated show object.

### `DELETE /shows/{internalShowId}`

Delete one show by UUID.

Success response (`204`): no body.

---

## Episodes

### `POST /episodes`

Create an episode.

Request body:

```json
{
  "showId": "3cb5e44c-9cb6-4eb1-b34d-9c57e513c127",
  "seasonNumber": 1,
  "episodeNumber": 1,
  "title": "Journey's End",
  "airDate": "2023-09-29",
  "runtimeMinutes": 24,
  "externalIds": {
    "anilist": 1,
    "tvdb": 2
  }
}
```

Success response (`201`): episode object.

### `GET /episodes`

List episodes.

Success response (`200`): array of episode objects.

### `GET /episodes/{internalEpisodeId}`

Get one episode by UUID.

Success response (`200`): episode object.

### `PUT /episodes/{internalEpisodeId}`

Update one episode by UUID (same body shape as create).

Success response (`200`): updated episode object.

### `DELETE /episodes/{internalEpisodeId}`

Delete one episode by UUID.

Success response (`204`): no body.

---

## Metadata

Provider type query param:
- `type=anidb` (default)
- `type=tvdb`

### `GET /metadata/search?query={q}&type={type}`

Search provider metadata by text query.

### `POST /metadata/show`

Create a show using the same JSON shape returned by `GET /metadata/search`, then enqueue a job linked to that show.
