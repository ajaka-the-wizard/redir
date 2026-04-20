CREATE TYPE media_type AS ENUM ('image','video');
CREATE TYPE upload_status AS ENUM ('pending','completed','failed');

CREATE TABLE IF NOT EXISTS medias(
    public_key Text PRIMARY KEY UNIQUE NOT NULL,
    inner_key TEXT UNIQUE NOT NULL,
    batch_id UUID NOT NULL,
    seq_id int NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    file_size BIGINT,
    status upload_status DEFAULT 'pending',
    file_type media_type NOT NULL,
    active BOOLEAN DEFAULT TRUE,
    public BOOLEAN NOT NULL,
    file_name TEXT NOT NULL,
    mime_type TEXT,
    hits INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (batch_id, seq_id)
);


CREATE INDEX idx_medias_batch_id ON medias(batch_id)