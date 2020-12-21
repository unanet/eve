alter table label_job_map
    add artifact_id int;

alter table label_job_map
    add constraint label_job_map_artifact_id_fk
        foreign key (artifact_id) references artifact
            on update cascade on delete cascade;



alter table label_service_map
    add artifact_id int;

alter table label_service_map
    add constraint label_service_map_artifact_id_fk
        foreign key (artifact_id) references artifact
            on update cascade on delete cascade;