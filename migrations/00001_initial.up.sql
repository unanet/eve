DROP TYPE IF EXISTS feed_type;
CREATE TYPE feed_type AS ENUM (
    'docker',
    'generic'
    );

DROP TYPE IF EXISTS provider_group;
CREATE TYPE provider_group AS ENUM (
    'unanet',
    'clearview'
    );

DROP TYPE IF EXISTS deployment_state;
CREATE TYPE deployment_state AS ENUM (
    'queued',
    'scheduled',
    'completed'
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
    feed_type feed_type
);
CREATE UNIQUE INDEX feed_name_uindex ON feed USING btree (name);
ALTER TABLE ONLY feed ADD CONSTRAINT feed_pk PRIMARY KEY (id);


CREATE TABLE artifact (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    feed_type feed_type NOT NULL,
    provider_group provider_group NOT NULL,
    function_pointer character varying(250),
    metadata jsonb DEFAULT '{}'::json NOT NULL
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
    metadata jsonb DEFAULT '{}'::json NOT NULL
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
    domain character varying(200) NOT NULL,
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
CREATE UNIQUE INDEX namespace_name_uindex ON namespace USING btree (name);
CREATE UNIQUE INDEX namespace_domain_uindex ON namespace USING btree (domain);
CREATE UNIQUE INDEX namespace_environment_id_alias_uindex ON namespace (environment_id, alias);
ALTER TABLE ONLY namespace ADD CONSTRAINT namespace_pk PRIMARY KEY (id);


CREATE TABLE service (
    id integer NOT NULL,
    namespace_id integer NOT NULL,
    artifact_id integer NOT NULL,
    override_version character varying(50),
    deployed_version character varying(50),
    metadata jsonb DEFAULT '{}'::json NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
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
CREATE UNIQUE INDEX database_instance_name_uindex ON database_instance USING btree (name);


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

/* ====================================== SEED DATA ============================================= */

INSERT INTO feed (id, name, promotion_order, feed_type) VALUES (1, 'docker-int', 0, 'docker');
INSERT INTO feed (id, name, promotion_order, feed_type) VALUES (2, 'generic-int', 0, 'generic');
INSERT INTO feed (id, name, promotion_order, feed_type) VALUES (3, 'docker-qa', 1, 'docker');
INSERT INTO feed (id, name, promotion_order, feed_type) VALUES (4, 'generic-qa', 1, 'generic');
INSERT INTO feed (id, name, promotion_order, feed_type) VALUES (5, 'docker-stage', 2, 'docker');
INSERT INTO feed (id, name, promotion_order, feed_type) VALUES (6, 'generic-stage', 2, 'generic');
INSERT INTO feed (id, name, promotion_order, feed_type) VALUES (7, 'docker-prod', 3, 'docker');
INSERT INTO feed (id, name, promotion_order, feed_type) VALUES (8, 'generic-prod', 3, 'generic');

INSERT INTO environment (id, name) VALUES(1, 'int');
INSERT INTO environment (id, name) VALUES(2, 'qa');
INSERT INTO environment (id, name) VALUES(3, 'demo');
INSERT INTO environment (id, name) VALUES(4, 'stage');
INSERT INTO environment (id, name) VALUES(5, 'prod');

INSERT INTO environment_feed_map(environment_id, feed_id) VALUES (1, 1);
INSERT INTO environment_feed_map(environment_id, feed_id) VALUES (1, 2);
INSERT INTO environment_feed_map(environment_id, feed_id) VALUES (2, 3);
INSERT INTO environment_feed_map(environment_id, feed_id) VALUES (2, 4);
INSERT INTO environment_feed_map(environment_id, feed_id) VALUES (3, 3);
INSERT INTO environment_feed_map(environment_id, feed_id) VALUES (3, 4);
INSERT INTO environment_feed_map(environment_id, feed_id) VALUES (4, 5);
INSERT INTO environment_feed_map(environment_id, feed_id) VALUES (4, 6);
INSERT INTO environment_feed_map(environment_id, feed_id) VALUES (4, 7);
INSERT INTO environment_feed_map(environment_id, feed_id) VALUES (4, 8);

/* ================== CLEARVIEW APPS ================== */
INSERT INTO artifact(id, name, feed_type, provider_group, function_pointer, metadata) VALUES (101, 'support', 'generic', 'clearview', 'https://cv-cloud-ops.azurewebsites.net/api/sites/support/create',
    '{"inject_vault_paths":"{{ .Plan.Namespace.ClusterName }}", "environment": "{{ .Plan.EnvironmentName }}", "namespace": "{{ .Plan.Namespace.Alias }}", "cluster": "{{ .Plan.Namespace.ClusterName }}", "artifact_name": "{{ .Service.ArtifactName }}", "artifact_version": "{{ .Service.AvailableVersion }}", "artifact_repo":"{{ .Service.ArtifactoryFeed }}", "artifact_path": "{{ .Service.ArtifactoryPath }}" }');
INSERT INTO artifact(id, name, feed_type, provider_group, function_pointer, metadata) VALUES (105, 'infocus-reports', 'generic', 'clearview', 'https://cv-cloud-ops.azurewebsites.net/api/sites/reports/create',
    '{"inject_vault_paths":"{{ .Plan.Namespace.ClusterName }}", "environment": "{{ .Plan.EnvironmentName }}", "namespace": "{{ .Plan.Namespace.Alias }}", "cluster": "{{ .Plan.Namespace.ClusterName }}", "artifact_name": "{{ .Service.ArtifactName }}", "artifact_version": "{{ .Service.AvailableVersion }}", "artifact_repo":"{{ .Service.ArtifactoryFeed }}", "artifact_path": "{{ .Service.ArtifactoryPath }}" }');
INSERT INTO artifact(id, name, feed_type, provider_group, function_pointer, metadata) VALUES (106, 'infocus-windows', 'generic', 'clearview', 'https://cv-windows-client.azurewebsites.net/api/setup/client',
    '{"inject_vault_paths":"{{ .Plan.Namespace.ClusterName }}", "environment": "{{ .Plan.EnvironmentName }}", "namespace": "{{ .Plan.Namespace.Alias }}", "cluster": "{{ .Plan.Namespace.ClusterName }}", "artifact_name": "{{ .Service.ArtifactName }}", "artifact_version": "{{ .Service.AvailableVersion }}", "artifact_repo":"{{ .Service.ArtifactoryFeed }}", "artifact_path": "{{ .Service.ArtifactoryPath }}" }');
INSERT INTO artifact(id, name, feed_type, provider_group) VALUES (120, 'infocus-cloud-client', 'docker', 'clearview');
INSERT INTO artifact(id, name, feed_type, provider_group) VALUES (121, 'infocus-documents', 'docker', 'clearview');
INSERT INTO artifact(id, name, feed_type, provider_group) VALUES (122, 'infocus-proxy', 'docker', 'clearview');
INSERT INTO artifact(id, name, feed_type, provider_group, metadata) VALUES (123, 'infocus-web', 'docker', 'clearview',
    '{"inject_vault_paths":"{{ .Plan.Namespace.ClusterName }}", "cloud_db_name": "cvs_{{ .Plan.EnvironmentName }}_cloud", "support_db_name": "cvs_{{ .Plan.EnvironmentName }}_support" }');

INSERT INTO artifact(id, name, feed_type, provider_group) VALUES (185, 'cvs-migrations', 'docker', 'clearview');

/* ================== UNANET APPS ================== */
INSERT INTO artifact(id, name, feed_type, provider_group) VALUES (201, 'unanet', 'docker', 'unanet');
INSERT INTO artifact(id, name, feed_type, provider_group) VALUES (202, 'platform', 'docker', 'unanet');
INSERT INTO artifact(id, name, feed_type, provider_group) VALUES (203, 'exago', 'docker', 'unanet');

INSERT INTO artifact(id, name, feed_type, provider_group) VALUES (285, 'sql-migration-scripts', 'docker', 'unanet');

INSERT INTO cluster(id, name, provider_group, sch_queue_url) VALUES (1, 'cvs-nonprod-zxrjdqr67u', 'clearview', 'https://sqs.us-east-2.amazonaws.com/580107804399/cvs-nonprod-zxrjdqr67u.fifo');

INSERT INTO namespace(id, name, alias, environment_id, requested_version, cluster_id, domain) VALUES (1, 'cvs-int', 'cvs', 1, '2020', 1, 'int.infocus.app');

INSERT INTO namespace(id, name, alias, environment_id, requested_version, cluster_id, domain) VALUES (2, 'cvs-prev-int', 'cvs-prev', 1, '2020.2', 1, 'prev-int.infocus.app');
INSERT INTO namespace(id, name, alias, environment_id, requested_version, cluster_id, domain) VALUES (3, 'cvs-curr-int', 'cvs-curr', 1, '2020.2', 1, 'curr-int.infocus.app');
INSERT INTO namespace(id, name, alias, environment_id, requested_version, cluster_id, domain) VALUES (4, 'cvs-next-int', 'cvs-next', 1, '2020', 1, 'next-int.infocus.app');

INSERT INTO namespace(id, name, alias, environment_id, requested_version, cluster_id, domain) VALUES (5, 'cvs-qa', 'cvs', 2, '2020', 1, 'qa.infocus.app');

INSERT INTO namespace(id, name, alias, environment_id, requested_version, cluster_id, domain) VALUES (6, 'cvs-prev-qa', 'cvs-prev', 2, '2020.2', 1, 'prev-qa.infocus.app');
INSERT INTO namespace(id, name, alias, environment_id, requested_version, cluster_id, domain) VALUES (7, 'cvs-curr-qa', 'cvs-curr', 2, '2020.2', 1, 'curr.qa-infocus.app');
INSERT INTO namespace(id, name, alias, environment_id, requested_version, cluster_id, domain) VALUES (8, 'cvs-next-qa', 'cvs-next', 2, '2020', 1, 'next-qa.infocus.app');

SELECT pg_catalog.setval('namespace_id_seq', 8, true);

INSERT INTO service(id, namespace_id, artifact_id) VALUES (1, 1, 101);

INSERT INTO service(id, namespace_id, artifact_id) VALUES (2, 2, 105);
INSERT INTO service(id, namespace_id, artifact_id) VALUES (3, 2, 106);
INSERT INTO service(id, namespace_id, artifact_id) VALUES (4, 2, 120);
INSERT INTO service(id, namespace_id, artifact_id) VALUES (5, 2, 121);
INSERT INTO service(id, namespace_id, artifact_id) VALUES (6, 2, 122);
INSERT INTO service(id, namespace_id, artifact_id) VALUES (7, 2, 123);

INSERT INTO service(id, namespace_id, artifact_id) VALUES (8, 3, 105);
INSERT INTO service(id, namespace_id, artifact_id) VALUES (9, 3, 106);
INSERT INTO service(id, namespace_id, artifact_id) VALUES (10, 3, 120);
INSERT INTO service(id, namespace_id, artifact_id) VALUES (11, 3, 121);
INSERT INTO service(id, namespace_id, artifact_id) VALUES (12, 3, 122);
INSERT INTO service(id, namespace_id, artifact_id) VALUES (13, 3, 123);

INSERT INTO service(id, namespace_id, artifact_id) VALUES (14, 4, 105);
INSERT INTO service(id, namespace_id, artifact_id) VALUES (15, 4, 106);
INSERT INTO service(id, namespace_id, artifact_id) VALUES (16, 4, 120);
INSERT INTO service(id, namespace_id, artifact_id) VALUES (17, 4, 121);
INSERT INTO service(id, namespace_id, artifact_id) VALUES (18, 4, 122);
INSERT INTO service(id, namespace_id, artifact_id) VALUES (19, 4, 123);

INSERT INTO service(id, namespace_id, artifact_id) VALUES (20, 5, 101);

INSERT INTO service(id, namespace_id, artifact_id, override_version, deployed_version) VALUES (21, 6, 105, NULL, NULL);
INSERT INTO service(id, namespace_id, artifact_id, override_version, deployed_version) VALUES (22, 6, 106, NULL, NULL);
INSERT INTO service(id, namespace_id, artifact_id, override_version, deployed_version) VALUES (23, 6, 120, NULL, NULL);
INSERT INTO service(id, namespace_id, artifact_id, override_version, deployed_version) VALUES (24, 6, 121, NULL, NULL);
INSERT INTO service(id, namespace_id, artifact_id, override_version, deployed_version) VALUES (25, 6, 122, NULL, NULL);
INSERT INTO service(id, namespace_id, artifact_id, override_version, deployed_version) VALUES (26, 6, 123, NULL, NULL);

INSERT INTO service(id, namespace_id, artifact_id, override_version, deployed_version) VALUES (27, 7, 105, NULL, NULL);
INSERT INTO service(id, namespace_id, artifact_id, override_version, deployed_version) VALUES (28, 7, 106, NULL, NULL);
INSERT INTO service(id, namespace_id, artifact_id, override_version, deployed_version) VALUES (29, 7, 120, NULL, NULL);
INSERT INTO service(id, namespace_id, artifact_id, override_version, deployed_version) VALUES (30, 7, 121, NULL, NULL);
INSERT INTO service(id, namespace_id, artifact_id, override_version, deployed_version) VALUES (31, 7, 122, NULL, NULL);
INSERT INTO service(id, namespace_id, artifact_id, override_version, deployed_version) VALUES (32, 7, 123, NULL, NULL);

INSERT INTO service(id, namespace_id, artifact_id, override_version, deployed_version) VALUES (33, 8, 105, NULL, NULL);
INSERT INTO service(id, namespace_id, artifact_id, override_version, deployed_version) VALUES (34, 8, 106, NULL, NULL);
INSERT INTO service(id, namespace_id, artifact_id, override_version, deployed_version) VALUES (35, 8, 120, NULL, NULL);
INSERT INTO service(id, namespace_id, artifact_id, override_version, deployed_version) VALUES (36, 8, 121, NULL, NULL);
INSERT INTO service(id, namespace_id, artifact_id, override_version, deployed_version) VALUES (37, 8, 122, NULL, NULL);
INSERT INTO service(id, namespace_id, artifact_id, override_version, deployed_version) VALUES (38, 8, 123, NULL, NULL);

SELECT pg_catalog.setval('service_id_seq', 38, true);

INSERT INTO database_type(id, name, migration_artifact_id) VALUES (101, 'cvs-cloud', 185);
INSERT INTO database_type(id, name, migration_artifact_id) VALUES (102, 'cvs-support', NULL);

INSERT INTO database_server(id, name) VALUES (101, 'cvs-nonprod-zxrjdqr67u');

INSERT INTO database_instance(id, name, database_type_id, database_server_id, namespace_id, migration_override_version, migration_deployed_version)
    VALUES(1, 'cvs-int-cloud', 101, 101, 1, NULL, NULL);
INSERT INTO database_instance(id, name, database_type_id, database_server_id, namespace_id, migration_override_version, migration_deployed_version)
    VALUES(2, 'cvs-int-support', 102, 101, 1, NULL, NULL);

SELECT pg_catalog.setval('database_instance_id_seq', 2, true);

/* ====================================== END SEED DATA ============================================= */

SELECT pg_catalog.setval('automation_job_id_seq', 1, false);

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
