delete from metadata;
delete from metadata_service_map;

WITH env_ins AS(
  INSERT INTO metadata (value, description) (SELECT metadata, FORMAT('env-%s', name) FROM environment e WHERE e.metadata != '{}')
  RETURNING id, description
)

INSERT INTO metadata_service_map(metadata_id, environment_id, description, stacking_order) (SELECT id, (SELECT id FROM environment WHERE name = ltrim(env_ins.description, 'env-')), description, 100 FROM env_ins);

WITH ns_ins AS(
    INSERT INTO metadata (value, description) (SELECT metadata, FORMAT('ns-%s-%s', name, id) FROM namespace ns WHERE ns.metadata != '{}')
    RETURNING id, description
)

INSERT INTO metadata_service_map(metadata_id, namespace_id, description, stacking_order) (SELECT id, (SELECT id FROM environment WHERE name = ltrim(env_ins.description, 'env-')), description, 100 FROM env_ins);


