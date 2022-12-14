CREATE TABLE `users` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT UNIQUE ,
    `username` VARCHAR(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL UNIQUE,
    `email` VARCHAR(50) NOT NULL UNIQUE ,
    `encrypted_password` VARCHAR(255) ,
    `given_name` VARCHAR(255) ,
    `family_name` VARCHAR(50) ,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `redacted_at` TIMESTAMP on update CURRENT_TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
    ) ENGINE = InnoDB;