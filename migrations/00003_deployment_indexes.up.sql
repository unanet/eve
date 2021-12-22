CREATE INDEX IF NOT EXISTS idx_deployment_namespace_id ON deployment(namespace_id);
CREATE INDEX IF NOT EXISTS idx_deployment_namespace_id ON deployment(environment_id);