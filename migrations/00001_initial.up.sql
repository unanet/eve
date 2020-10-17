DROP TYPE IF EXISTS feed_type;
CREATE TYPE feed_type AS ENUM (
    'docker',
    'generic'
    );


DROP TYPE IF EXISTS provider_group;
CREATE TYPE provider_group AS ENUM (
    'unanet',
    'clearview',
    'ops'
    );


DROP TYPE IF EXISTS deployment_state;
CREATE TYPE deployment_state AS ENUM (
    'queued',
    'scheduled',
    'completed'
    );


DROP TYPE IF EXISTS deployment_cron_state;
CREATE TYPE deployment_cron_state AS ENUM (
    'idle',
    'running'
    );

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION jsonb_merge(orig jsonb, delta jsonb)
RETURNS jsonb LANGUAGE sql AS $$
    SELECT
        jsonb_object_agg(
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


CREATE TABLE feed (
    id integer NOT NULL,
    name character varying(25) NOT NULL,
    promotion_order integer DEFAULT 0 NOT NULL,
    feed_type feed_type,
    alias character varying(25) DEFAULT '' NOT NULL
);
CREATE UNIQUE INDEX feed_name_uindex ON feed USING btree (name);
ALTER TABLE ONLY feed ADD CONSTRAINT feed_pk PRIMARY KEY (id);


CREATE TABLE artifact (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    feed_type feed_type NOT NULL,
    provider_group provider_group NOT NULL,
    function_pointer character varying(250),
    image_tag character varying (25) DEFAULT '$version' NOT NULL,
    service_port integer default 8080 NOT NULL,
    metrics_port integer default 0 NOT NULL,
    service_account character varying (50) DEFAULT 'unanet' NOT NULL,
    run_as character varying (20) DEFAULT '1101' NOT NULL,
    metadata jsonb DEFAULT '{}'::json NOT NULL,
    liveliness_probe jsonb DEFAULT '{}'::json NOT NULL,
    readiness_probe jsonb DEFAULT '{}'::json NOT NULL,
    resource_limits jsonb DEFAULT '{}'::json NOT NULL,
    resource_requests jsonb DEFAULT '{}'::json NOT NULL,
    utilization_limits jsonb DEFAULT '{}'::json NOT NULL,
    autoscaling jsonb DEFAULT '{}'::json NOT NULL,
    pod_resource jsonb DEFAULT '{}'::json NOT NULL
);
CREATE UNIQUE INDEX artifact_name_uindex ON artifact USING btree (name);
ALTER TABLE ONLY artifact ADD CONSTRAINT artifact_pk PRIMARY KEY (id);


CREATE TABLE automation_job (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    parameters jsonb DEFAULT '{}'::json NOT NULL
);
CREATE SEQUENCE automation_job_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE automation_job_id_seq OWNED BY automation_job.id;
ALTER TABLE ONLY automation_job ADD CONSTRAINT automation_job_pk PRIMARY KEY (id);


CREATE TABLE environment (
    id integer NOT NULL,
    name character varying(25) NOT NULL,
    alias character varying(25) NOT NULL,
    metadata jsonb DEFAULT '{}'::json NOT NULL,
    description character varying(1024) NOT NULL
);
CREATE UNIQUE INDEX environment_name_uindex ON environment USING btree (name);
ALTER TABLE ONLY environment ADD CONSTRAINT environment_pk PRIMARY KEY (id);


CREATE TABLE cluster (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    sch_queue_url character varying(200) NOT NULL,
    provider_group provider_group NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);
CREATE UNIQUE INDEX cluster_name_uindex ON cluster USING btree (name);
ALTER TABLE ONLY cluster ADD CONSTRAINT cluster_pk PRIMARY KEY (id);


CREATE TABLE namespace (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    alias character varying(50) NOT NULL,
    environment_id integer NOT NULL,
    requested_version character varying(50) NOT NULL,
    explicit_deploy_only boolean DEFAULT false NOT NULL,
    cluster_id integer NOT NULL,
    metadata jsonb DEFAULT '{}'::json NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);
CREATE SEQUENCE namespace_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE namespace_id_seq OWNED BY namespace.id;
ALTER TABLE ONLY namespace ALTER COLUMN id SET DEFAULT nextval('namespace_id_seq'::regclass);
CREATE UNIQUE INDEX namespace_name_cluster_id_uindex ON namespace (name, cluster_id);
CREATE UNIQUE INDEX namespace_environment_id_cluster_id_alias_uindex ON namespace (environment_id, cluster_id, alias);
ALTER TABLE ONLY namespace ADD CONSTRAINT namespace_pk PRIMARY KEY (id);


CREATE TABLE service (
    id integer NOT NULL,
    name varchar(50) NOT NULL,
    namespace_id integer NOT NULL,
    artifact_id integer NOT NULL,
    override_version character varying(50),
    deployed_version character varying(50),
    metadata jsonb DEFAULT '{}'::json NOT NULL,
    sticky_sessions bool DEFAULT false NOT NULL,
    count int DEFAULT 2 NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    resource_limits jsonb DEFAULT '{}'::json NOT NULL,
    resource_requests jsonb DEFAULT '{}'::json NOT NULL,
    utilization_limits jsonb DEFAULT '{}'::json NOT NULL,
    min_pod int DEFAULT 2 NOT NULL,
    max_pod int DEFAULT 10 NOT NULL,
    autoscaling jsonb DEFAULT '{}'::json NOT NULL,
    pod_resource jsonb DEFAULT '{}'::json NOT NULL

);
CREATE SEQUENCE service_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE service_id_seq OWNED BY service.id;
ALTER TABLE ONLY service ALTER COLUMN id SET DEFAULT nextval('service_id_seq'::regclass);
CREATE UNIQUE INDEX service_namespace_id_name_uindex ON service (name, namespace_id);
ALTER TABLE ONLY service ADD CONSTRAINT service_pk PRIMARY KEY (id);


CREATE TABLE database_server (
    id integer NOT NULL,
    name character varying(50),
    metadata jsonb DEFAULT '{}'::json NOT NULL
);
ALTER TABLE ONLY database_server ADD CONSTRAINT database_server_pk PRIMARY KEY (id);
CREATE UNIQUE INDEX database_server_name_uindex ON database_server USING btree (name);


CREATE TABLE database_type (
    id integer NOT NULL,
    name character varying(50),
    migration_artifact_id integer
);
ALTER TABLE ONLY database_type ADD CONSTRAINT database_type_pk PRIMARY KEY (id);
CREATE UNIQUE INDEX database_type_name_uindex ON database_type USING btree (name);


CREATE TABLE database_instance (
    id integer NOT NULL,
    name character varying(50),
    database_type_id integer NOT NULL,
    database_server_id integer NOT NULL,
    namespace_id integer NOT NULL,
    migration_override_version character varying(50),
    migration_deployed_version character varying(50),
    metadata jsonb DEFAULT '{}'::json NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);
CREATE SEQUENCE database_instance_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE database_instance_id_seq OWNED BY database_instance.id;
ALTER TABLE ONLY database_instance ALTER COLUMN id SET DEFAULT nextval('database_instance_id_seq'::regclass);
ALTER TABLE ONLY database_instance ADD CONSTRAINT database_instance_pk PRIMARY KEY (id);
CREATE UNIQUE INDEX database_instance_namespace_id_name_uindex ON database_instance (name, namespace_id);


CREATE TABLE automation_job_service_map (
    service_id integer NOT NULL,
    automation_job_id integer NOT NULL,
    parameters jsonb DEFAULT '{}'::json NOT NULL
);
CREATE UNIQUE INDEX automation_job_service_map_service_id_automation_job_id_uindex
    ON automation_job_service_map (service_id, automation_job_id);


CREATE TABLE environment_feed_map (
    environment_id int NOT NULL,
    feed_id int NOT NULL
);
CREATE UNIQUE INDEX environment_feed_map_environment_id_feed_id_uindex
    ON environment_feed_map (environment_id, feed_id);


CREATE TABLE deployment (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    environment_id integer NOT NULL,
    namespace_id integer NOT NULL,
    req_id character varying(100),
    message_id character varying(100),
    receipt_handle character varying(1024),
    plan_options jsonb NOT NULL,
    plan_location jsonb,
    state deployment_state NOT NULL,
    "user" character varying(50) NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


CREATE TABLE deployment_cron (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    description character varying(100),
    plan_options jsonb NOT NULL,
    schedule character varying(25) NOT NULL,
    state deployment_cron_state DEFAULT 'idle' NOT NULL,
    last_run timestamp without time zone DEFAULT now() NOT NULL,
    disabled bool DEFAULT false NOT NULL,
    exec_order int DEFAULT 0 NOT NULL
);


CREATE TABLE deployment_cron_job (
    deployment_cron_id UUID NOT NULL,
    deployment_id UUID NOT NULL
);

SELECT pg_catalog.setval('automation_job_id_seq', 1, false);
SELECT pg_catalog.setval('database_instance_id_seq', 1, false);
SELECT pg_catalog.setval('service_id_seq', 1, false);
SELECT pg_catalog.setval('namespace_id_seq', 1, false);

ALTER TABLE ONLY database_type
    ADD CONSTRAINT database_type_migration_artifact_id_fk FOREIGN KEY (migration_artifact_id) REFERENCES artifact(id);

ALTER TABLE ONLY database_instance
    ADD CONSTRAINT database_instance_namespace_id_fk FOREIGN KEY (namespace_id) REFERENCES namespace(id);
ALTER TABLE ONLY database_instance
    ADD CONSTRAINT database_instance_database_type_id_fk FOREIGN KEY (database_type_id) REFERENCES database_type(id);
ALTER TABLE ONLY database_instance
    ADD CONSTRAINT database_instance_database_sever_id_fk FOREIGN KEY (database_server_id) REFERENCES database_server(id);

ALTER TABLE ONLY namespace
    ADD CONSTRAINT namespace_environment_id_fk FOREIGN KEY (environment_id) REFERENCES environment(id);
ALTER TABLE ONLY namespace
    ADD CONSTRAINT namespace_cluster_id_fk FOREIGN KEY (cluster_id) REFERENCES cluster(id);

ALTER TABLE ONLY service
    ADD CONSTRAINT service_artifact_id_fk FOREIGN KEY (artifact_id) REFERENCES artifact(id);
ALTER TABLE ONLY service
    ADD CONSTRAINT service_namespace_id_fk FOREIGN KEY (namespace_id) REFERENCES namespace(id);

ALTER TABLE ONLY automation_job_service_map
    ADD CONSTRAINT automation_job_service_map_automation_job_id_fk FOREIGN KEY (automation_job_id) REFERENCES automation_job(id);
ALTER TABLE ONLY automation_job_service_map
    ADD CONSTRAINT automation_job_service_map_service_id FOREIGN KEY (service_id) REFERENCES service(id);

ALTER TABLE ONLY environment_feed_map
    ADD CONSTRAINT environment_feed_map_environment_id FOREIGN KEY (environment_id) REFERENCES environment(id);
ALTER TABLE ONLY environment_feed_map
    ADD CONSTRAINT environment_feed_map_feed_id FOREIGN KEY(feed_id) REFERENCES feed(id);

ALTER TABLE ONLY deployment
    ADD CONSTRAINT deployment_environment_id FOREIGN KEY (environment_id) REFERENCES environment(id);
ALTER TABLE ONLY deployment
    ADD CONSTRAINT deployment_namespace_id FOREIGN KEY (namespace_id) REFERENCES namespace(id);

ALTER TABLE ONLY deployment_cron_job
    ADD CONSTRAINT deployment_cron_job_deployment_id FOREIGN KEY (deployment_id) REFERENCES deployment(id) ON DELETE CASCADE;;
ALTER TABLE ONLY deployment_cron_job
    ADD CONSTRAINT deployment_cron_job_deployment_cron_id FOREIGN KEY (deployment_cron_id) REFERENCES deployment_cron(id) ON DELETE CASCADE;
