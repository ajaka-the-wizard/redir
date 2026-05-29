# redir

redir is a Go backend for a media upload and redirect service. It combines user authentication, product-scoped API keys, Redis-backed sessions and caching, Postgres persistence, S3-compatible object storage, presigned asset redirects, request logging, rate limiting, and access metrics.

This repository is an active backend project. The codebase is already shaped like a real service, but some production-hardening work is still ongoing around documentation, response consistency, cache invalidation rules, and test ergonomics.

## Current status

- HTTP API: Gin server with versioned routes under `/api/v1`.
- Auth: registration, login, logout, email verification, session cookies, Google OAuth, and GitHub OAuth scaffolding.
- Sessions and tokens: Redis-backed session storage and verification tokens with TTLs.
- Products: authenticated users can create products, toggle visibility, and generate private API keys.
- Client API: external clients authenticate with product IDs and bearer keys.
- Uploads: clients can upload multipart media batches to S3-compatible object storage.
- Assets: public asset keys resolve to presigned storage URLs and redirect callers.
- Metrics: asset redirects capture browser, OS, device, IP, and media metadata.
- Storage layers: Postgres repositories for durable data, Redis cache/store helpers for transient and cached data.
- Middleware: request IDs, request-scoped structured logging, latency/status logging, recovery, validation, ownership checks, and rate limiting.
- Tests: includes DB/Redis-backed integration coverage for the auth lifecycle.

## What works today

- User registration with bcrypt password hashing.
- Email verification flow backed by Redis tokens.
- Login and logout with Redis-backed session cookies.
- Google OAuth login flow.
- Product creation and private key generation.
- Product ownership checks for protected product operations.
- API-key validation for client upload routes.
- Multipart upload flow using AWS SDK v2 transfer manager.
- Public asset redirect through cached presigned URLs.
- Redis cache-aside reads for products and media.
- Postgres persistence through `pgx` repositories.
- Request-scoped logs with request IDs and latency/status summaries.
- Route-specific rate limiting for auth, user/product, and client routes.

## Known limitations

- Local setup docs still need more executable examples, especially for migrations and external services.
- Some exported names need Go-style cleanup and shortening.
- Cache invalidation rules should be documented and tightened for products, users, media, and presigned URLs.
- Some handler unit tests are currently commented out while integration tests cover the heavier auth flow.
- Local test setup depends on external services such as Postgres and Redis.
- `go.mod` dependency metadata should be cleaned up with `go mod tidy`.

## Architecture

The server starts from `cmd/api/main.go`, which calls `internal.Listen()`. Startup loads environment variables, initializes Redis, Postgres, S3-compatible storage, a user-agent parser, repositories, and the store layer before mounting Gin routes.

The main package boundaries are:

- `internal/routes`: route registration and route grouping.
- `internal/handlers`: request handlers for auth, users, products, clients, and assets.
- `internal/middlewares`: validation, auth, ownership checks, rate limiting, request IDs, logging, and recovery.
- `internal/store`: thin orchestration layer between handlers/middleware, Redis, and repositories.
- `internal/cache`: Redis-backed sessions, verification tokens, cached media/products, and presigned URLs.
- `internal/repository`: Postgres queries and persistence logic using `pgx`.
- `internal/configs`: environment loading, storage setup, and user-agent parser setup.
- `internal/models` and `internal/domain`: persistent models and request/domain types.
- `internal/utils`: shared helpers for IDs, keys, cookies, hashing, and context extraction.

## API overview

All primary routes are mounted under `/api/v1`.

### Auth

- `POST /auth/register`
- `POST /auth/login`
- `GET /auth/verify`
- `POST /auth/logout`
- `GET /auth/oauth/google`
- `GET /auth/oauth/google/callback`
- `GET /auth/oauth/github`

### Users

- `GET /users/me`

### Products

- `POST /product`
- `POST /product/:id`
- `PUT /product/:id`
- `PUT /product/:id/assets/:assetId`

### Client API

- `GET /client/ping`
- `POST /client/upload`
- `PUT /client/commit/:batchId`

Client routes expect product credentials through headers:

- `X-Product`: product ID
- `Authorization`: `Bearer <private-key>`

### Assets

- `GET /assets/:assetId`

Public asset requests validate the public key, check asset visibility, generate or reuse a cached presigned URL, save access metrics, and redirect to object storage.

## Configuration

The app currently loads configuration from `.env` using `godotenv`.

Required or supported environment variables include:

- `DATABASE_URL`
- `PORT`
- `SERVER_ADDRESS`
- `ENVIRONMENT`
- `GOOGLE_CLIENT_ID`
- `GOOGLE_CLIENT_SECRET`
- `GOOGLE_REDIRECT_URL`
- `GITHUB_CLIENT_ID`
- `GITHUB_CLIENT_SECRET`
- `GITHUB_REDIRECT_URL`
- `CLIENT_DASHBOARD`
- `CLIENT_LOGIN_URL`
- `STORAGE_SERVICE_ACCESS_KEY_ID`
- `STORAGE_SERVICE_SECRET_ACCESS_KEY`
- `STORAGE_SERVICE_ENDPOINT`
- `BUCKET_NAME`
- `BUCKET_ROOT`
- `DATA_GET_PATH`
- `REDIS_ADDR`
- `REDIS_PASSWORD`

`ENVIRONMENT=production` enables production mode behavior where configured.

## Local development

Build or run the server:

```sh
go run cmd/api/main.go
```

Or use the make target:

```sh
make server
```

Clean module metadata:

```sh
go mod tidy
```

Run tests:

```sh
go test ./...
```

Some tests are integration tests and expect services such as Postgres and Redis to be configured through `.env`. They may take longer than small unit tests because they exercise real request flows and backing services.

## Database

SQL migrations live in `migrations/` and cover users, media, metrics, product tables, encryption helpers, and later media/metrics schema changes.

Migration execution is not yet wrapped in a documented make target, so run migrations using your chosen migration tool against `DATABASE_URL`.

## Why this project matters

redir demonstrates practical backend engineering beyond simple CRUD. It includes authentication, sessions, OAuth, caching, API-key validation, object storage, multipart uploads, public redirect flows, structured logging, rate limiting, database access, and integration testing across real infrastructure.

It is a strong portfolio project because it shows the shape of a real service: user-facing auth flows, external client APIs, storage integration, cache/database coordination, operational middleware, and a path toward production hardening.

## Roadmap

- Standardize API response shapes.
- Clean up exported names and Go idioms.
- Document cache invalidation rules.
- Add Docker Compose for local Postgres, Redis, and S3-compatible storage.
- Add migration commands to the makefile.
- Add OpenAPI documentation or a concise route reference.
- Run `go mod tidy` and keep dependency metadata clean.
