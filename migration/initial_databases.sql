CREATE TABLE users (
    id serial PRIMARY KEY,
    first_name VARCHAR ( 50 ) UNIQUE NOT NULL,
    last_name VARCHAR ( 50 ) NOT NULL,
    email VARCHAR ( 255 ) UNIQUE NOT NULL,
    mobile_phone VARCHAR ( 50 ) UNIQUE NOT NULL
);

drop table users;