alter table pod_resources_map
    add artifact_id integer;

alter table pod_resources_map
    add constraint pod_resources_map_artifact_id_fk
        foreign key (artifact_id) references artifact
            on update cascade on delete cascade;