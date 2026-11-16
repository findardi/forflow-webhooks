-- +goose Up
-- +goose StatementBegin
create table if not exists deliveries(
    id uuid primary key default uuidv7(),
    message_id uuid not null references messages(id),
    endpoint_id uuid not null references endpoints(id),
    status text not null default 'pending'
        check (status in ('pending', 'succeeded', 'failed', 'canceled')),
    created_at timestamptz not null default now(),
    updated_at timestamptz default null,
    unique (message_id, endpoint_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists deliveries;
-- +goose StatementEnd
