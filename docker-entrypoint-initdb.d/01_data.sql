INSERT
INTO roles (name)
VALUES ('USER'),
       ('ADMIN');

INSERT
INTO users (login, password)
VALUES ('user1','$2a$10$5h6GDvR0EBtCECFgptg6iuiOu0jkc/qJ8if9jt39NY9ir602nOcXu'),
       ('user2','$2a$10$5h6GDvR0EBtCECFgptg6iuiOu0jkc/qJ8if9jt39NY9ir602nOcXu');

INSERT
INTO userroles (user_id, role_id)
VALUES (1, 1),
       (2, 1),
       (2, 2);


INSERT
INTO shops (name, address, lon, lat, working_hours)
VALUES ('Магазин на диване', 'Москва, Останкино', '324234' , '5465476', '8 - 20'),
       ('Магазин для взрослых', 'Ростов, кремль', '12334' , '5465476', '8 - 20'),
       ('Так себе магазин', 'Одесса, привоз', '8394' , '632542', '10 - 22'),
       ('Магазин где есть все', 'Краснодар, центр', '86932' , '324234', '0 - 24');

INSERT
INTO categories (name, uri_name)
VALUES ('Стройматериалы', 'Стройматериалы-1'),
       ('Игрушки', 'Игрушки-2'),
       ('Продукты', 'Продукты-3'),
       ('Тряпки', 'Тряпки-4'),
       ('Товаря для дома', 'Товары для дома-5');



