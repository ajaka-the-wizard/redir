CREATE TYPE media_type AS ENUM ('image','video');
CREATE TYPE upload_status AS ENUM ('pending','completed','failed');

CREATE TABLE IF NOT EXISTS medias(
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    bucket_key TEXT UNIQUE NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    bucket TEXT NOT NULL,
    file_size BIGINT,
    status upload_status DEFAULT 'pending',
    file_type media_type NOT NULL,
    active BOOLEAN DEFAULT TRUE,
    file_name TEXT NOT NULL,
    mime_type TEXT,
    hits INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);