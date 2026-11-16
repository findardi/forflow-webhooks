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
