CREATE TABLE
    users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        email VARCHAR(255) UNIQUE NOT NULL,
        phone_number VARCHAR(20) UNIQUE NOT NULL,
        sex int NOT NULL,
        address text NOT NULL,
        birth_date time,
        
    )ENGINE=INNODB DEFAULT CHARSET=utf8mb4;

CREATE TABLE hobbies (
        id INT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(255) NOT NULL
    )ENGINE=INNODB DEFAULT CHARSET=utf8mb4;

CREATE TABLE
    user_hobbies (
        id INT AUTO_INCREMENT PRIMARY KEY,
        user_id int NOT NULL,
        hobby_id int UNIQUE NOT NULL,
        CONSTRAINT `fk_user_id`
            FOREIGN KEY (`user_id`)
            REFERENCES `users` (`id`)
            ON DELETE CASCADE
            ON UPDATE CASCADE,
        CONSTRAINT `fk_hobby_id`
            FOREIGN KEY (`hobby_id`)
            REFERENCES `hobbies` (`id`)
            ON UPDATE CASCADE
    )ENGINE=INNODB DEFAULT CHARSET=utf8mb4;
