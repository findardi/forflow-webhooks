-- +goose Up
-- +goose StatementBegin
alter table endpoints add column if not exists deleted_at timestamptz default null;

-- Uniqueness only among live endpoints, so a deleted URL can be registered again.
drop index if exists uq_endpoints_application_id_url;
create unique index if not exists uq_endpoints_application_id_url
    on endpoints (application_id, url) where deleted_at is null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Fails if soft-deleted rows duplicate a live (application_id, url); purge them first.
drop index if exists uq_endpoints_application_id_url;
create unique index if not exists uq_endpoints_application_id_url on endpoints (application_id, url);
alter table endpoints drop column if exists deleted_at;
-- +goose StatementEnd
