create table metadata
(
    id serial not null
        constraint metadata_pk
            primary key,
    description varchar(100) not null,
    value jsonb default '{}'::json not null,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    migrated_from int default null
);

create unique index metadata_description_uindex
    on metadata (description);

create unique index metadata_id_uindex
    on metadata (id);

create table metadata_job_map
(
    description varchar(100) not null,
    metadata_id integer not null
        constraint metadata_job_map_metadata_id_fk
            references metadata
            on update cascade on delete cascade,
    environment_id integer
        constraint metadata_job_map_environment_id_fk
            references environment
            on update cascade on delete cascade,
    artifact_id integer
        constraint metadata_job_map_artifact_id_fk
            references artifact
            on update cascade on delete cascade,
    namespace_id integer
        constraint metadata_job_map_namespace_id_fk
            references namespace
            on update cascade on delete cascade,
    job_id integer
        constraint metadata_job_map_job_id_fk
            references job
            on update cascade on delete cascade,
    stacking_order integer default 0 not null,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null
    constraint chk_only_at_least_one_must_be_set
        check (num_nonnulls(job_id, environment_id, namespace_id, artifact_id) >= 1),
    constraint chk_only_service_id_or_environment_id_or_namespace_id_nonnull
        check (num_nonnulls(job_id, environment_id, namespace_id) <= 1),
    constraint chk_only_service_id_or_artifact_id_nonnull
        check (num_nonnulls(job_id, artifact_id) <= 1)
);

create unique index metadata_job_map_description_uindex
    on metadata_job_map (description);

create table metadata_service_map
(
    description varchar(100) not null,
    metadata_id integer not null
        constraint metadata_service_map_metadata_id_fk
            references metadata
            on update cascade on delete cascade,
    environment_id integer
        constraint metadata_service_map_environment_id_fk
            references environment
            on update cascade on delete cascade,
    artifact_id integer
        constraint metadata_service_map_artifact_id_fk
            references artifact
            on update cascade on delete cascade,
    namespace_id integer
        constraint metadata_service_map_namespace_id_fk
            references namespace
            on update cascade on delete cascade,
    service_id integer
        constraint metadata_service_map_service_id_fk
            references service
            on update cascade on delete cascade,
    stacking_order integer default 0 not null,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null
    constraint chk_only_at_least_one_must_be_set
        check (num_nonnulls(service_id, environment_id, namespace_id, artifact_id) >= 1),
    constraint chk_only_service_id_or_environment_id_or_namespace_id_nonnull
        check (num_nonnulls(service_id, environment_id, namespace_id) <= 1),
    constraint chk_only_service_id_or_artifact_id_nonnull
        check (num_nonnulls(service_id, artifact_id) <= 1)
);

create unique index metadata_service_map_description_uindex
    on metadata_service_map (description);

create table metadata_history
(
    id serial not null
        constraint metadata_history_pk
            primary key,
    metadata_id integer not null
        constraint metadata_history_metadata_id_fk
            references metadata
            on update cascade on delete cascade,
    description varchar(100) not null,
    value jsonb not null,
    created timestamp,
    created_by varchar(32),
    deleted timestamp,
    deleted_by varchar(32)
);

CREATE OR REPLACE FUNCTION metadata_insert() RETURNS trigger AS
$$
BEGIN
    INSERT INTO metadata_history
    (metadata_id, description, value, created, created_by)
    VALUES
    (NEW.id, NEW.description, NEW.value, current_timestamp, current_user);
    RETURN NEW;
END;
$$
    LANGUAGE plpgsql;

CREATE TRIGGER metadata_insert_trigger
    AFTER INSERT ON metadata
    FOR EACH ROW EXECUTE PROCEDURE metadata_insert();

CREATE OR REPLACE FUNCTION metadata_delete() RETURNS trigger AS
$$
BEGIN
    UPDATE metadata_history
    SET deleted = current_timestamp, deleted_by = current_user
    WHERE deleted IS NULL and metadata_id = OLD.id;
    RETURN NULL;
END;
$$
    LANGUAGE plpgsql;

CREATE TRIGGER metadata_delete_trigger
    AFTER DELETE ON metadata
    FOR EACH ROW EXECUTE PROCEDURE metadata_delete();


CREATE OR REPLACE FUNCTION metadata_update() RETURNS trigger AS
$$
BEGIN

    UPDATE metadata_history
    SET deleted = current_timestamp, deleted_by = current_user
    WHERE deleted IS NULL and metadata_id = OLD.id;

    INSERT INTO metadata_history
    (metadata_id, description, value, created, created_by)
    VALUES
    (NEW.id, NEW.description, NEW.value, current_timestamp, current_user);
    RETURN NEW;

END;
$$
    LANGUAGE plpgsql;

CREATE TRIGGER metadata_update_trigger
    AFTER UPDATE ON metadata
    FOR EACH ROW EXECUTE PROCEDURE metadata_update();

