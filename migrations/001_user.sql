-- +goose Up
create table users(
        id uuid primary key,
        name varchar not null,
        last_name varchar not null,
        middle_name varchar,
        email varchar unique not null,
        password varchar not null,
        created_at timestamp not null default now(),
        updated_at timestamp not null default now(),
        unique(name,last_name)
);

-- +goose Down
drop table users;