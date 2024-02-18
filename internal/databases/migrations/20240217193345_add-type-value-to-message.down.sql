DROP TYPE public.message_type;
CREATE TYPE public.message_type AS ENUM ('created', 'delivered');

ALTER TABLE public.messages
    DROP COLUMN IF EXISTS type,
    DROP COLUMN IF EXISTS value;
