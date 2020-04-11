CREATE TABLE cloud_provider (
    id integer NOT NULL,
    name character varying(25) NOT NULL
);

CREATE SEQUENCE cloud_provider_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE cloud_provider_id_seq OWNED BY cloud_provider.id;

CREATE TABLE cloud_service (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    cloud_provider_id integer NOT NULL,
    description character varying(250)
);

CREATE SEQUENCE cloud_service_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE cloud_service_id_seq OWNED BY cloud_service.id;

CREATE TABLE cloud_service_instance (
    id integer NOT NULL,
    address character varying(50) NOT NULL,
    port integer NOT NULL,
    cloud_service_id integer NOT NULL,
    cloud_service_type_id integer NOT NULL
);

CREATE TABLE cloud_service_instance_context (
    id integer NOT NULL,
    name character varying(50),
    cloud_service_instance_id integer NOT NULL,
    environment_id integer,
    customer_id integer
);

CREATE SEQUENCE cloud_service_instance_context_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE cloud_service_instance_context_id_seq OWNED BY cloud_service_instance_context.id;

CREATE SEQUENCE cloud_service_instance_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE cloud_service_instance_id_seq OWNED BY cloud_service_instance.id;

CREATE TABLE cloud_service_type (
    id integer NOT NULL,
    name character varying NOT NULL,
    default_port integer NOT NULL
);

CREATE SEQUENCE cloud_service_type_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE cloud_service_type_id_seq OWNED BY cloud_service_type.id;

CREATE TABLE customer (
    id integer NOT NULL,
    name character varying(200) NOT NULL,
    subdomain character varying(50) NOT NULL
);

CREATE SEQUENCE customer_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE customer_id_seq OWNED BY customer.id;

CREATE TABLE docker_image (
    id integer NOT NULL,
    name character varying(200) NOT NULL
);

CREATE SEQUENCE docker_image_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE docker_image_id_seq OWNED BY docker_image.id;

CREATE TABLE docker_repository (
    id integer NOT NULL,
    address character varying(250) NOT NULL
);

CREATE SEQUENCE docker_repository_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE docker_repository_id_seq OWNED BY docker_repository.id;

CREATE TABLE docker_repository_type_map (
    docker_type_id integer NOT NULL,
    docker_repository_id integer NOT NULL
);

CREATE TABLE docker_type (
    id integer NOT NULL,
    name character varying(25) NOT NULL,
    cloud_service_restriction_id integer
);

CREATE SEQUENCE docker_type_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE docker_type_id_seq OWNED BY docker_type.id;

CREATE TABLE environment (
    id integer NOT NULL,
    name character varying(25) NOT NULL
);

CREATE SEQUENCE environment_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE environment_id_seq OWNED BY environment.id;

CREATE TABLE k8s_cluster (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    cloud_service_id integer NOT NULL
);

CREATE SEQUENCE k8s_cluster_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE k8s_cluster_id_seq OWNED BY k8s_cluster.id;

CREATE TABLE k8s_namespace (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    environment_id integer NOT NULL,
    default_tag character varying(50),
    k8s_cluster_id integer NOT NULL
);

CREATE TABLE k8s_namespace_customer_map (
    k8s_namespace_id integer NOT NULL,
    customer_id integer NOT NULL
);

CREATE SEQUENCE k8s_namespace_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE k8s_namespace_id_seq OWNED BY k8s_namespace.id;

CREATE TABLE k8s_service (
    id integer NOT NULL,
    k8s_namespace_id integer NOT NULL,
    provisioned_tag character varying(250) NOT NULL,
    docker_image_id integer NOT NULL,
    current_docker_image_digest character varying(80),
    current_tag character varying(50)
);

CREATE SEQUENCE k8s_service_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE k8s_service_id_seq OWNED BY k8s_service.id;

ALTER TABLE ONLY cloud_provider ALTER COLUMN id SET DEFAULT nextval('cloud_provider_id_seq'::regclass);
ALTER TABLE ONLY cloud_service ALTER COLUMN id SET DEFAULT nextval('cloud_service_id_seq'::regclass);
ALTER TABLE ONLY cloud_service_instance ALTER COLUMN id SET DEFAULT nextval('cloud_service_instance_id_seq'::regclass);
ALTER TABLE ONLY cloud_service_instance_context ALTER COLUMN id SET DEFAULT nextval('cloud_service_instance_context_id_seq'::regclass);
ALTER TABLE ONLY cloud_service_type ALTER COLUMN id SET DEFAULT nextval('cloud_service_type_id_seq'::regclass);
ALTER TABLE ONLY customer ALTER COLUMN id SET DEFAULT nextval('customer_id_seq'::regclass);
ALTER TABLE ONLY docker_image ALTER COLUMN id SET DEFAULT nextval('docker_image_id_seq'::regclass);
ALTER TABLE ONLY docker_repository ALTER COLUMN id SET DEFAULT nextval('docker_repository_id_seq'::regclass);
ALTER TABLE ONLY docker_type ALTER COLUMN id SET DEFAULT nextval('docker_type_id_seq'::regclass);
ALTER TABLE ONLY environment ALTER COLUMN id SET DEFAULT nextval('environment_id_seq'::regclass);
ALTER TABLE ONLY k8s_cluster ALTER COLUMN id SET DEFAULT nextval('k8s_cluster_id_seq'::regclass);
ALTER TABLE ONLY k8s_namespace ALTER COLUMN id SET DEFAULT nextval('k8s_namespace_id_seq'::regclass);
ALTER TABLE ONLY k8s_service ALTER COLUMN id SET DEFAULT nextval('k8s_service_id_seq'::regclass);

