# Go Skeleton

This is a clean Go service skeleton extracted from the original project shape.
Business modules were intentionally removed; the only domain-like code left is
the `Example` flow used to demonstrate the app layers.

## Structure

- `cmd/api`: HTTP API process.
- `cmd/worker`: Asynq worker process.
- `cmd/migrate`: minimal GORM migration entrypoint for the example table.
- `config`: environment loading and dependency registry.
- `internal`: application wiring, routes, middleware, and example layers.
- `pkg`: reusable infrastructure helpers.

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

## Verify

```sh
go test ./...
go vet ./...
```
