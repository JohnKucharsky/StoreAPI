-- +goose Up
create table orders(
        id serial primary key,
        user_id uuid not null references users(id) on update cascade on delete restrict,
        address_id int not null references address(id) on update cascade on delete restrict,
        total varchar not null,
        payment varchar not null,
        created_at timestamp not null default now(),
        updated_at timestamp not null default now()
);

-- +goose Down
drop table orders;