/* ====================================== SEED DATA FOR DEV PURPOSES ============================================= */

INSERT INTO cloud_provider (id, name) VALUES (1, 'AWS');
INSERT INTO cloud_provider (id, name) VALUES (2, 'Azure');

SELECT pg_catalog.setval('cloud_provider_id_seq', 2, true);

/* ====================================== END SEED DATA ============================================= */

SELECT pg_catalog.setval('cloud_service_id_seq', 1, false);
SELECT pg_catalog.setval('cloud_service_instance_context_id_seq', 1, false);
SELECT pg_catalog.setval('cloud_service_instance_id_seq', 1, false);
SELECT pg_catalog.setval('cloud_service_type_id_seq', 1, false);
SELECT pg_catalog.setval('customer_id_seq', 1, false);
SELECT pg_catalog.setval('docker_image_id_seq', 1, false);
SELECT pg_catalog.setval('docker_repository_id_seq', 1, false);
SELECT pg_catalog.setval('docker_type_id_seq', 1, false);
SELECT pg_catalog.setval('environment_id_seq', 1, false);
SELECT pg_catalog.setval('k8s_cluster_id_seq', 1, false);
SELECT pg_catalog.setval('k8s_namespace_id_seq', 1, false);
SELECT pg_catalog.setval('k8s_service_id_seq', 1, false);

ALTER TABLE ONLY cloud_provider
    ADD CONSTRAINT cloud_provider_pk PRIMARY KEY (id);

ALTER TABLE ONLY cloud_service_instance_context
    ADD CONSTRAINT cloud_service_instance_context_pk PRIMARY KEY (id);

ALTER TABLE ONLY cloud_service_instance
    ADD CONSTRAINT cloud_service_instance_pk PRIMARY KEY (id);

ALTER TABLE ONLY cloud_service
    ADD CONSTRAINT cloud_service_pk PRIMARY KEY (id);

ALTER TABLE ONLY cloud_service_type
    ADD CONSTRAINT cloud_service_type_pk PRIMARY KEY (id);

ALTER TABLE ONLY customer
    ADD CONSTRAINT customer_pk PRIMARY KEY (id);

ALTER TABLE ONLY docker_image
    ADD CONSTRAINT docker_image_pk PRIMARY KEY (id);

ALTER TABLE ONLY docker_repository
    ADD CONSTRAINT docker_repository_pk PRIMARY KEY (id);

ALTER TABLE ONLY docker_type
    ADD CONSTRAINT docker_type_pk PRIMARY KEY (id);

ALTER TABLE ONLY environment
    ADD CONSTRAINT environment_pk PRIMARY KEY (id);

ALTER TABLE ONLY k8s_cluster
    ADD CONSTRAINT k8s_cluster_pk PRIMARY KEY (id);

ALTER TABLE ONLY k8s_namespace
    ADD CONSTRAINT k8s_namespace_pk PRIMARY KEY (id);

ALTER TABLE ONLY k8s_service
    ADD CONSTRAINT k8s_service_pk PRIMARY KEY (id);

CREATE UNIQUE INDEX cloud_provider_name_uindex ON cloud_provider USING btree (name);
CREATE UNIQUE INDEX cloud_service_instance_context_name_uindex ON cloud_service_instance_context USING btree (name);
CREATE UNIQUE INDEX cloud_service_name_uindex ON cloud_service USING btree (name);
CREATE UNIQUE INDEX cloud_service_type_name_uindex ON cloud_service_type USING btree (name);
CREATE UNIQUE INDEX customer_name_uindex ON customer USING btree (name);
CREATE UNIQUE INDEX customer_subdomain_uindex ON customer USING btree (subdomain);
CREATE UNIQUE INDEX docker_image_name_uindex ON docker_image USING btree (name);
CREATE UNIQUE INDEX docker_repository_address_uindex ON docker_repository USING btree (address);
CREATE UNIQUE INDEX docker_type_name_uindex ON docker_type USING btree (name);
CREATE UNIQUE INDEX environment_name_uindex ON environment USING btree (name);
CREATE UNIQUE INDEX k8s_cluster_name_uindex ON k8s_cluster USING btree (name);
CREATE UNIQUE INDEX k8s_namespace_name_uindex ON k8s_namespace USING btree (name);

ALTER TABLE ONLY cloud_service
    ADD CONSTRAINT cloud_service_cloud_provider_id_fk FOREIGN KEY (cloud_provider_id) REFERENCES cloud_provider(id);

ALTER TABLE ONLY cloud_service_instance_context
    ADD CONSTRAINT cloud_service_instance_context_cloud_service_instance_id_fk FOREIGN KEY (cloud_service_instance_id) REFERENCES cloud_service_instance(id);

ALTER TABLE ONLY cloud_service_instance_context
    ADD CONSTRAINT cloud_service_instance_context_customer_id_fk FOREIGN KEY (customer_id) REFERENCES customer(id);

ALTER TABLE ONLY cloud_service_instance_context
    ADD CONSTRAINT cloud_service_instance_context_environment_id_fk FOREIGN KEY (environment_id) REFERENCES environment(id);

ALTER TABLE ONLY k8s_cluster
    ADD CONSTRAINT k8s_cluster_cloud_service_id_fk FOREIGN KEY (cloud_service_id) REFERENCES cloud_service(id);

ALTER TABLE ONLY k8s_namespace
    ADD CONSTRAINT k8s_namespace_environment_id_fk FOREIGN KEY (environment_id) REFERENCES environment(id);

ALTER TABLE ONLY k8s_namespace
    ADD CONSTRAINT k8s_namespace_k8s_cluster_id_fk FOREIGN KEY (k8s_cluster_id) REFERENCES k8s_cluster(id);
