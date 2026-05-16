CREATE TABLE IF NOT EXISTS metrics(
    id BIGSERIAL PRIMARY KEY,
    media_id TEXT REFERENCES medias(public_key) ON DELETE CASCADE,
    ip INET,
    browser TEXT,
    os TEXT,
    country CHAR(2),
    referrer TEXT,
    captured_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);


CREATE INDEX idx_metrics_cap_at ON metrics(captured_at);
CREATE INDEX idx_metrics_time ON metrics(media_id, captured_at DESC);