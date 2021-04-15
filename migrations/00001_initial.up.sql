CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create type feed_type as enum ('docker', 'generic');

create type provider_group as enum ('unanet', 'clearview', 'ops', 'cosential', 'connect');

create type deployment_state as enum ('queued', 'scheduled', 'completed');

create type deployment_cron_state as enum ('idle', 'running');

create type definition_order as enum ('main', 'pre', 'post');

create table if not exists feed
(
    id              integer               not null,
    name            varchar(25)           not null,
    promotion_order integer     default 0 not null,
    feed_type       feed_type,
    alias           varchar(25) default ''::character varying,
    constraint feed_pk
        primary key (id)
);

create unique index if not exists feed_name_uindex
    on feed (name);

create table if not exists artifact
(
    id               integer                                           not null,
    name             varchar(50)                                       not null,
    feed_type        feed_type                                         not null,
    provider_group   provider_group                                    not null,
    image_tag        varchar(25) default '$version'::character varying not null,
    service_port     integer     default 8080                          not null,
    metrics_port     integer     default 0                             not null,
    service_account  varchar(50) default 'unanet'::character varying   not null,
    run_as           integer     default 1101                          not null,
    liveliness_probe jsonb       default '{}'::json                    not null,
    readiness_probe  jsonb       default '{}'::json                    not null,
    constraint artifact_pk
        primary key (id)
);

create unique index if not exists artifact_name_uindex
    on artifact (name);

create table if not exists environment
(
    id          integer       not null,
    name        varchar(25)   not null,
    alias       varchar(25)   not null,
    description varchar(1024) not null,
    updated_at  timestamp default now(),
    constraint environment_pk
        primary key (id)
);

create unique index if not exists environment_name_uindex
    on environment (name);

create table if not exists cluster
(
    id             integer                 not null,
    name           varchar(50)             not null,
    sch_queue_url  varchar(200)            not null,
    provider_group provider_group          not null,
    created_at     timestamp default now() not null,
    updated_at     timestamp default now() not null,
    constraint cluster_pk
        primary key (id)
);

create unique index if not exists cluster_name_uindex
    on cluster (name);

create table if not exists namespace
(
    id                serial                  not null,
    alias             varchar(50)             not null,
    environment_id    integer                 not null,
    requested_version varchar(50)             not null,
    explicit_deploy   boolean   default false not null,
    cluster_id        integer                 not null,
    created_at        timestamp default now() not null,
    updated_at        timestamp default now() not null,
    name              varchar(50)             not null,
    constraint namespace_pk
        primary key (id),
    constraint namespace_environment_id_fk
        foreign key (environment_id) references environment,
    constraint namespace_cluster_id_fk
        foreign key (cluster_id) references cluster
);

create unique index if not exists namespace_environment_id_cluster_id_alias_uindex
    on namespace (environment_id, cluster_id, alias);

create unique index if not exists namespace_name_cluster_id_uindex
    on namespace (name, cluster_id);

create table if not exists service
(
    id                 serial                                         not null,
    namespace_id       integer                                        not null,
    artifact_id        integer                                        not null,
    override_version   varchar(50),
    deployed_version   varchar(50),
    created_at         timestamp    default now()                     not null,
    updated_at         timestamp    default now()                     not null,
    name               varchar(50)  default 'blah'::character varying not null,
    sticky_sessions    boolean      default false                     not null,
    count              integer      default 2                         not null,
    success_exit_codes varchar(100) default '0'::character varying    not null,
    explicit_deploy    boolean      default false                     not null,
    constraint service_pk
        primary key (id),
    constraint service_artifact_id_fk
        foreign key (artifact_id) references artifact,
    constraint service_namespace_id_fk
        foreign key (namespace_id) references namespace
);

create unique index if not exists service_namespace_id_name_uindex
    on service (name, namespace_id);

create table if not exists environment_feed_map
(
    environment_id integer not null,
    feed_id        integer not null,
    constraint environment_feed_map_environment_id
        foreign key (environment_id) references environment,
    constraint environment_feed_map_feed_id
        foreign key (feed_id) references feed
);

create unique index if not exists environment_feed_map_environment_id_feed_id_uindex
    on environment_feed_map (environment_id, feed_id);

create table if not exists deployment
(
    id             uuid      default uuid_generate_v4() not null,
    environment_id integer                              not null,
    namespace_id   integer                              not null,
    req_id         varchar(100),
    message_id     varchar(100),
    receipt_handle varchar(1024),
    plan_options   jsonb                                not null,
    plan_location  jsonb,
    state          deployment_state                     not null,
    "user"         varchar(50)                          not null,
    created_at     timestamp default now()              not null,
    updated_at     timestamp default now()              not null,
    constraint deployment_pkey
        primary key (id),
    constraint deployment_environment_id
        foreign key (environment_id) references environment,
    constraint deployment_namespace_id
        foreign key (namespace_id) references namespace
);

