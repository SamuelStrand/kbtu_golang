# Practice 5 (KBTU Go)

Implements:
- Pagination with `limit` & `offset`
- Dynamic filtering by **any** user field: `id`, `name`, `email`, `gender`, `birth_date`
- Dynamic sorting with `order_by` (+ default) and `order_dir`
- Many-to-Many table `user_friends` with constraint `user_id <> friend_id`
- Seed data: **20 users** + friendships (User01 & User02 have **3 common friends**)
- `GetCommonFriends` implemented with **ONE SQL query using JOIN** (no N+1)
- REST endpoints on port `:8080`

## Run

### Option A: with Docker Postgres
```bash
cp .env.example .env

docker compose up -d db

go mod tidy

go run ./cmd/api
```

### Option B: with your local Postgres
1. Create database `practice5`
2. Set `.env`
3. Run:
```bash
cp .env.example .env

go mod tidy

go run ./cmd/api
```

## Endpoints

### GET /users
Query params:
- `limit` (default 10, max 100)
- `offset` (default 0)
- `order_by` one of: `id|name|email|gender|birth_date` (default `id`)
- `order_dir` `asc|desc` (default `asc`)

Filters (all optional):
- `id` UUID
- `name` (ILIKE contains)
- `email` (ILIKE contains)
- `gender` (exact)
- `birth_date` (YYYY-MM-DD)

Example:
```bash
curl "http://localhost:8080/users?limit=5&offset=0&order_by=name&order_dir=asc&gender=male&email=@mail.com"
```

### GET /users/common-friends
Query params:
- `user1` UUID (required)
- `user2` UUID (required)

Example:
```bash
curl "http://localhost:8080/users/common-friends?user1=<UUID>&user2=<UUID>"
```

## Notes
- On startup the app creates tables if missing and seeds data if users < 20.
