# Space Bounty Hunters API

REST API for space bounty hunters: a database of wanted criminals plus a journal of hunt attempts.

**Stack:** Go (`net/http`, `pgx`) + PostgreSQL 16, run with Docker Compose

## Run

```bash
cp .env.example .env
docker compose up --build
```

API: http://localhost:8080 . Postgres is also exposed on `localhost:5432`.

## Database setup / import

Docker Compose does it automatically: on first start with an empty volume Postgres executes
[db/schema.sql](db/schema.sql) (tables) and [db/seed.sql](db/seed.sql) (7 criminals, 12 hunt logs).

To reset the data: `docker compose down -v && docker compose up --build`.

Manual import into any PostgreSQL instance:

```bash
createdb bounty
psql bounty -f db/schema.sql
psql bounty -f db/seed.sql
```

then start the API locally (needs Go 1.26):

```bash
POSTGRES_USER=... POSTGRES_PASSWORD=... POSTGRES_DB=bounty DB_HOST=localhost go run ./cmd/api
```

## Models

**wanted_criminals** — `id`, `name` (nullable), `nickname`, `reward`, `latitude`/`longitude` (nullable, both or neither),
`last_known_note` (text alternative to coordinates), `difficulty` (1–10), `armed`, `has_crew`, `escape_master`.

**hunt_logs** — `id`, `criminal_id` (FK → wanted_criminals, `ON DELETE CASCADE`), `success`, `notes`,
`created_at` (set by the database). One criminal → many hunt logs.

## Endpoints

Trailing slash is optional.

| Method | Path | Description |
|---|---|---|
| GET | `/api/criminals/` | List. Query: `min_reward`, `max_reward`, `armed`, `has_crew`, `escape_master`, `sort` (`reward` / `-reward`). Filters combine with AND. |
| GET | `/api/criminals/{id}/` | One criminal |
| GET | `/api/hunts/` | List hunt logs. Query: `criminal_id` |
| GET | `/api/hunts/{id}/` | One hunt log |
| POST | `/api/hunts/` | Create. Body: `criminal_id`, `success` (required), `notes`. `created_at` is server-generated; any client value is ignored |
| PATCH | `/api/hunts/{id}/` | Update `success` and/or `notes` |
| DELETE | `/api/hunts/{id}/` | Delete (204) |

Example: `GET /api/criminals/?min_reward=5000&armed=true&sort=-reward`

### Status codes

`200` OK · `201` created · `204` deleted · `400` invalid id / query param / body / non-existent `criminal_id` · `404` resource not found.
Errors are JSON: `{"error": "..."}`.

## Postman

Import [postman/space-bounty-hunters.postman_collection.json](postman/space-bounty-hunters.postman_collection.json)
and run the collection in order (variable `base_url` = `http://localhost:8080`). It covers listing, filtering, sorting,
details, hunt log create/read/update/delete and several invalid requests, each with test assertions.

CLI run: `newman run postman/space-bounty-hunters.postman_collection.json`
