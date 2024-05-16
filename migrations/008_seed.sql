-- +goose Up
insert into address (city, street, house, floor, entrance)
values ('Krasnodar', 'Red', '1b', 4, 2);

insert into shelf (name)
values ('FirstShelf'), ('SecondShelf'), ('ThirdShelf');

insert into product (name, serial, price, model, picture_url) values
 ('Phone', '34ko34j', 250, 'iPhone', 'https://picture.234kjh.com'),
 ('Notebook', 'sudo34j', 850, 'Honor', 'https://picture.4kjh.com'),
 ('Tab', '32444j', 450, 'Samsung', 'https://picture.2sjh.com'),
 ('Vacuum', '42ssg', 650, null, 'https://picture.23rh.com'),
 ('Watch', '34csso34j', 1600, 'Rolex', 'https://picture.2dkjh.com');

insert into shelf_product (shelf_id, product_id, product_qty) values
    (2,1),(3,2,2),(1,3),(2,4,4),(1,5);

insert into orders (address_id, payment) VALUES
    (1,'cash');

insert into order_product (product_id, order_id) VALUES
 (1,1),(2,1),(3,1),(4,1),(5,1);

-- +goose Down