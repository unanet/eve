alter table annotation_job_map
    add cluster_id int;

alter table annotation_job_map
    add constraint annotation_job_map_cluster_id_fk
        foreign key (cluster_id) references cluster
            on update cascade on delete cascade;



alter table annotation_service_map
    add cluster_id int;

alter table annotation_service_map
    add constraint annotation_service_map_cluster_id_fk
        foreign key (cluster_id) references cluster
            on update cascade on delete cascade;