# Go Skeleton

This is a clean Go service skeleton extracted from the original project shape.
Business modules were intentionally removed; the only domain-like code left is
the `Example` flow used to demonstrate the app layers.

## Structure

- `cmd/api`: HTTP API process.
- `cmd/worker`: Asynq worker process.
- `cmd/migrate`: minimal GORM migration entrypoint for the example table.
- `config`: environment loading and typed configuration values.
- `internal/bootstrap`: process-level resource initialization and lifecycle.
- `internal`: application wiring, routes, middleware, and example layers.
- `pkg`: reusable infrastructure helpers, including generic JWT auth.

## Run

```sh
cp .env.example .env
go run ./cmd/api
```

Run the worker when Redis is configured:

```sh
go run ./cmd/worker
```

Run the example migration when Postgres is configured:

```sh
go run ./cmd/migrate
```

## Runtime Dependencies

- The API process requires `POSTGRES`.
- Redis is optional for the API process. When configured, it enables cache and queue publishing.
- The worker process requires `REDIS_ADDR`.
- Postgres is optional for the worker process.

## Deployment Notes

- Swagger is not enabled in this skeleton. Deployment does not require `swag init`.
- If Swagger is added later, generate docs during development or CI build, not at service startup.
- `CORS_ALLOW_ORIGINS` is a comma-separated allow list. Empty means no CORS allow headers.
- API business errors use the JSON envelope `code`, `msg`, and `reason`; most API errors are returned with HTTP 200 by convention.
- `/health` uses real HTTP status codes and returns 503 when required dependencies are unavailable.

## Verify

```sh
go test ./...
go vet ./...
golangci-lint run
```
