ALTER TABLE public.templates
    ADD COLUMN IF NOT EXISTS created_at timestamp,
    ADD COLUMN IF NOT EXISTS updated_at timestamp;
