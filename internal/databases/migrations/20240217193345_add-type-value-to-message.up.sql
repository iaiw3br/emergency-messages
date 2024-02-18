DROP TYPE message_type;
CREATE TYPE message_type AS ENUM ('sms', 'email');

ALTER TABLE public.messages
    ADD COLUMN IF NOT EXISTS type  message_type,
    ADD COLUMN IF NOT EXISTS value text;


