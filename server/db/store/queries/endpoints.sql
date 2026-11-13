-- name: CreateEndpoint :one
insert into endpoints
    (application_id, url)
values
    ($1, $2)
returning *;

-- name: ListEndpointsByApplicationID :many
select * from endpoints
where application_id = $1
order by created_at asc;
