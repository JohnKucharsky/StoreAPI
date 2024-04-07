-- +goose Up
create table order_product(
        product_id int not null references product(id) on delete cascade,
        order_id int not null references orders(id) on delete cascade
);

-- +goose Down
drop table order_product;