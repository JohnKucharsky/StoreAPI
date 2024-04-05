-- +goose Up
create table shelf_product(
        shelf_id int not null references shelf(id) on delete cascade,
        product_id int not null references product(id) on delete cascade,
        unique (shelf_id,product_id)
);

-- +goose Down
drop table shelf_product;