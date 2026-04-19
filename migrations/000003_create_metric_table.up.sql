CREATE TABLE IF NOT EXISTS metrics(
    id BIGSERIAL PRIMARY KEY,
    media_id UUID REFERENCES medias(public_key) ON DELETE CASCADE,
    ip INET,
    browser TEXT,
    os TEXT,
    country CHAR(3),
    referrer TEXT,
    captured_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);


CREATE INDEX idx_metrics ON metrics(media_id);
CREATE INDEX idx_metrics_req_at ON metrics(captured_at);
CREATE INDEX idx_metrics_time ON metrics(media_id, captured_at DESC);