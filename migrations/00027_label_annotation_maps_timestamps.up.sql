ALTER TABLE annotation_job_map
    ADD created_at timestamp default now() not null,
    ADD updated_at timestamp default now() not null;


ALTER TABLE annotation_service_map
    ADD created_at timestamp default now() not null,
    ADD updated_at timestamp default now() not null;


ALTER TABLE label_service_map
    ADD created_at timestamp default now() not null,
    ADD updated_at timestamp default now() not null;


ALTER TABLE label_job_map
    ADD created_at timestamp default now() not null,
    ADD updated_at timestamp default now() not null;