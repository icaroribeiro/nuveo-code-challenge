-- Enable pgcrypto extension to generate UUID value.
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create types
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status') THEN
        CREATE TYPE status AS ENUM ('inserted', 'consumed');
    END IF;
    -- More types here...
END$$;

CREATE TABLE IF NOT EXISTS workflows (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    status status DEFAULT 'inserted',
    data JSONB NOT NULL,
    steps VARCHAR (255) ARRAY NOT NULL
);