INSERT INTO
    users (id, name, email, phone_number, sex, address)
VALUES (
        1,
        'taro yamada',
        'taro-yamada@example.com',
        'xxx-xxxx-xxxx',
        0,
        "東京都新宿区大久保"
        ""
    ), (
        2,
        'hanako sato',
        'sato-hanako@example.com',
        'yyy-yyyy-yyyy',
        1,
        "東京都港区芝浦"
        ""
    );
INSERT INTO
    hobbies (id, name)
VALUES (
        1,
        "スポーツ"        
    ), (
        2,
        "パソコン"
    ),(
        3,
        "読書"
    );
INSERT INTO
    user_hobbies (id, user_id, hobby_id)
VALUES (
        1,
        1,
        1
    ), (
        2,
        1,
        2
    ),(
        3,
        2,
        3
);


