-- name: CreateMessages :one
insert into messages
    (application_id, payload)
select sqlc.arg(application_id), sqlc.arg(payload)
where exists (
    select 1 from applications
    where id = sqlc.arg(application_id) and deleted_at is null
)
returning *;

-- name: ListMessagesByEndpoint :many
select m.id, m.application_id, m.payload, m.created_at,
    d.id as delivery_id, d.status as delivery_status
from messages m
join deliveries d on d.message_id = m.id
where d.endpoint_id = sqlc.arg(endpoint_id)
and m.application_id = sqlc.arg(application_id)
order by m.created_at asc;

-- name: ListMessagesByApplication :many
select m.id, m.application_id, m.payload, m.created_at,
    d.id as delivery_id, d.status as delivery_status, e.url as endpoint_url
from messages m
left join deliveries d on d.message_id = m.id
left join endpoints e on e.id = d.endpoint_id
where m.application_id = $1
order by m.created_at asc, m.id asc, e.url desc;