create table if not exists deployment_cron
(
    id           uuid                  default uuid_generate_v4()            not null,
    plan_options jsonb                                                       not null,
    schedule     varchar(25)                                                 not null,
    state        deployment_cron_state default 'idle'::deployment_cron_state not null,
    last_run     timestamp             default now()                         not null,
    disabled     boolean               default false                         not null,
    description  varchar(100),
    exec_order   integer               default 0                             not null,
    constraint deployment_cron_pkey
        primary key (id)
);

create table if not exists deployment_cron_job
(
    deployment_cron_id uuid not null,
    deployment_id      uuid not null,
    constraint deployment_cron_job_deployment_id
        foreign key (deployment_id) references deployment
            on delete cascade,
    constraint deployment_cron_job_deployment_cron_id
        foreign key (deployment_cron_id) references deployment_cron
            on delete cascade
);

create table if not exists job
(
    id                 serial                                      not null,
    name               varchar(50)                                 not null,
    artifact_id        integer                                     not null,
    namespace_id       integer,
    override_version   varchar(50),
    deployed_version   varchar(50),
    created_at         timestamp    default now()                  not null,
    updated_at         timestamp    default now()                  not null,
    success_exit_codes varchar(100) default '0'::character varying not null,
    explicit_deploy    boolean      default false                  not null,
    constraint job_pk
        primary key (id),
    constraint job_artifact_id_fk
        foreign key (artifact_id) references artifact,
    constraint job_namespace_id_fk
        foreign key (namespace_id) references namespace
);

create unique index if not exists job_namespace_id_name_uindex
    on job (name, namespace_id);

create table if not exists metadata
(
    id            serial                       not null,
    description   varchar(100)                 not null,
    value         jsonb     default '{}'::json not null,
    created_at    timestamp default now()      not null,
    updated_at    timestamp default now()      not null,
    migrated_from integer,
    constraint metadata_pk
        primary key (id)
);

create unique index if not exists metadata_description_uindex
    on metadata (description);

create unique index if not exists metadata_id_uindex
    on metadata (id);

create table if not exists metadata_job_map
(
    description    varchar(100)            not null,
    metadata_id    integer                 not null,
    environment_id integer,
    artifact_id    integer,
    namespace_id   integer,
    job_id         integer,
    stacking_order integer   default 0     not null,
    created_at     timestamp default now() not null,
    updated_at     timestamp default now() not null,
    cluster_id     integer,
    constraint metadata_job_map_metadata_id_fk
        foreign key (metadata_id) references metadata
            on update cascade on delete cascade,
    constraint metadata_job_map_environment_id_fk
        foreign key (environment_id) references environment
            on update cascade on delete cascade,
    constraint metadata_job_map_artifact_id_fk
        foreign key (artifact_id) references artifact
            on update cascade on delete cascade,
    constraint metadata_job_map_namespace_id_fk
        foreign key (namespace_id) references namespace
            on update cascade on delete cascade,
    constraint metadata_job_map_job_id_fk
        foreign key (job_id) references job
            on update cascade on delete cascade,
    constraint metadata_job_map_cluster_id_fk
        foreign key (cluster_id) references cluster
            on update cascade on delete cascade
);

create unique index if not exists metadata_job_map_description_uindex
    on metadata_job_map (description);

alter table metadata_job_map
    add constraint chk_only_service_id_or_artifact_id_nonnull
        check (num_nonnulls(job_id, artifact_id) <= 1);

alter table metadata_job_map
    add constraint chk_only_at_least_one_must_be_set
        check (num_nonnulls(job_id, cluster_id, environment_id, namespace_id, artifact_id) >= 1);

alter table metadata_job_map
    add constraint chk_no_more_than_one_can_be_set
        check (num_nonnulls(job_id, cluster_id, environment_id, namespace_id) <= 1);

create table if not exists metadata_service_map
(
    description    varchar(100)            not null,
    metadata_id    integer                 not null,
    environment_id integer,
    artifact_id    integer,
    namespace_id   integer,
    service_id     integer,
    stacking_order integer   default 0     not null,
    created_at     timestamp default now() not null,
    updated_at     timestamp default now() not null,
    cluster_id     integer,
    constraint metadata_service_map_metadata_id_fk
        foreign key (metadata_id) references metadata
            on update cascade on delete cascade,
    constraint metadata_service_map_environment_id_fk
        foreign key (environment_id) references environment
            on update cascade on delete cascade,
    constraint metadata_service_map_artifact_id_fk
        foreign key (artifact_id) references artifact
            on update cascade on delete cascade,
    constraint metadata_service_map_namespace_id_fk
        foreign key (namespace_id) references namespace
            on update cascade on delete cascade,
    constraint metadata_service_map_service_id_fk
        foreign key (service_id) references service
            on update cascade on delete cascade,
    constraint metadata_service_map_cluster_id_fk
        foreign key (cluster_id) references cluster
            on update cascade on delete cascade
);

create unique index if not exists metadata_service_map_description_uindex
    on metadata_service_map (description);

