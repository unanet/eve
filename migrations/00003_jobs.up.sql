CREATE TABLE job (
   id integer NOT NULL,
   name character varying(50),
   artifact_id integer NOT NULL,
   namespace_id integer NOT NULL,
   override_version character varying(50),
   deployed_version character varying(50),
   metadata jsonb DEFAULT '{}'::json NOT NULL,
   created_at timestamp without time zone DEFAULT now() NOT NULL,
   updated_at timestamp without time zone DEFAULT now() NOT NULL
);
CREATE SEQUENCE job_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE job_id_seq OWNED BY job.id;
ALTER TABLE ONLY job ALTER COLUMN id SET DEFAULT nextval('job_id_seq'::regclass);
ALTER TABLE ONLY job ADD CONSTRAINT job_pk PRIMARY KEY (id);
CREATE UNIQUE INDEX job_namespace_id_name_uindex ON job (name, namespace_id);

SELECT pg_catalog.setval('job_id_seq', 1, false);

ALTER TABLE ONLY job
    ADD CONSTRAINT job_artifact_id_fk FOREIGN KEY (artifact_id) REFERENCES artifact(id);
ALTER TABLE ONLY job
    ADD CONSTRAINT job_namespace_id_fk FOREIGN KEY (namespace_id) REFERENCES namespace(id);

