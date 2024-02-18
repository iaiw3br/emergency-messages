ALTER TABLE public.messages
    ALTER COLUMN status TYPE text USING status::message_type;

DROP TYPE message_type;