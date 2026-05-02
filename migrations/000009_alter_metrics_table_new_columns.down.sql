ALTER TABLE metrics 
    DROP COLUMN IF EXISTS os_version,
    DROP COLUMN IF EXISTS browser_version,
    DROP COLUMN IF EXISTS device_brand,
    DROP COLUMN IF EXISTS device_model,
    DROP COLUMN IF EXISTS device;