-- name: CreateEndpoint :one
insert into endpoints
    (application_id, url, secret)
select sqlc.arg(application_id), sqlc.arg(url), sqlc.arg(secret)
where exists (
    select 1 from applications
    where id = sqlc.arg(application_id) and deleted_at is null
)
returning id, application_id, url, created_at, deleted_at;

-- name: ListEndpointsByApplicationID :many
select id, application_id, url, created_at, deleted_at from endpoints
where application_id = $1 and deleted_at is null
order by created_at asc;

-- name: DeleteEndpoint :execrows
update endpoints
set deleted_at = now()
where id = $1 and application_id = $2 and deleted_at is null;

-- name: DeleteEndpointsByApplication :exec
update endpoints
set deleted_at = now()
where application_id = $1 and deleted_at is null;

-- name: RotateSecret :execrows
update endpoints
    set previous_secret = secret,
        previous_secret_expires_at = now() + interval '24 hours',
        secret = sqlc.arg(secret)
where id = sqlc.arg(id) and application_id = sqlc.arg(application_id) and deleted_at is null;
