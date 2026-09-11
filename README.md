# Akademi API

Backend REST API for **Akademi**, a school administration platform. Built with Go and PostgreSQL, following a layered architecture (handler → service → repository → model) with JWT authentication and role-based access control.

**Live API:** [akademi-api.onrender.com](https://akademi-api.onrender.com)
**Frontend repo:** [akademi](https://github.com/TakedaB/akademi)

> **Note:** deployed on Render's free tier. The Postgres instance is a temporary free database that expires periodically — if the API is unresponsive, it may need to be re-provisioned.

## Features

- JWT authentication with role-based claims
- Role-based access control (RBAC) across 4 roles: `diretoria` (school board), `financeiro` (finance staff), `professor` (teacher), `aluno` (student)
- Transactional user creation — creating a Teacher or Student also creates their login account (`users` + `teachers`/`students` in a single DB transaction)
- Auto-generated student enrollment numbers (year + atomic sequential counter, e.g. `2026004`)
- In-memory rate limiting on login (5 attempts / 15 minutes per IP)
- CORS configured for both local development and Vercel preview/production deployments

## Tech Stack

- **Language:** Go 1.27
- **Database:** PostgreSQL (via `lib/pq`)
- **Auth:** JWT (`golang-jwt/jwt/v5`), bcrypt password hashing
- **Routing:** standard library `net/http` (`http.ServeMux` with Go 1.22+ method-based routing)
- **Hosting:** Render (Web Service + managed Postgres)

## Architecture

```
cmd/api/            entry point, route registration
internal/
  handler/           HTTP handlers (request/response, no business logic)
  service/           business logic, validation, orchestration
  repository/        SQL queries, no business logic
  model/              structs shared across layers
  middleware/         auth (JWT), RBAC, CORS, rate limiting
```

Each domain (Students, Teachers, Finance, Auth) follows the same four-layer pattern. Teachers and Students are modeled as extensions of `users` (one-to-one via `user_id` FK) — both roles log in with their own credentials. Finance records belong to a Student.

## Role-Based Access Control

| Resource                   | Aluno | Professor | Financeiro | Diretoria |
| -------------------------- | :---: | :-------: | :--------: | :-------: |
| View Students              |  ❌   |    ✅     |     ✅     |    ✅     |
| Create/Edit/Delete Student |  ❌   |    ❌     |     ✅     |    ✅     |
| View Teachers              |  ✅   |    ✅     |     ✅     |    ✅     |
| Create/Delete Teacher      |  ❌   |    ❌     |     ❌     |    ✅     |
| View/Manage Finance        |  ❌   |    ❌     |     ✅     |    ✅     |
| View own profile (`/me`)   |  ✅   |    ✅     |     ✅     |    ✅     |

## API Endpoints

| Method                | Path                     | Auth                     | Description                              |
| --------------------- | ------------------------ | ------------------------ | ---------------------------------------- |
| POST                  | `/login`                 | Public (rate-limited)    | Returns JWT                              |
| GET                   | `/me`                    | Any authenticated user   | Current user's profile                   |
| GET/POST/PUT/DELETE   | `/students`              | Role-gated               | Student CRUD (creation is transactional) |
| GET/POST/DELETE       | `/teachers`              | Role-gated               | Teacher CRUD (creation is transactional) |
| GET/POST/PATCH/DELETE | `/finance`               | `diretoria`/`financeiro` | Charges CRUD, status updates             |
| GET                   | `/students/{id}/finance` | `diretoria`/`financeiro` | Charges for a specific student           |
| GET                   | `/health`                | Public                   | Health check                             |

## Local Setup

**Requirements:** Go 1.27+, PostgreSQL, Docker (optional, for local Postgres)

```bash
git clone https://github.com/TakedaB/akademi-api.git
cd akademi-api

# start a local Postgres (or use your own)
docker run --name akademi-postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres

# apply the schema
docker exec -i akademi-postgres psql -U postgres -d postgres -c "CREATE DATABASE akademi;"
docker exec -i akademi-postgres psql -U postgres -d akademi < schema.sql

# configure environment
cp .env.example .env   # set JWT_SECRET at minimum

go run cmd/api/main.go
```

Server starts on `:8080` by default (configurable via `PORT` env var).

## Environment Variables

| Variable       | Required      | Description                                                                       |
| -------------- | ------------- | --------------------------------------------------------------------------------- |
| `DATABASE_URL` | In production | Postgres connection string. Falls back to a local dev connection string if unset. |
| `JWT_SECRET`   | Yes           | Secret used to sign JWTs. Must be a long, random string.                          |
| `PORT`         | No            | Port to listen on (defaults to `8080`; set automatically by Render).              |

## Known Limitations / Future Improvements

- **Rate limiting is in-memory** — resets on server restart. A production deployment with multiple instances would need a shared store (Redis).
- **No token revocation** — JWTs are stateless and remain valid until expiry, even after "logout" (which only clears client-side storage). A production system would need a token blacklist or short-lived access tokens with refresh tokens.
- **No password strength validation** — any password is currently accepted.
- **CORS accepts any `*.vercel.app` origin** — convenient for preview deployments during development, but broader than a typical production policy.
- **`class_assigned` on Teacher is a single value** — a teacher can only be assigned one class at a time. Supporting multiple classes would require a many-to-many relationship (a join table) instead of a single column. Deferred as a portfolio-scope simplification.
- **Free-tier Postgres expires periodically** — not a code limitation, but worth noting for anyone evaluating the live deployment.
