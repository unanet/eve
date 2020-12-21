alter table label_job_map
    add cluster_id int;

alter table label_job_map
    add constraint label_job_map_cluster_id_fk
        foreign key (cluster_id) references cluster
            on update cascade on delete cascade;



alter table label_service_map
    add cluster_id int;

alter table label_service_map
    add constraint label_service_map_cluster_id_fk
        foreign key (cluster_id) references cluster
            on update cascade on delete cascade;