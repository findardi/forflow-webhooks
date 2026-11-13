-- name: CreateEndpoint :one
insert into endpoints
    (application_id, url)
values
    ($1, $2)
returning *;

-- name: GetByAppIdandURL :one
select * from endpoints
where 
    application_id = $1
and
    url = $2;