-- +goose Up
-- +goose StatementBegin
alter table deliveries drop constraint if exists deliveries_status_check;
alter table deliveries add constraint deliveries_status_check 
        check (status in ('pending', 'processing', 'succeeded', 'failed', 'canceled'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
update deliveries set status = 'pending' where status = 'processing';
alter table deliveries drop constraint if exists deliveries_status_check;
alter table deliveries add constraint deliveries_status_check
    check (status in ('pending', 'succeeded', 'failed', 'canceled'));
-- +goose StatementEnd
