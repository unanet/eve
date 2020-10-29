create table pod_autoscale
(
    id serial not null,
    description varchar not null,
    data jsonb default '{}'::json not null,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null
);

create unique index pod_autoscale_description_uindex
    on pod_autoscale (description);

create unique index pod_autoscale_id_uindex
    on pod_autoscale (id);

alter table pod_autoscale
    add constraint pod_autoscale_pk
        primary key (id);

INSERT INTO pod_autoscale(description, data)
    VALUES
        ('default','{"enabled": true, "replicas": {"max": 10, "min": 2}, "utilization": {"cpu": 70, "memory": 110}}');