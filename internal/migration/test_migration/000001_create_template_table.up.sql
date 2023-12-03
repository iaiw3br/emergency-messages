CREATE TABLE IF NOT EXISTS templates
(
    id      serial PRIMARY KEY,
    subject VARCHAR(50) NOT NULL,
    text    text        NOT NULL
);