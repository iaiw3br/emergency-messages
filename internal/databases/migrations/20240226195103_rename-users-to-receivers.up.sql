ALTER TABLE users RENAME TO receivers;

ALTER TABLE messages RENAME COLUMN user_id TO receiver_id;