ALTER TABLE service
    ADD COLUMN service_port integer default 80 NOT NULL,
    ADD COLUMN metrics_port integer default NULL;

ALTER TABLE namespace
    ADD COLUMN name character varying(50);

