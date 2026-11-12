-- name: CreateEndpoint :one
insert into endpoints
    (application_id, url)
values
    ($1, $2)
returning *;