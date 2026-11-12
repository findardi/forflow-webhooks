# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Forflow: webhooks-as-a-service learning project. It delivers webhooks reliably on behalf of other applications (queues, retries, idempotency, signing, observability). Status: design / early scaffolding. There is no Go source, no HTTP server, and no test suite yet.

Monorepo layout:
- `server/`: Go backend (module `github.com/findardi/forflow-webhooks/server`). Postgres through sqlc + pgx/v5. Goose for migrations.
- `web/`: SvelteKit dashboard (Svelte 5, TypeScript, Bun, Vite). Runes mode is forced for all project files in `vite.config.ts`.
- `docs/data-model/`: ERD images. These are the source of truth for the schema until migrations exist.
- `deploy/`, `.github/workflows/`: empty placeholders.

## Setup

Required tools: `go`, `bun`, `docker`, `goose`, `sqlc`, `golangci-lint`.

```sh
cp env.example .env   # Makefile includes .env; DATABASE_URL, GOOSE_DRIVER, GOOSE_MIGRATION_DIR come from it
make db-up            # Postgres 18 on localhost:5433 (compose.dev.yaml)
cd web && bun install
```

## Commands

Run all `make` targets from the repo root.

| Task | Command |
|------|---------|
| Start / stop dev DB | `make db-up` / `make db-down` |
| New migration | `make migrate-create name=<name>` (goose, SQL, sequential numbering via `-s`) |
| Apply / roll back / status | `make migrate-up` / `make migrate-down` / `make migrate-status` |
| Generate Go DB code | `make generate` (runs `sqlc generate` in `server/`) |
| Format | `make format` (gofmt + prettier) |
| Lint | `make lint` (`go mod tidy -diff`, `golangci-lint run`, prettier check, `svelte-check`) |
| Web dev server | `cd web && bun run dev` |
| Web type check | `cd web && bun run check` |

## Database workflow

- Migrations go in `server/db/migrations/`. sqlc reads the same directory as its schema, so every migration changes the generated types.
- Hand-written queries go in `server/db/store/queries/*.sql`. Generated code goes to `server/db/store/sqlc/` (package `sqlc`). Never edit the generated files. Change the SQL and run `make generate`.
- After a schema change: `make migrate-create`, write the up/down SQL, `make migrate-up`, then `make generate`.

## Domain model (from `docs/data-model/`)

Tenancy:
- `organizations` 1:N `environments` 1:N `applications`
- `users` N:M `organizations` through `user_members`

Core delivery pipeline:
- `environments` own `api_keys` (N:1) and `event_types` (1:N).
- `applications` own `endpoints` (subscriber URLs) and `messages` (events sent in by the customer).
- `messages` 1:N `deliveries`; `endpoints` 1:N `deliveries`. One delivery = one message routed to one endpoint.
- `deliveries` 1:N `attempts`. One row per HTTP try; retries create new attempts.

Event subscription (fan-out):
- `endpoints` N:M `event_types` through `endpoint_event_types`.
- `event_types` 1:N `messages`.
- A message fans out to every endpoint of its application that subscribes to the message's event type.

## Communication

- Max 250 words per response, excluding code blocks.

## Web conventions

- Prettier: tabs, single quotes, no trailing commas, print width 100.
- `#lib` / `#lib/*` import aliases (package.json `imports`) map to `src/lib`.
- `adapter-auto` is a placeholder. Choose a real adapter when deployment is decided.
