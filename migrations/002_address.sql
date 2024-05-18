-- +goose Up
create table address(
        id serial primary key,
        city varchar not null,
        street varchar not null,
        house varchar not null,
        floor int,
        entrance int,
        additional_info varchar,
        created_at timestamp not null default now(),
        updated_at timestamp not null default now(),
        unique(city,street,house,floor,entrance)
);

-- +goose Down
drop table address;