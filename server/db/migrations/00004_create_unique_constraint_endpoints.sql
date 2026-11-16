-- +goose Up
-- +goose StatementBegin
create unique index if not exists uq_endpoints_application_id_url on endpoints (application_id, url);
drop index if exists idx_application_id_endpoints;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
create index if not exists idx_application_id_endpoints on endpoints (application_id);
drop index if exists uq_endpoints_application_id_url;
-- +goose StatementEnd
