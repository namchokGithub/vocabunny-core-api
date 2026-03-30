# vocabunny-core-api

Production-ready Go backend starter for an enterprise monolith pattern:

`protocol -> handler -> service -> repository -> database`

Core design rules:

- `domain` and `port` are the center of the system
- handlers never talk to repositories directly
- business logic stays in services
- query/database logic stays in repositories
- composition root is `protocol/init.go`
- all HTTP route registration is centralized in `protocol/http.go`
- config loads from env under `configs/`
- infrastructure bootstrap stays in `infrastructure/`

## Stack

- Go 1.24+
- Echo
- GORM
- PostgreSQL
- Redis
- JWT
- cron
- websocket manager
- file storage abstraction

## Run

```bash
cp .env-example .env
go mod tidy
go run ./cmd
```

## Project Structure

```text
cmd/
configs/
infrastructure/
internal/
  core/
    domain/
    helper/
    port/
    service/
  handler/
  repository/
pkg/
protocol/
```

## Request Flow

1. `protocol/http.go` wires routes into module handlers.
2. handler binds and validates request DTO, then transforms request into domain input.
3. service executes business rules, transactions, duplicate checks, and orchestration.
4. repository transforms domain input into GORM model and executes database queries.
5. handler transforms domain output into response DTO and returns via shared presenter.

## Adding a New Module

Follow the same pattern as `users`:

1. Add domain entities/value objects to `internal/core/domain`.
2. Add interfaces to `internal/core/port`.
3. Implement business logic in `internal/core/service/<module>_service.go`.
4. Implement GORM model, repository queries, and transformers in `internal/repository/`.
5. Implement DTOs, request/response transformers, and handler in `internal/handler/`.
6. Register dependencies through `internal/repository/repository.go`, `internal/core/service/service.go`, `internal/handler/handler.go`, and `protocol/http.go`.
