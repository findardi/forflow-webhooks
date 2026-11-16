-- +goose Up
-- +goose StatementBegin
create table if not exists applications(
    id uuid primary key default uuidv7(),
    name text not null,
    created_at timestamptz not null default now(),
    deleted_at timestamptz default null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists applications;
-- +goose StatementEnd
