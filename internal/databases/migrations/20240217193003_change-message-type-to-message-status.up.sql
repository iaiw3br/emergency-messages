CREATE TYPE public.message_status AS ENUM ('created', 'delivered');
ALTER TYPE public.message_status ADD VALUE 'failed';

ALTER TABLE public.messages
    ADD COLUMN new_status public.message_status;

UPDATE public.messages
SET new_status = status::TEXT::public.message_status;

ALTER TABLE public.messages
    DROP COLUMN status;

ALTER TABLE public.messages
    RENAME COLUMN new_status TO status;