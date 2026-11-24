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
    select p.id from deliveries p
    join endpoints pe on pe.id = p.endpoint_id and pe.deleted_at is null
    where p.status = 'pending'
        and p.next_attempt_at <= now()
        and (select count(*) from deliveries x
            where x.endpoint_id = p.endpoint_id
            and x.status = 'processing') < sqlc.arg(max_per_endpoint)::int
    order by p.next_attempt_at
    limit 1
    for update skip locked
)
and m.id = d.message_id
and e.id = d.endpoint_id
returning d.id as delivery_id, m.id as message_id, e.url, m.payload, d.attempt_count, e.secret,
    (case when e.previous_secret_expires_at > now() then e.previous_secret end)::bytea as previous_secret;

-- name: MarkSuccess :execrows
update deliveries
set status = 'succeeded', attempt_count = attempt_count + 1, updated_at = now()
where id = $1 and status = 'processing';

-- name: MarkFailed :execrows
update deliveries
set status = 'failed', attempt_count = attempt_count + 1, updated_at = now()
where id = $1 and status = 'processing';

-- name: RetryAttempt :execrows
update deliveries
set status = 'pending',
    attempt_count = attempt_count + 1,
    next_attempt_at = $2,
    updated_at = now()
where id = $1 and status = 'processing';

-- name: GetDeliveryByApplication :one
select d.id, d.message_id, d.endpoint_id, d.status, d.attempt_count, d.next_attempt_at,
    d.created_at, d.updated_at, e.url as endpoint_url, m.payload
from deliveries d
join messages m on m.id = d.message_id
join endpoints e on e.id = d.endpoint_id
where d.id = sqlc.arg(id)
and m.application_id = sqlc.arg(application_id);

-- name: ResetProcessing :execrows
update deliveries set status = 'pending', updated_at = now()
where status = 'processing';