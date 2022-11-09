CREATE TABLE `u1803158_default`.`users` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT ,
    `username` VARCHAR(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
    `email` VARCHAR(50) NOT NULL ,
    `encrypted_password` VARCHAR(255) ,
    `given_name` VARCHAR(255) ,
    `family_name` VARCHAR(50) ,
    `created_at` TIMESTAMP NOT NULL,
    `redacted_at` TIMESTAMP,
    PRIMARY KEY (`id`)
    ) ENGINE = InnoDB;