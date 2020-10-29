ALTER TABLE service
    DROP COLUMN IF EXISTS autoscaling,
    DROP COLUMN IF EXISTS pod_resource;
