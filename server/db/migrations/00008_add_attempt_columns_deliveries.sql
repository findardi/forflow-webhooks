-- +goose Up
-- +goose StatementBegin
alter table deliveries 
    add column next_attempt_at timestamptz not null default now(),
    add column attempt_count int not null default 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table deliveries 
    drop column attempt_count,
    drop column next_attempt_at;
-- +goose StatementEnd
