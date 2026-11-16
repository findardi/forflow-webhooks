-- +goose Up
-- +goose StatementBegin
create table if not exists endpoints(
    id uuid primary key default uuidv7(),
    application_id uuid not null references applications (id),
    url text not null,
    created_at timestamptz not null default now()
);

create index if not exists idx_application_id_endpoints on endpoints (application_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index if exists idx_application_id_endpoints;
drop table if exists endpoints;
-- +goose StatementEnd
