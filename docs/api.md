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
