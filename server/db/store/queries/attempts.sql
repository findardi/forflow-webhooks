-- name: CreateAttempt :exec
insert into attempts
    (delivery_id, response_status, error_kind)
values
    ($1, $2, $3);

-- name: ListAttemptsByDelivery :many
select id, delivery_id, response_status, error_kind, created_at
from attempts
where delivery_id = $1
order by created_at asc;
