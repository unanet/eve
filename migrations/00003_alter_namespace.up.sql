DROP INDEX namespace_environment_id_alias_uindex;
CREATE UNIQUE INDEX namespace_environment_id_cluster_id_alias_uindex ON namespace (environment_id, cluster_id, alias);