ALTER TABLE environment ADD COLUMN updated_at TIMESTAMP;
ALTER TABLE environment ALTER COLUMN updated_at SET DEFAULT now();