create table if not exists definition_type
(
    id serial not null,
    name varchar(50) not null,
    description varchar(250) not null,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null
);

create unique index if not exists definition_type_name_uindex
    on definition_type (name);

create unique index if not exists definition_type_description_uindex
    on definition_type (description);

create unique index if not exists definition_type_id_uindex
    on definition_type (id);

alter table definition_type
    add constraint definition_type_pk
        primary key (id);


INSERT INTO definition_type(name,description)
VALUES
('appsv1.Deployment','Kubernetes v1 Deployment Config'),
('batchv1.Job','Kubernetes v1 Job Config'),
('v2beta2.HorizontalPodAutoscaler','Kubernetes v2 Horizontal Pod Autoscaler Config'),
('apiv1.Service','Kubernetes v1 Service Config');


create table if not exists definition
(
    id serial not null,
    description varchar(200) not null,
    definition_type_id integer,
    data jsonb default '{}'::json not null,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    constraint definition_service_map_definition_id_fk
        foreign key (definition_type_id) references definition_type
            on update cascade on delete cascade
);

create unique index if not exists definition_description_uindex
    on definition (description);

create unique index if not exists definition_id_uindex
    on definition (id);

alter table definition
    add constraint definition_pk
        primary key (id);


INSERT INTO definition(description,definition_type_id,data)
VALUES
('deployment:nodeSelector:node-group:shared',1,'{"spec": {"template": {"spec": {"nodeSelector": {"node-group": "shared"}}}}}'),
('deployment:labels:proxy:unanet',1,'{"spec": {"template": {"metadata": {"labels": {"proxy": "unanet"}}}}}'),
('deployment:labels:proxy:disabled',1,'{"spec": {"template": {"metadata": {"labels": {"proxy": "disabled"}}}}}'),
('deployment:labels:proxy:uae',1,'{"spec": {"template": {"metadata": {"labels": {"proxy": "uae"}}}}}');








create table if not exists definition_service_map
(
    description    varchar(200)            not null,
    definition_id    integer                 not null,
    environment_id integer,
    artifact_id    integer,
    namespace_id   integer,
    service_id     integer,
    cluster_id     integer,
    stacking_order integer   default 0     not null,
    created_at     timestamp default now() not null,
    updated_at     timestamp default now() not null,
    constraint definition_service_map_definition_id_fk
        foreign key (definition_id) references definition
            on update cascade on delete cascade,
    constraint definition_service_map_environment_id_fk
        foreign key (environment_id) references environment
            on update cascade on delete cascade,
    constraint definition_service_map_artifact_id_fk
        foreign key (artifact_id) references artifact
            on update cascade on delete cascade,
    constraint definition_service_map_namespace_id_fk
        foreign key (namespace_id) references namespace
            on update cascade on delete cascade,
    constraint definition_service_map_service_id_fk
        foreign key (service_id) references service
            on update cascade on delete cascade,
    constraint definition_service_map_cluster_id_fk
        foreign key (cluster_id) references cluster
            on update cascade on delete cascade
);

create unique index if not exists definition_service_map_description_uindex
    on definition_service_map(description);

alter table definition_service_map
    add constraint chk_only_service_id_or_artifact_id_nonnull
        check (num_nonnulls(service_id, artifact_id) <= 1);

alter table definition_service_map
    add constraint chk_only_at_least_one_must_be_set
        check (num_nonnulls(service_id, cluster_id, environment_id, namespace_id, artifact_id) >= 1);

alter table definition_service_map
    add constraint chk_no_more_than_one_can_be_set
        check (num_nonnulls(service_id, cluster_id, environment_id, namespace_id) <= 1);



create table if not exists definition_job_map
(
    description    varchar(200)            not null,
    definition_id    integer                 not null,
    environment_id integer,
    artifact_id    integer,
    namespace_id   integer,
    job_id     integer,
    cluster_id     integer,
    stacking_order integer   default 0     not null,
    created_at     timestamp default now() not null,
    updated_at     timestamp default now() not null,
    constraint definition_job_map_definition_id_fk
        foreign key (definition_id) references definition
            on update cascade on delete cascade,
    constraint definition_job_map_environment_id_fk
        foreign key (environment_id) references environment
            on update cascade on delete cascade,
    constraint definition_job_map_artifact_id_fk
        foreign key (artifact_id) references artifact
            on update cascade on delete cascade,
    constraint definition_job_map_namespace_id_fk
        foreign key (namespace_id) references namespace
            on update cascade on delete cascade,
    constraint definition_job_map_job_id_fk
        foreign key (job_id) references job
            on update cascade on delete cascade,
    constraint definition_job_map_cluster_id_fk
        foreign key (cluster_id) references cluster
            on update cascade on delete cascade
);


create unique index if not exists definition_job_map_description_uindex
    on definition_job_map(description);

alter table definition_job_map
    add constraint chk_only_job_id_or_artifact_id_nonnull
        check (num_nonnulls(job_id, artifact_id) <= 1);

alter table definition_job_map
    add constraint chk_only_at_least_one_must_be_set
        check (num_nonnulls(job_id, cluster_id, environment_id, namespace_id, artifact_id) >= 1);

alter table definition_job_map
    add constraint chk_no_more_than_one_can_be_set
        check (num_nonnulls(job_id, cluster_id, environment_id, namespace_id) <= 1);

