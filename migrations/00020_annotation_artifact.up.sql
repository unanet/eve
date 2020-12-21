alter table annotation_job_map
    add artifact_id int;

alter table annotation_job_map
    add constraint annotation_job_map_artifact_id_fk
        foreign key (artifact_id) references artifact
            on update cascade on delete cascade;



alter table annotation_service_map
    add artifact_id int;

alter table annotation_service_map
    add constraint annotation_service_map_artifact_id_fk
        foreign key (artifact_id) references artifact
            on update cascade on delete cascade;