-- name: CreateApplication :one
insert into applications
    (name)
values
    ($1)
returning *;

-- name: ListApplications :many
select * from applications
where deleted_at is null
order by created_at asc;

-- name: DeleteApplication :execrows
update applications
set deleted_at = now()
where id = $1 and deleted_at is null;
