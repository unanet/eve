create table pod_resources_map
(
    description varchar(100) not null,
    pod_resources_id integer not null
        constraint pod_resources_map_pod_resources_id_fk
            references pod_resources
            on update cascade on delete cascade,
    service_id integer
        constraint pod_resources_map_service_id_fk
            references service
            on update cascade on delete cascade,
    environment_id integer
        constraint pod_resources_map_environment_id_fk
            references environment
            on update cascade on delete cascade,
    namespace_id integer
        constraint pod_resources_map_namespace_id_fk
            references namespace
            on update cascade on delete cascade,
    stacking_order integer default 0 not null
);

create unique index pod_resources_stacking_order_uindex
    on pod_resources_map (pod_resources_id, service_id, environment_id, namespace_id, stacking_order);

