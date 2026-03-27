CREATE TABLE IF NOT EXISTS products(
    id SERIAL PRIMARY KEY,
    product_id INTEGER UNIQUE NOT NULL DEFAULT pseudo_encrypt(nextval('product_id_seq')::int),
    product_name TEXT NOT NULL,
    user_id UUID REFERENCES users(id) NOT NULL ON DELETE CASCADE,
    private_key TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- CREATE OR REPLACE FUNCTION pseudo_encrypt(VALUE int) RETURNS int AS $$
-- DECLARE 
--     l1 int;
--     l2 int;
--     r1 int;
--     r2 int;
--     i int:=0;
-- BEGIN
--     l1:=(VALUE >> 16) & 65535;
--     r1:= VALUE & 65535;
--     WHILE i < 3 LOOP
--         l2 := r1;
--         r2 := l1;
--         l1 := l2;
--         r1 := r2;
--         i := i + 1;
--     END LOOP;
--     RETURN ((l1 << 16) | r1);
-- END;
-- $$ LANGUAGE plpgsql STRICT IMMUTABLE;