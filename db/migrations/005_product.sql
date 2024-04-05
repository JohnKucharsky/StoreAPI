-- +goose Up
create table product(
        id serial primary key,
        main_shelf_id int not null references shelf(id) on update cascade on delete restrict,
        name varchar not null,
        serial varchar not null unique,
        price bigint not null,
        model varchar,
        picture_url varchar not null,
        created_at timestamp not null default now(),
        updated_at timestamp not null default now()
);

-- +goose Down
drop table product;