ALTER TABLE public.templates
    DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE public.users
    DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE IF EXISTS public.messages
    DROP COLUMN IF EXISTS deleted_at;