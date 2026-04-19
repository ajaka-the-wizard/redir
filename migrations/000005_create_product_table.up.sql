CREATE TABLE IF NOT EXISTS products(
    id SERIAL PRIMARY KEY,
    product_id INTEGER UNIQUE NOT NULL DEFAULT pseudo_encrypt(nextval('products_id_seq')::int),
    product_name TEXT NOT NULL,
    public BOOLEAN NOT NULL DEFAULT TRUE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    private_key TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
