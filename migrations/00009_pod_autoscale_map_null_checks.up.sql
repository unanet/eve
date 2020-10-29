ALTER TABLE pod_autoscale_map
    ADD CONSTRAINT chk_only_one_is_not_null CHECK (num_nonnulls(service_id, environment_id, namespace_id) = 1),
    ADD CONSTRAINT chk_two_are_null CHECK (num_nulls(service_id, environment_id, namespace_id) = 2);
