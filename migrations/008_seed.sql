-- +goose Up
insert into address (city, street, house, floor, entrance)
values ('Krasnodar', 'Red', '1b', 4, 2);

insert into users (id, name, last_name, middle_name, email, password)
values ('018f8d76-7f77-701d-8b43-42a7be65212a', 'Name','Last','Middle','test@mail.com',
        '$argon2id$v=19$m=65536,t=3,p=4$OFoNuVrpjugGLiezadJy1g$KVbw11kKeb5haI72uekAOZFsBJQ3OqGBkESkwRvoAmI');

insert into shelf (name)
values ('FirstShelf'), ('SecondShelf'), ('ThirdShelf');

insert into product (name, serial, price, model, picture_url) values ('Carbonated Water - Strawberry', '233358583-8', 324, 'Motorola A6188', 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Wine - Red, Cabernet Merlot', '968454794-3', 126, 'Samsung F200', 'http://dummyimage.com/50x50.png/dddddd/000000');
insert into product (name, serial, price, model, picture_url) values ('Sauce - White, Mix', '009819058-X', 280, null, 'http://dummyimage.com/50x50.png/cc0000/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Juice - Tomato, 48 Oz', '330007551-3', 188, null, 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Ham Black Forest', '261954356-8', 740, 'Motorola W205', 'http://dummyimage.com/50x50.png/ff4444/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Wine - Alicanca Vinho Verde', '337831355-2', 388, 'vivo Y50', 'http://dummyimage.com/50x50.png/dddddd/000000');
insert into product (name, serial, price, model, picture_url) values ('Chips Potato Salt Vinegar 43g', '948730024-4', 483, null, 'http://dummyimage.com/50x50.png/dddddd/000000');
insert into product (name, serial, price, model, picture_url) values ('Cranberries - Fresh', '159186005-9', 601, 'ZTE Grand X Plus Z826', 'http://dummyimage.com/50x50.png/ff4444/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Wine - Cave Springs Dry Riesling', '052347553-5', 509, 'Siemens Xelibri 1', 'http://dummyimage.com/50x50.png/ff4444/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Grapes - Black', '414305465-5', 501, 'Yezz Andy 5E2I', 'http://dummyimage.com/50x50.png/cc0000/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Ham - Smoked, Bone - In', '068831670-0', 513, 'ZTE Blade S6', 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Beans - Fava, Canned', '943999533-5', 612, 'LG KP210', 'http://dummyimage.com/50x50.png/dddddd/000000');
insert into product (name, serial, price, model, picture_url) values ('Plasticforkblack', '457209086-6', 115, null, 'http://dummyimage.com/50x50.png/dddddd/000000');
insert into product (name, serial, price, model, picture_url) values ('Beef - Short Loin', '533899652-9', 250, null, 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Glove - Cutting', '780437625-1', 513, 'OnePlus 8', 'http://dummyimage.com/50x50.png/ff4444/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Nut - Chestnuts, Whole', '851921076-7', 789, 'Philips S200', 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Mop Head - Cotton, 24 Oz', '398075909-1', 185, null, 'http://dummyimage.com/50x50.png/cc0000/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Pastry - Banana Tea Loaf', '389960292-7', 554, 'Samsung S7550 Blue Earth', 'http://dummyimage.com/50x50.png/cc0000/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Puff Pastry - Sheets', '371995723-3', 92, null, 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Pepper - Cubanelle', '791316246-5', 266, null, 'http://dummyimage.com/50x50.png/dddddd/000000');
insert into product (name, serial, price, model, picture_url) values ('Pepper - Cayenne', '912107818-1', 732, 'LG GW620', 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Bread Bowl Plain', '816259584-8', 178, null, 'http://dummyimage.com/50x50.png/cc0000/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Bay Leaf Fresh', '235924875-8', 178, 'Sony Xperia Z3v', 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Herb Du Provence - Primerba', '636835008-3', 237, null, 'http://dummyimage.com/50x50.png/ff4444/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Calypso - Pineapple Passion', '549955265-7', 783, 'Samsung Galaxy Note 8.0', 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Flour - Fast / Rapid', '468276578-7', 698, 'Samsung U320 Haven', 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Miso Paste White', '372498038-8', 159, 'Lenovo S820', 'http://dummyimage.com/50x50.png/cc0000/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Trout - Rainbow, Frozen', '767374807-X', 391, 'LG Optimus L5 II E460', 'http://dummyimage.com/50x50.png/cc0000/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Butcher Twine 4r', '259774742-5', 405, 'ZTE Grand Memo II LTE', 'http://dummyimage.com/50x50.png/cc0000/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Leeks - Large', '904417854-7', 232, 'Amoi A210', 'http://dummyimage.com/50x50.png/ff4444/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Soup - Campbells, Minestrone', '915285378-0', 91, 'Samsung Ch@t 322', 'http://dummyimage.com/50x50.png/dddddd/000000');
insert into product (name, serial, price, model, picture_url) values ('Wine - Red, Pinot Noir, Chateau', '570561662-7', 323, null, 'http://dummyimage.com/50x50.png/ff4444/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Tea - Honey Green Tea', '897390111-7', 764, 'Asus Memo Pad ME172V', 'http://dummyimage.com/50x50.png/cc0000/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Liners - Banana, Paper', '773133079-6', 260, 'alcatel Idol 2 Mini', 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Dried Figs', '531454496-2', 57, 'Motorola DROID X2', 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Beans - Yellow', '715752327-2', 644, 'Asus V55', 'http://dummyimage.com/50x50.png/cc0000/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Dried Figs', '329940951-4', 678, null, 'http://dummyimage.com/50x50.png/dddddd/000000');
insert into product (name, serial, price, model, picture_url) values ('Chocolate Bar - Oh Henry', '256945030-8', 753, null, 'http://dummyimage.com/50x50.png/dddddd/000000');
insert into product (name, serial, price, model, picture_url) values ('Napkin Colour', '054327282-6', 137, 'Motorola V171', 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Lamb - Pieces, Diced', '465293897-7', 517, 'Motorola DROID RAZR XT912', 'http://dummyimage.com/50x50.png/ff4444/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Corn Shoots', '136516734-8', 551, 'Lava 3G 402+', 'http://dummyimage.com/50x50.png/cc0000/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Island Oasis - Sweet And Sour Mix', '847716715-X', 414, 'BLU C6', 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Pork Salted Bellies', '206461036-7', 97, 'vivo iQOO 7', 'http://dummyimage.com/50x50.png/ff4444/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Squid - U 5', '859222438-1', 485, 'ZTE Grand X Pro', 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Pur Source', '915951319-5', 642, 'Philips S308', 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Cheese - Le Cheve Noir', '139913322-5', 402, 'Sagem MY C3-2', 'http://dummyimage.com/50x50.png/cc0000/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Sauce Bbq Smokey', '184928209-9', 361, null, 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Cookies - Englishbay Wht', '988880624-6', 119, null, 'http://dummyimage.com/50x50.png/dddddd/000000');
insert into product (name, serial, price, model, picture_url) values ('Milkettes - 2%', '013381609-5', 56, null, 'http://dummyimage.com/50x50.png/dddddd/000000');
insert into product (name, serial, price, model, picture_url) values ('Chicken - Wieners', '857561393-6', 776, 'LG KM555E', 'http://dummyimage.com/50x50.png/ff4444/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Longos - Greek Salad', '664362325-1', 371, 'Nokia 3310 3G', 'http://dummyimage.com/50x50.png/ff4444/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Oil - Coconut', '264566449-0', 350, null, 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Mix Pina Colada', '107692813-7', 455, 'alcatel OT-819 Soul', 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Milk - Condensed', '882159273-1', 720, 'Karbonn A21', 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Red Currants', '131476093-9', 662, 'i-mobile TV 535', 'http://dummyimage.com/50x50.png/cc0000/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Chocolate Bar - Oh Henry', '198497093-3', 201, 'Samsung Z140', 'http://dummyimage.com/50x50.png/dddddd/000000');
insert into product (name, serial, price, model, picture_url) values ('Gatorade - Orange', '475231576-9', 342, null, 'http://dummyimage.com/50x50.png/ff4444/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Wine - Pinot Noir Pond Haddock', '889535990-9', 169, 'alcatel OT-385', 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Chips - Assorted', '624232285-4', 479, 'Sony Xperia 1', 'http://dummyimage.com/50x50.png/ff4444/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Bacardi Breezer - Tropical', '251420181-0', 259, 'Unnecto Rush', 'http://dummyimage.com/50x50.png/dddddd/000000');
insert into product (name, serial, price, model, picture_url) values ('Oil - Olive, Extra Virgin', '716025713-8', 353, null, 'http://dummyimage.com/50x50.png/ff4444/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Chocolate Eclairs', '293707875-3', 747, 'Vodafone 710', 'http://dummyimage.com/50x50.png/ff4444/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Muffin Orange Individual', '928166567-0', 50, 'BLU Touchbook M7 Pro', 'http://dummyimage.com/50x50.png/dddddd/000000');
insert into product (name, serial, price, model, picture_url) values ('Chinese Foods - Thick Noodles', '108224642-5', 782, 'Celkon C350', 'http://dummyimage.com/50x50.png/ff4444/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Chinese Foods - Chicken Wing', '289231081-4', 738, 'Microsoft Lumia 1030', 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Wine - Magnotta - Bel Paese White', '058067002-3', 463, 'Eten M500', 'http://dummyimage.com/50x50.png/ff4444/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Brandy - Bar', '187136659-3', 106, null, 'http://dummyimage.com/50x50.png/cc0000/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Syrup - Golden, Lyles', '128325879-X', 325, 'Xiaomi Redmi 9AT', 'http://dummyimage.com/50x50.png/dddddd/000000');
insert into product (name, serial, price, model, picture_url) values ('Pasta - Rotini, Dry', '542867215-3', 315, 'BLU Energy X 2', 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Pepsi - Diet, 355 Ml', '912917652-2', 485, 'Mitac MIO A702', 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Irish Cream - Butterscotch', '257760234-0', 542, 'LG Optimus 2 AS680', 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Rum - Spiced, Captain Morgan', '476602930-5', 498, 'T-Mobile Shadow', 'http://dummyimage.com/50x50.png/dddddd/000000');
insert into product (name, serial, price, model, picture_url) values ('Soup - Campbells Tomato Ravioli', '987799575-1', 268, 'Nokia 5.3', 'http://dummyimage.com/50x50.png/ff4444/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Onion - Dried', '436258340-8', 106, 'Celkon Millennia Epic Q550', 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Lid Tray - 12in Dome', '118045993-8', 277, 'Plum Z708', 'http://dummyimage.com/50x50.png/ff4444/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Wine - Segura Viudas Aria Brut', '328664394-7', 211, 'BenQ P31', 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Pear - Packum', '595659014-9', 592, 'Gionee F205', 'http://dummyimage.com/50x50.png/cc0000/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Chip - Potato Dill Pickle', '928905890-0', 603, 'Spice KT-5353', 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Wine - Gato Negro Cabernet', '692155523-2', 495, null, 'http://dummyimage.com/50x50.png/cc0000/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Lamb - Leg, Boneless', '322196501-X', 144, null, 'http://dummyimage.com/50x50.png/dddddd/000000');
insert into product (name, serial, price, model, picture_url) values ('Pizza Pizza Dough', '189762066-7', 255, 'Nokia 5140', 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Juice - Happy Planet', '923342693-9', 269, null, 'http://dummyimage.com/50x50.png/cc0000/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Salmon - Fillets', '411241286-X', 509, null, 'http://dummyimage.com/50x50.png/cc0000/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Green Scrubbie Pad H.duty', '135498000-X', 177, 'ZTE F103', 'http://dummyimage.com/50x50.png/ff4444/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Lychee', '077914980-7', 85, 'Nokia 2.1', 'http://dummyimage.com/50x50.png/cc0000/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Guava', '686647122-2', 529, 'Motorola Moto Z Play', 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Cookie Dough - Oatmeal Rasin', '407409762-1', 524, 'Philips 598', 'http://dummyimage.com/50x50.png/cc0000/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Tea - Darjeeling, Azzura', '081988170-8', 376, null, 'http://dummyimage.com/50x50.png/ff4444/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Rosemary - Primerba, Paste', '490726065-2', 232, null, 'http://dummyimage.com/50x50.png/dddddd/000000');
insert into product (name, serial, price, model, picture_url) values ('Quinoa', '375567440-8', 674, 'Samsung P940', 'http://dummyimage.com/50x50.png/cc0000/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Bacardi Breezer - Strawberry', '286856550-6', 784, null, 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Bread - Bagels, Plain', '088743655-2', 496, 'Samsung Galaxy Z Fold2 5G', 'http://dummyimage.com/50x50.png/ff4444/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Limes', '337728757-4', 792, null, 'http://dummyimage.com/50x50.png/dddddd/000000');
insert into product (name, serial, price, model, picture_url) values ('Veal - Ground', '858813329-6', 709, 'Sharp 770SH', 'http://dummyimage.com/50x50.png/cc0000/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Shrimp - Black Tiger 8 - 12', '884046769-6', 174, null, 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Fireball Whisky', '030817211-6', 62, 'Celkon C5055', 'http://dummyimage.com/50x50.png/dddddd/000000');
insert into product (name, serial, price, model, picture_url) values ('Wine - Prem Select Charddonany', '728087413-4', 597, 'Amoi M630', 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Wine - Mas Chicet Rose, Vintage', '330494678-0', 399, 'Bird S789', 'http://dummyimage.com/50x50.png/cc0000/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Chinese Lemon Pork', '351913685-6', 131, 'Sony Xperia Z5 Premium Dual', 'http://dummyimage.com/50x50.png/dddddd/000000');
insert into product (name, serial, price, model, picture_url) values ('Nacho Chips', '680419772-0', 81, 'Samsung Galaxy Xcover 4s', 'http://dummyimage.com/50x50.png/dddddd/000000');
insert into product (name, serial, price, model, picture_url) values ('Veal - Slab Bacon', '093967094-1', 353, null, 'http://dummyimage.com/50x50.png/ff4444/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Sauerkraut', '137040284-8', 410, 'Samsung Google Nexus 10 P8110', 'http://dummyimage.com/50x50.png/ff4444/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Island Oasis - Mango Daiquiri', '583191021-0', 47, 'Celkon C619', 'http://dummyimage.com/50x50.png/cc0000/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Cup Translucent 9 Oz', '514312995-8', 182, 'Palm Treo 270', 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Veal - Loin', '236671682-6', 496, 'Sagem myMobileTV', 'http://dummyimage.com/50x50.png/ff4444/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Soup - Campbells, Lentil', '791007001-2', 681, 'Telit G83', 'http://dummyimage.com/50x50.png/dddddd/000000');
insert into product (name, serial, price, model, picture_url) values ('Rum - Cream, Amarula', '699183776-8', 378, 'Panasonic Eluga U', 'http://dummyimage.com/50x50.png/ff4444/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Melon - Watermelon Yellow', '908643373-1', 67, 'LG KS10', 'http://dummyimage.com/50x50.png/ff4444/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Cocoa Feuilletine', '381981498-1', 647, 'BLU Energy X LTE', 'http://dummyimage.com/50x50.png/dddddd/000000');
insert into product (name, serial, price, model, picture_url) values ('Soup - Tomato Mush. Florentine', '690254145-0', 599, 'Nokia 6650', 'http://dummyimage.com/50x50.png/dddddd/000000');
insert into product (name, serial, price, model, picture_url) values ('Peach - Halves', '753887289-2', 324, 'Nokia N800', 'http://dummyimage.com/50x50.png/ff4444/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Lamb - Loin, Trimmed, Boneless', '754402495-4', 788, 'Nokia 3555', 'http://dummyimage.com/50x50.png/ff4444/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Carbonated Water - Orange', '416976737-3', 25, 'BenQ T51', 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Rappini - Andy Boy', '499323734-0', 453, 'LG Spirit', 'http://dummyimage.com/50x50.png/dddddd/000000');
insert into product (name, serial, price, model, picture_url) values ('Beer - Muskoka Cream Ale', '394558818-9', 282, null, 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Cherries - Maraschino,jar', '240075912-X', 51, null, 'http://dummyimage.com/50x50.png/5fa2dd/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Munchies Honey Sweet Trail Mix', '110157478-X', 414, 'Bird A130', 'http://dummyimage.com/50x50.png/dddddd/000000');
insert into product (name, serial, price, model, picture_url) values ('Water - Evian 355 Ml', '361115972-7', 226, 'Gigabyte GSmart G1317 Rola', 'http://dummyimage.com/50x50.png/dddddd/000000');
insert into product (name, serial, price, model, picture_url) values ('Curry Paste - Madras', '191473175-1', 31, 'ZTE Blade V7', 'http://dummyimage.com/50x50.png/cc0000/ffffff');
insert into product (name, serial, price, model, picture_url) values ('Fiddlehead - Frozen', '772938757-3', 385, null, 'http://dummyimage.com/50x50.png/ff4444/ffffff');


insert into shelf_product (shelf_id, product_id, product_qty) values
    (3,2,2),(2,4,4);

insert into shelf_product (shelf_id, product_id) values
    (2,1),(1,3),(1,5);

insert into orders (address_id,payment,user_id) VALUES
    (1,'cash','018f8d76-7f77-701d-8b43-42a7be65212a');

insert into order_product (product_id, order_id) VALUES
 (1,1),(2,1),(3,1),(4,1),(5,1);

-- +goose Down