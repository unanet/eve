CREATE TABLE IF NOT EXISTS policies (
     p_type character varying(256) DEFAULT ''::character varying NOT NULL,
     v0 character varying(256) DEFAULT ''::character varying NOT NULL,
     v1 character varying(256) DEFAULT ''::character varying NOT NULL,
     v2 character varying(256) DEFAULT ''::character varying NOT NULL,
     v3 character varying(256) DEFAULT ''::character varying NOT NULL,
     v4 character varying(256) DEFAULT ''::character varying NOT NULL,
     v5 character varying(256) DEFAULT ''::character varying NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_policies_p_type ON public.policies USING btree (p_type);
CREATE INDEX IF NOT EXISTS idx_policies_v0 ON public.policies USING btree (v0);
CREATE INDEX IF NOT EXISTS idx_policies_v1 ON public.policies USING btree (v1);
CREATE INDEX IF NOT EXISTS idx_policies_v2 ON public.policies USING btree (v2);
CREATE INDEX IF NOT EXISTS idx_policies_v3 ON public.policies USING btree (v3);
CREATE INDEX IF NOT EXISTS idx_policies_v4 ON public.policies USING btree (v4);
CREATE INDEX IF NOT EXISTS idx_policies_v5 ON public.policies USING btree (v5);

INSERT INTO public.policies (p_type, v0, v1, v2, v3, v4, v5) VALUES ('p', 'user', '/*', 'GET', '', '', '');
INSERT INTO public.policies (p_type, v0, v1, v2, v3, v4, v5) VALUES ('p', 'service', '/*', '*', '', '', '');
INSERT INTO public.policies (p_type, v0, v1, v2, v3, v4, v5) VALUES ('p', 'admin', '/*', '*', '', '', '');