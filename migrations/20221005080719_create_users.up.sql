CREATE TABLE users (
    id bigserial not null primary key,
    username varchar unique,
    email varchar not null unique,
    encrypted_password varchar not null,
    given_name varchar,
    family_name varchar
);