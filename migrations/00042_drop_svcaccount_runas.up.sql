ALTER TABLE artifact
DROP COLUMN IF EXISTS service_account,
DROP COLUMN IF EXISTS run_as,
DROP COLUMN IF EXISTS readiness_probe,
DROP COLUMN IF EXISTS liveliness_probe;