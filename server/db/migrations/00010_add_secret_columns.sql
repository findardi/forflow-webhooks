-- +goose Up
-- +goose StatementBegin
alter table endpoints
    add column secret bytea not null
    default uuid_send(gen_random_uuid()) || uuid_send(gen_random_uuid())
    constraint endpoints_secret_len_check check (octet_length(secret) between 24 and 64);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table endpoints
    drop column secret;
-- +goose StatementEnd
