# Product

<!-- impeccable:product-schema 1 -->

## Platform

web

## Users

- Developer integrators: backend engineers who use Forflow to send webhooks to their own customers. They register applications and endpoints, and later debug failed deliveries.
- The builder of Forflow: the project owner is the first daily user while the API grows. The dashboard is built for integrators, but it must also work as a practical tool for driving the API during development.

## Product Purpose

Forflow is webhooks as a service. It delivers webhooks reliably for other applications: queues, retries, idempotency, signing, and observability. The dashboard lets a developer manage what Forflow delivers and to where. Success: a developer can set up an application and its endpoints in seconds and can trust what the dashboard shows.

## Positioning

Undecided. Forflow is a learning project that targets the same job as hosted webhook platforms. Do not claim scale, customers, uptime, or features that the API does not have.

## Operating Context

- Not multi-tenant from the user's view. On first login, Forflow creates an organization for the user automatically. There is no org switcher.
- Planned hierarchy (docs/data-model/01-tenancy-erd.png): organization 1:N environments 1:N applications. Users join organizations through user members.
- Core pipeline (docs/data-model/02-core-erd.png): applications own endpoints and messages. A message fans out to deliveries (one per endpoint), and each delivery has attempts (one per HTTP try).

## Capabilities and Constraints

API available today (Go, `server/internal/app/app.go`):

- `POST /applications` with `{ name }`. Name is required, trimmed, max 100 characters. 400 on invalid name.
- `GET /applications`. Lists all applications, oldest first.
- `POST /applications/{appID}/endpoints` with `{ url }`. 400 invalid URL or app ID, 404 application not found, 409 endpoint URL already exists for that application.
- `GET /applications/{appID}/endpoints`. Lists endpoints, oldest first.
- `DELETE /applications/{appID}` and `DELETE /applications/{appID}/endpoints/{endpointID}`. Soft delete; deleting an application deletes its endpoints.
- `POST /applications/{appID}/messages` with `{ payload }` (any JSON value except null, request body under 4 KB). Creates one `pending` delivery per live endpoint in the same transaction. 400 invalid payload or app ID, 404 application not found.
- `GET /applications/{appID}/endpoints/{endpointID}/messages`. Messages delivered to that endpoint with their `delivery_status` (`pending`, `succeeded`, `failed`, `canceled`), oldest first. There is no list per application.
- JSON responses use snake_case. Errors are plain text bodies.
- Records: application `id` (UUIDv7), `name`, `created_at`, `deleted_at`. Endpoint `id`, `application_id`, `url`, `created_at`. Message `id`, `application_id`, `payload`, `created_at`.

Scope rules:

- The dashboard covers only the APIs above. Do not design screens or fake data for delivery status, attempts, event types, or API keys until the API exists.
- No auth for now. Login, environments, and the auto-created organization come later.
- No update endpoints exist yet. No delivery worker exists yet, so every delivery stays `pending`.

Stack: SvelteKit (Svelte 5 runes), TypeScript, Bun, Vite. Prettier: tabs, single quotes, no trailing commas, width 100.

## Evidence on Hand

- ERD images in `docs/data-model/`.
- No logo, brand assets, customers, testimonials, or metrics. Do not invent them.

## Product Principles

1. Show only what is real. Every screen maps to a working API call.
2. Developer-first: IDs, URLs, and timestamps are primary content, easy to read and copy.
3. Fast setup: creating an application or endpoint takes one short form.
4. Grow by addition: the layout must take environments, messages, and deliveries later without a rebuild.
