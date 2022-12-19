CREATE TABLE `lots` (
    `lot_id` INT UNSIGNED NOT NULL AUTO_INCREMENT ,
    `user_id` INT UNSIGNED NOT NULL,
    `type_of_estate` VARCHAR(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
    `rooms` INT(2) NOT NULL ,
    `area` INT(4) NOT NULL ,
    `floor` INT(3) NOT NULL ,
    `max_floor` INT(3),
    `city` VARCHAR(255) NOT NULL ,
    `district` VARCHAR(255) NOT NULL ,
    `street` VARCHAR(255) NOT NULL ,
    `building` VARCHAR(50) NOT NULL ,
    `price` INT NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `redacted_at` TIMESTAMP on update CURRENT_TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`lot_id`),
    FOREIGN KEY (`user_id`) REFERENCES users(user_id)
    ) ENGINE = InnoDB;