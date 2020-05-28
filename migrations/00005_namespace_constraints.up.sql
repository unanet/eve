ALTER TABLE namespace ALTER COLUMN name SET NOT NULL;
CREATE UNIQUE INDEX namespace_name_cluster_id_uindex ON namespace (name, cluster_id);