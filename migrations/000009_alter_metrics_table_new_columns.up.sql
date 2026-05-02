ALTER TABLE metrics 
    ADD COLUMN os_version TEXT,
    ADD COLUMN browser_version TEXT,
    ADD COLUMN device_brand TEXT,
    ADD COLUMN device_model TEXT,
    ADD COLUMN device TEXT;
