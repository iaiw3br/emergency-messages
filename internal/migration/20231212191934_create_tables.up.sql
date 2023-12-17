CREATE TABLE IF NOT EXISTS users
(
    id           uuid  PRIMARY KEY,
    first_name   VARCHAR(50) UNIQUE NOT NULL,
    last_name    VARCHAR(50)        NOT NULL,
    email        VARCHAR(255) UNIQUE,
    mobile_phone VARCHAR(50) UNIQUE,
    city         VARCHAR(100)       NOT NULL
);

CREATE TABLE IF NOT EXISTS messages
(
    id      uuid  PRIMARY KEY,
    status  message_type NOT NULL,
    subject VARCHAR(50) NOT NULL,
    text    text        NOT NULL,
    user_id uuid  NOT NULL,

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);
CREATE TYPE message_type AS ENUM ('created', 'delivered');

CREATE TABLE IF NOT EXISTS templates
(
    id uuid  PRIMARY KEY,
    subject VARCHAR(50) NOT NULL,
    text    text        NOT NULL
);
