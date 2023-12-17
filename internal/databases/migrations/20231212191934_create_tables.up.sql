CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS public.users
(
    id      UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    first_name   VARCHAR(50) UNIQUE NOT NULL,
    last_name    VARCHAR(50)        NOT NULL,
    email        VARCHAR(255) UNIQUE,
    mobile_phone VARCHAR(50) UNIQUE,
    city         VARCHAR(100)       NOT NULL
);
CREATE TYPE public.message_type AS ENUM ('created', 'delivered');

CREATE TABLE IF NOT EXISTS public.messages
(
    id      UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    status  message_type NOT NULL,
    subject VARCHAR(50)  NOT NULL,
    text    text         NOT NULL,
    user_id uuid         NOT NULL,

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);


CREATE TABLE IF NOT EXISTS public.templates
(
    id      UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    subject VARCHAR(50) NOT NULL,
    text    text        NOT NULL
);