alter table metadata_service_map
    add constraint chk_only_service_id_or_artifact_id_nonnull
        check (num_nonnulls(service_id, artifact_id) <= 1);

alter table metadata_service_map
    add constraint chk_only_at_least_one_must_be_set
        check (num_nonnulls(service_id, cluster_id, environment_id, namespace_id, artifact_id) >= 1);

alter table metadata_service_map
    add constraint chk_no_more_than_one_can_be_set
        check (num_nonnulls(service_id, cluster_id, environment_id, namespace_id) <= 1);

create table if not exists metadata_history
(
    metadata_id integer      not null,
    description varchar(100) not null,
    value       jsonb        not null,
    created     timestamp,
    created_by  varchar(32),
    deleted     timestamp,
    deleted_by  varchar(32)
);

create table if not exists definition_type
(
    id               serial                  not null,
    name             varchar(50)             not null,
    description      varchar(250)            not null,
    created_at       timestamp default now() not null,
    updated_at       timestamp default now() not null,
    class            varchar(50)             not null,
    version          varchar(50)             not null,
    kind             varchar(50)             not null,
    definition_order definition_order        not null,
    constraint definition_type_pk
        primary key (id)
);

create unique index if not exists definition_type_name_uindex
    on definition_type (name);

create unique index if not exists definition_type_description_uindex
    on definition_type (description);

create unique index if not exists definition_type_id_uindex
    on definition_type (id);

create table if not exists definition
(
    id                 serial                       not null,
    description        varchar(200)                 not null,
    definition_type_id integer,
    data               jsonb     default '{}'::json not null,
    created_at         timestamp default now()      not null,
    updated_at         timestamp default now()      not null,
    constraint definition_pk
        primary key (id),
    constraint definition_service_map_definition_id_fk
        foreign key (definition_type_id) references definition_type
            on update cascade on delete cascade
);

create unique index if not exists definition_description_uindex
    on definition (description);

create unique index if not exists definition_id_uindex
    on definition (id);

create table if not exists definition_service_map
(
    description    varchar(200)            not null,
    definition_id  integer                 not null,
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
    on definition_service_map (description);

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
    definition_id  integer                 not null,
    environment_id integer,
    artifact_id    integer,
    namespace_id   integer,
    job_id         integer,
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
    on definition_job_map (description);

alter table definition_job_map
    add constraint chk_only_job_id_or_artifact_id_nonnull
        check (num_nonnulls(job_id, artifact_id) <= 1);

alter table definition_job_map
    add constraint chk_only_at_least_one_must_be_set
        check (num_nonnulls(job_id, cluster_id, environment_id, namespace_id, artifact_id) >= 1);

alter table definition_job_map
    add constraint chk_no_more_than_one_can_be_set
        check (num_nonnulls(job_id, cluster_id, environment_id, namespace_id) <= 1);

create function jsonb_merge(orig jsonb, delta jsonb) returns jsonb
    language sql
as
$$
SELECT jsonb_object_agg(
               coalesce(keyOrig, keyDelta),
               CASE
                   WHEN valOrig ISNULL THEN valDelta
                   WHEN valDelta ISNULL THEN valOrig
                   WHEN (jsonb_typeof(valOrig) <> 'object' OR jsonb_typeof(valDelta) <> 'object') THEN valDelta
                   ELSE jsonb_merge(valOrig, valDelta)
                   END
           )
FROM jsonb_each(orig) e1(keyOrig, valOrig)
         FULL JOIN jsonb_each(delta) e2(keyDelta, valDelta) ON keyOrig = keyDelta
$$;

create function metadata_insert() returns trigger
    language plpgsql
as
$$
BEGIN
    INSERT INTO metadata_history
    (metadata_id, description, value, created, created_by)
    VALUES (NEW.id, NEW.description, NEW.value, current_timestamp, current_user);
    RETURN NEW;
END;
$$;

create trigger metadata_insert_trigger
    after insert
    on metadata
    for each row
execute procedure metadata_insert();

create function metadata_delete() returns trigger
    language plpgsql
as
$$
BEGIN
    UPDATE metadata_history
    SET deleted    = current_timestamp,
        deleted_by = current_user
    WHERE deleted IS NULL
      and metadata_id = OLD.id;
    RETURN NULL;
END;
$$;

create trigger metadata_delete_trigger
    after delete
    on metadata
    for each row
execute procedure metadata_delete();

create function metadata_update() returns trigger
    language plpgsql
as
$$
BEGIN

    UPDATE metadata_history
    SET deleted    = current_timestamp,
        deleted_by = current_user
    WHERE deleted IS NULL
      and metadata_id = OLD.id;

    INSERT INTO metadata_history
    (metadata_id, description, value, created, created_by)
    VALUES (NEW.id, NEW.description, NEW.value, current_timestamp, current_user);
    RETURN NEW;

END;
$$;

create trigger metadata_update_trigger
    after update
    on metadata
    for each row
execute procedure metadata_update();
