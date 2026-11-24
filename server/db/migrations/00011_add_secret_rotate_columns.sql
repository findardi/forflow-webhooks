-- +goose Up
-- +goose StatementBegin
alter table endpoints
    add column previous_secret bytea
        constraint endpoints_previous_secret_len_check
        check (octet_length(previous_secret) between 24 and 64),
    add column previous_secret_expires_at timestamptz,
    add constraint endpoints_previous_secret_check 
        check ((previous_secret is null) = (previous_secret_expires_at is null));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table endpoints
    drop column previous_secret_expires_at,
    drop column previous_secret;
-- +goose StatementEnd
