-- name: CreateApplication :one
insert into applications
    (name)
values
    ($1)
returning *;