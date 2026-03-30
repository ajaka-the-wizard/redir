CREATE TABLE IF NOT EXISTS metrics(
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    media_id UUID REFERENCES medias(id) ON DELETE CASCADE,
    upload_duration_ms INTEGER,
    error_message TEXT,
    captured_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);