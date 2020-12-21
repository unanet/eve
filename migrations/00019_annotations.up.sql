create table annotation
(
    id serial not null,
    description varchar not null,
    data jsonb default '{}'::json not null,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null
);

create unique index annotation_description_uindex
    on annotation (description);

create unique index annotation_id_uindex
    on annotation (id);

alter table annotation
    add constraint annotation_pk
        primary key (id);

INSERT INTO annotation(description, data)
VALUES
('applinks','{"proxy.unanet.io/paths":"/.well-known/apple-app-site-association,/.well-known/assetlinks.json", "proxy.unanet.io/exact":"true"}'),
('unanet','{"proxy.unanet.io/paths":"unanet"}'),
('unanet-analytics','{"proxy.unanet.io/paths":"analytics,analytics-api"}');




create table annotation_service_map
(
    description varchar(100) not null,
    annotation_id integer not null
        constraint annotation_service_map_label_id_fk
            references annotation
            on update cascade on delete cascade,
    service_id integer
        constraint annotation_service_map_service_id_fk
            references service
            on update cascade on delete cascade,
    environment_id integer
        constraint annotation_service_map_environment_id_fk
            references environment
            on update cascade on delete cascade,
    namespace_id integer
        constraint annotation_service_map_namespace_id_fk
            references namespace
            on update cascade on delete cascade,
    stacking_order integer default 0 not null
);

create unique index annotation_service_map_stacking_order_uindex
    on annotation_service_map (annotation_id, service_id, environment_id, namespace_id, stacking_order);


ALTER TABLE annotation_service_map
    ADD CONSTRAINT chk_only_one_is_not_null CHECK (num_nonnulls(service_id, environment_id, namespace_id) = 1),
    ADD CONSTRAINT chk_two_are_null CHECK (num_nulls(service_id, environment_id, namespace_id) = 2);



create table annotation_job_map
(
    description varchar(100) not null,
    annotation_id integer not null
        constraint annotation_job_map_label_id_fk
            references annotation
            on update cascade on delete cascade,
    job_id integer
        constraint annotation_job_map_service_id_fk
            references job
            on update cascade on delete cascade,
    environment_id integer
        constraint annotation_job_map_environment_id_fk
            references environment
            on update cascade on delete cascade,
    namespace_id integer
        constraint annotation_job_map_namespace_id_fk
            references namespace
            on update cascade on delete cascade,
    stacking_order integer default 0 not null
);

create unique index annotation_job_map_stacking_order_uindex
    on annotation_job_map (annotation_id, job_id, environment_id, namespace_id, stacking_order);


ALTER TABLE annotation_job_map
    ADD CONSTRAINT chk_only_one_is_not_null CHECK (num_nonnulls(job_id, environment_id, namespace_id) = 1),
    ADD CONSTRAINT chk_two_are_null CHECK (num_nulls(job_id, environment_id, namespace_id) = 2);
