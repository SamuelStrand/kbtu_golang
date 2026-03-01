# Practice 3 (Go + Postgres)

## What is implemented
- Postgres connection + auto migrations (golang-migrate)
- Repository -> Usecase -> Handler layers
- CRUD endpoints for users (JSON only)
- Middleware:
  - request/response logging (standard `log`)
  - X-API-KEY authentication
- Healthcheck endpoint
- Optional EASY: config from `.env`, swagger docs at `/swagger`
- Optional MEDIUM:
  - soft delete (`deleted_at`)
  - pagination (`limit`, `offset`) for GET `/users`
  - transaction support in `CreateUser` (creates user + audit log)

## Run
1) Create DB (example: `mydb`) in Postgres.
2) Copy env:

```bash
cp .env.example .env
```

3) Start server:

```bash
go mod tidy
go run ./cmd/api
```

Server starts on `:8080`.

## Auth
All `/users` endpoints require header:

- `X-API-KEY: dev-api-key`

(You can change it in `.env`)

## Endpoints
- `GET /health`
- `GET /swagger` (UI) and `GET /swagger/swagger.json` (spec)

Protected:
- `GET /users?limit=20&offset=0`
- `GET /users/{id}`
- `POST /users`
- `PUT /users/{id}`
- `PATCH /users/{id}`
- `DELETE /users/{id}` (soft delete)
