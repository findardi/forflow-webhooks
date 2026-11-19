-- name: CreateDeliveriesForMessage :execrows
insert into deliveries
    (message_id, endpoint_id)
select 
    m.id, e.id
from messages m
join endpoints e on e.application_id = m.application_id
where m.id = $1
and e.deleted_at is null;

-- name: CancelStatusDeliveries :exec
update deliveries
set status = 'canceled', updated_at = now()
where endpoint_id = $1 and status = 'pending';

-- name: CancelDeliveriesByApplication :exec
update deliveries
set status = 'canceled', updated_at = now()
where status = 'pending'
and endpoint_id in (select id from endpoints where application_id = $1);

-- name: GetClaim :one
update deliveries d
set status = 'processing', updated_at = now()
from messages m, endpoints e
where d.id = (
    select id from deliveries
    where status = 'pending'
    order by created_at
    limit 1
    for update skip locked
)
and m.id = d.message_id
and e.id = d.endpoint_id
returning d.id as delivery_id, m.id as message_id, e.url, m.payload;

-- name: MarkSuccess :execrows
update deliveries
set status = 'succeeded', updated_at = now()
where id = $1 and status = 'processing';

-- name: MarkFailed :execrows
update deliveries
set status = 'failed', updated_at = now()
where id = $1 and status = 'processing';