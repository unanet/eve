create table metadata
(
    id serial not null
        constraint metadata_pk
            primary key,
    value jsonb default '{}'::json not null,
    description varchar(100) not null
);

create unique index metadata_description_uindex
    on metadata (description);

create unique index metadata_id_uindex
    on metadata (id);

create table metadata_job_map
(
    job_id integer
        constraint metadata_job_map_job_id_fk
            references job
            on update cascade on delete cascade,
    environment_id integer
        constraint metadata_job_map_environment_id_fk
            references environment
            on update cascade on delete cascade,
    namespace_id integer
        constraint metadata_job_map_namespace_id_fk
            references namespace
            on update cascade on delete cascade,
    artifact_id integer
        constraint metadata_job_map_artifact_id_fk
            references artifact
            on update cascade on delete cascade,
    metadata_id integer not null
        constraint metadata_job_map_metadata_id_fk
            references metadata
            on update cascade on delete cascade,
    stacking_order integer default 0 not null,
    description varchar(100) not null
);

create table metadata_service_map
(
    service_id integer
        constraint metadata_service_map_service_id_fk
            references service
            on update cascade on delete cascade,
    namespace_id integer
        constraint metadata_service_map_namespace_id_fk
            references namespace
            on update cascade on delete cascade,
    environment_id integer
        constraint metadata_service_map_environment_id_fk
            references environment
            on update cascade on delete cascade,
    metadata_id integer not null
        constraint metadata_service_map_metadata_id_fk
            references metadata
            on update cascade on delete cascade,
    artifact_id integer
        constraint metadata_service_map_artifact_id_fk
            references artifact
            on update cascade on delete cascade,
    description varchar(100) not null,
    stacking_order integer default 0 not null
);