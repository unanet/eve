create table if not exists definition_history
(
    definition_id integer      not null,
    description varchar(100) not null,
    data       jsonb        not null,
    created     timestamp,
    created_by  varchar(32),
    deleted     timestamp,
    deleted_by  varchar(32)
);

create function definition_insert() returns trigger
    language plpgsql
as
$$
BEGIN
    INSERT INTO definition_history
    (definition_id, description, data, created, created_by)
    VALUES (NEW.id, NEW.description, NEW.data, current_timestamp, current_user);
    RETURN NEW;
END;
$$;

create trigger definition_insert_trigger
    after insert
    on definition
    for each row
    execute procedure definition_insert();

create function definition_delete() returns trigger
    language plpgsql
as
$$
BEGIN
    UPDATE definition_history
    SET deleted    = current_timestamp,
        deleted_by = current_user
    WHERE deleted IS NULL
      and definition_id = OLD.id;
    RETURN NULL;
END;
$$;

create trigger definition_delete_trigger
    after delete
    on definition
    for each row
    execute procedure definition_delete();

create function definition_update() returns trigger
    language plpgsql
as
$$
BEGIN
    UPDATE definition_history
    SET deleted    = current_timestamp,
        deleted_by = current_user
    WHERE deleted IS NULL
      and definition_id = OLD.id;

    INSERT INTO definition_history
    (definition_id, description, data, created, created_by)
    VALUES (NEW.id, NEW.description, NEW.data, current_timestamp, current_user);
    RETURN NEW;
END;
$$;

create trigger definition_update_trigger
    after update
    on definition
    for each row
    execute procedure definition_update();
