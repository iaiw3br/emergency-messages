ALTER TABLE public.templates
    ADD COLUMN deleted_at timestamp;
ALTER TABLE public.users
    ADD COLUMN deleted_at timestamp;
ALTER TABLE public.messages
    ADD COLUMN deleted_at timestamp;