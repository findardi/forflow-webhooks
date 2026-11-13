-- name: CreateApplication :one
insert into applications
    (name)
values
    ($1)
returning *;

-- name: ListApplications :many
select * from applications
order by created_at asc;
