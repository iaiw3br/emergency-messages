ALTER TABLE users
    ADD COLUMN contacts jsonb;

UPDATE users
SET contacts = jsonb_build_array(
        jsonb_build_object('type', 'email', 'value', email, 'is_active', true),
        jsonb_build_object('type', 'mobile_phone', 'value', mobile_phone, 'is_active', false)
               );

ALTER TABLE users
    DROP COLUMN email;

ALTER TABLE users
    DROP COLUMN mobile_phone;