alter table metadata_job_map
    add cluster_id int default null
    constraint metadata_job_map_cluster_id_fk
        references cluster
        on update cascade on delete cascade;

alter table metadata_service_map
    add cluster_id int default null
    constraint metadata_service_map_cluster_id_fk
        references cluster
        on update cascade on delete cascade;

alter table metadata_job_map
    drop constraint chk_only_at_least_one_must_be_set,
    drop constraint chk_only_service_id_or_environment_id_or_namespace_id_nonnull;

alter table metadata_service_map
    drop constraint chk_only_at_least_one_must_be_set,
    drop constraint chk_only_service_id_or_environment_id_or_namespace_id_nonnull;

alter table metadata_job_map
    add constraint chk_only_at_least_one_must_be_set
        check (num_nonnulls(job_id, cluster_id, environment_id, namespace_id, artifact_id) >= 1),
    add constraint chk_no_more_than_one_can_be_set
        check (num_nonnulls(job_id, cluster_id, environment_id, namespace_id) <= 1);

alter table metadata_service_map
    add constraint chk_only_at_least_one_must_be_set
        check (num_nonnulls(service_id, cluster_id, environment_id, namespace_id, artifact_id) >= 1),
    add constraint chk_no_more_than_one_can_be_set
        check (num_nonnulls(service_id, cluster_id, environment_id, namespace_id) <= 1);