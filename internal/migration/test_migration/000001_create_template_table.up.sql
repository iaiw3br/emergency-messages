CREATE TABLE users
(
    id           serial PRIMARY KEY,
    first_name   VARCHAR(50) UNIQUE NOT NULL,
    last_name    VARCHAR(50)        NOT NULL,
    email        VARCHAR(255) UNIQUE,
    mobile_phone VARCHAR(50) UNIQUE,
    city         VARCHAR(100)       NOT NULL
);

CREATE TABLE messages
(
    id      serial PRIMARY KEY,
    status  VARCHAR(50) NOT NULL,
    subject VARCHAR(50) NOT NULL,
    text    text        NOT NULL,
    user_id int         NOT NULL,

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE templates
(
    id      serial PRIMARY KEY,
    subject VARCHAR(50) NOT NULL,
    text    text        NOT NULL
);