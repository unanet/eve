create table label
(
    id serial not null,
    description varchar not null,
    data jsonb default '{}'::json not null,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null
);

create unique index label_description_uindex
    on label (description);

create unique index label_id_uindex
    on label (id);

alter table label
    add constraint label_pk
        primary key (id);

INSERT INTO label(description, data)
VALUES
('unanet-proxy','{"proxy":"unanet"}');



create table label_service_map
(
    description varchar(100) not null,
    label_id integer not null
        constraint label_service_map_label_id_fk
            references label
            on update cascade on delete cascade,
    service_id integer
        constraint label_service_map_service_id_fk
            references service
            on update cascade on delete cascade,
    environment_id integer
        constraint label_service_map_environment_id_fk
            references environment
            on update cascade on delete cascade,
    namespace_id integer
        constraint label_service_map_namespace_id_fk
            references namespace
            on update cascade on delete cascade,
    stacking_order integer default 0 not null
);

create unique index label_service_map_stacking_order_uindex
    on label_service_map (label_id, service_id, environment_id, namespace_id, stacking_order);


ALTER TABLE label_service_map
    ADD CONSTRAINT chk_only_one_is_not_null CHECK (num_nonnulls(service_id, environment_id, namespace_id) = 1),
    ADD CONSTRAINT chk_two_are_null CHECK (num_nulls(service_id, environment_id, namespace_id) = 2);



create table label_job_map
(
    description varchar(100) not null,
    label_id integer not null
        constraint label_job_map_label_id_fk
            references label
            on update cascade on delete cascade,
    job_id integer
        constraint label_job_map_service_id_fk
            references job
            on update cascade on delete cascade,
    environment_id integer
        constraint label_job_map_environment_id_fk
            references environment
            on update cascade on delete cascade,
    namespace_id integer
        constraint label_job_map_namespace_id_fk
            references namespace
            on update cascade on delete cascade,
    stacking_order integer default 0 not null
);

create unique index label_job_map_stacking_order_uindex
    on label_job_map (label_id, job_id, environment_id, namespace_id, stacking_order);


ALTER TABLE label_job_map
    ADD CONSTRAINT chk_only_one_is_not_null CHECK (num_nonnulls(job_id, environment_id, namespace_id) = 1),
    ADD CONSTRAINT chk_two_are_null CHECK (num_nulls(job_id, environment_id, namespace_id) = 2);






































