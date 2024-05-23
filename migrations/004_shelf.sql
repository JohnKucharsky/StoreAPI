-- +goose Up
create table shelf(
        id serial primary key,
        name varchar not null unique,
        created_at timestamp not null default now(),
        updated_at timestamp not null default now()
);

-- +goose Down
drop table shelf;