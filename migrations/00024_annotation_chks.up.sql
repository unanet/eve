ALTER TABLE annotation_service_map
    drop constraint chk_only_one_is_not_null,
    drop constraint chk_two_are_null;


ALTER TABLE annotation_job_map
    drop constraint chk_only_one_is_not_null,
    drop constraint chk_two_are_null;


ALTER TABLE annotation_service_map
    add constraint chk_only_at_least_one_must_be_set
        check (num_nonnulls(service_id, cluster_id, environment_id, namespace_id, artifact_id) >= 1),
    add constraint chk_no_more_than_one_can_be_set
        check (num_nonnulls(service_id, cluster_id, environment_id, namespace_id) <= 1);

alter table annotation_job_map
    add constraint chk_only_at_least_one_must_be_set
        check (num_nonnulls(job_id, cluster_id, environment_id, namespace_id, artifact_id) >= 1),
    add constraint chk_no_more_than_one_can_be_set
        check (num_nonnulls(job_id, cluster_id, environment_id, namespace_id) <= 1);