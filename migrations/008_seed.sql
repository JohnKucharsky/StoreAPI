-- +goose Up
insert into address (city, street, house, floor, entrance)
values ('Krasnodar', 'Red', '1b', 4, 2);

insert into users (id, name, last_name, middle_name, email, password)
values ('018f8d76-7f77-701d-8b43-42a7be65212a', 'Name','Last','Middle','test@mail.com',
        '$argon2id$v=19$m=65536,t=3,p=4$OFoNuVrpjugGLiezadJy1g$KVbw11kKeb5haI72uekAOZFsBJQ3OqGBkESkwRvoAmI');

insert into shelf (name)
values ('FirstShelf'), ('SecondShelf'), ('ThirdShelf');

insert into product (name, serial, price, model, picture_url) values
 ('Phone', '34ko34j', 250, 'iPhone', 'https://unsplash.com/photos/gold-iphone-6-sweUF7FcyP4'),
 ('Notebook', 'sudo34j', 850, 'Honor', 'https://unsplash.com/photos/black-and-silver-retractable-pen-on-blank-book-3ym6i13Y9LU'),
 ('Tab', '32444j', 450, 'Samsung', 'https://picture.2sjh.com'),
 ('Vacuum', '42ssg', 650, null, 'https://picture.23rh.com'),
 ('Watch', '34csso34j', 1600, 'Rolex', 'https://picture.2dkjh.com');

insert into shelf_product (shelf_id, product_id, product_qty) values
    (3,2,2),(2,4,4);

insert into shelf_product (shelf_id, product_id) values
    (2,1),(1,3),(1,5);

insert into orders (address_id,payment,user_id) VALUES
    (1,'cash','018f8d76-7f77-701d-8b43-42a7be65212a');

insert into order_product (product_id, order_id) VALUES
 (1,1),(2,1),(3,1),(4,1),(5,1);

-- +goose Down