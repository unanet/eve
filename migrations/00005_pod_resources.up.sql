create table pod_resources
(
    id serial not null,
    description varchar not null,
    data jsonb default '{}'::json not null,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null
);

create unique index pod_resources_description_uindex
    on pod_resources (description);

create unique index pod_resources_id_uindex
    on pod_resources (id);

alter table pod_resources
    add constraint pod_resources_pk
        primary key (id);

INSERT INTO pod_resources(description, data)
    VALUES
        ('default subcontractor' ,'{"limit": {"cpu": "1000m", "memory": "500M"}, "request": {"cpu": "250m", "memory": "50M"}}'),
        ('default unanet-app','{"limit": {"cpu": "1000m", "memory": "2500M"}, "request": {"cpu": "250m", "memory": "1650M"}}'),
        ('default unanet-analytics','{"limit": {"cpu": "1000m", "memory": "3000M"}, "request": {"cpu": "250m", "memory": "2000M"}}'),
        ('default platform','{"limit": {"cpu": "1000m", "memory": "2000M"}, "request": {"cpu": "250m", "memory": "1000M"}}'),
        ('default unanet-proxy','{"limit": {"cpu": "1000m", "memory": "250M"}, "request": {"cpu": "100m", "memory": "30M"}}'),
        ('limit 1GB','{"limit": {"memory": "1000M"}}'),
        ('limit 2GB','{"limit": {"memory": "2000M"}}'),
        ('limit 2.5GB','{"limit": {"memory": "2500M"}}'),
        ('limit 3GB','{"limit": {"memory": "3000M"}}'),
        ('limit 1CPU','{"limit": {"cpu": "1000m"}}'),
        ('request 1CPU','{"request": {"cpu": "1000m"}}'),
        ('request 1GB','{"request": {"memory": "1000M"}}'),
        ('request 2GB','{"request": {"memory": "2000M"}}'),
        ('request 2.5GB','{"request": {"memory": "2500M"}}'),
        ('request 3GB','{"request": {"memory": "3000M"}}');
