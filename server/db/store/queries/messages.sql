-- name: CreateMessages :one
insert into messages
    (application_id, payload)
values
    ($1, $2)
returning *;