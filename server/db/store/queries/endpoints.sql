-- name: CreateEndpoint :one
insert into endpoints
    (application_id, url)
select sqlc.arg(application_id), sqlc.arg(url)
where exists (
    select 1 from applications
    where id = sqlc.arg(application_id) and deleted_at is null
)
returning *;

-- name: ListEndpointsByApplicationID :many
select * from endpoints
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
