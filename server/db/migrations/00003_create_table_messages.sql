-- +goose Up
-- +goose StatementBegin
create table if not exists messages(
    id uuid primary key default uuidv7(),
    application_id uuid not null references applications (id),
    payload bytea not null,
    created_at timestamptz not null default now()
);

create index if not exists idx_application_id_messages on messages (application_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index if exists idx_application_id_messages;
drop table if exists messages;
-- +goose StatementEnd
