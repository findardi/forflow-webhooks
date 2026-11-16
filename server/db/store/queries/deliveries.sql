-- name: CreateDeliveriesForMessage :execrows
insert into deliveries
    (message_id, endpoint_id)
select 
    m.id, e.id
from messages m
join endpoints e on e.application_id = m.application_id
where m.id = $1
and e.deleted_at is null;