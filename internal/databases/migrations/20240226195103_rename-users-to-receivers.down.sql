ALTER TABLE receivers RENAME TO users;

ALTER TABLE messages RENAME COLUMN receiver_id TO user_id;