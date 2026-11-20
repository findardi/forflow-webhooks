-- name: CreateAttempt :exec
insert into attempts
    (delivery_id, response_status, error_kind)
values
    ($1, $2, $3);
