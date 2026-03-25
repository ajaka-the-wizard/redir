CREATE TYPE provider_enum AS ENUM ('google','github','custom');

CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    full_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    provider provider_enum,
    provider_sub TEXT,
    verified BOOLEAN DEFAULT FALSE,
    paid BOOLEAN DEFAULT FALSE,
    active BOOLEAN DEFAULT TRUE,
    admin BOOLEAN DEFAULT FALSE,
    password TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);