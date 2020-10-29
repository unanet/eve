ALTER TABLE artifact
    DROP COLUMN IF EXISTS autoscaling,
    DROP COLUMN IF EXISTS pod_resource;
