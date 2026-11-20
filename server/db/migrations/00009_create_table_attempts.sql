-- +goose Up
-- +goose StatementBegin
create table if not exists attempts (
    id uuid primary key default uuidv7(),
    delivery_id uuid not null references deliveries(id),
    response_status int,
    error_kind text,
    created_at timestamptz not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists attempts;
-- +goose StatementEnd
