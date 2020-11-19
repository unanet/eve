delete from metadata_history;
delete from metadata;
delete from metadata_service_map;

WITH env_ins AS(
  INSERT INTO metadata (value, description, migrated_from) (SELECT e.metadata, FORMAT('env:%s', e.name), e.id FROM environment e WHERE e.metadata != '{}')
  RETURNING id, description, migrated_from
)

INSERT INTO metadata_service_map(metadata_id, environment_id, description, stacking_order) (SELECT id, migrated_from, description, 100 FROM env_ins);

WITH art_ins AS(
    INSERT INTO metadata (value, description, migrated_from) (SELECT a.metadata, FORMAT('art:%s', a.name), a.id FROM artifact a WHERE a.metadata != '{}')
    RETURNING id, description, migrated_from
)

INSERT INTO metadata_service_map(metadata_id, artifact_id, description, stacking_order) (SELECT id, migrated_from, description, 200 FROM art_ins);

WITH ns_ins AS(
    INSERT INTO metadata (value, description, migrated_from) (SELECT ns.metadata, FORMAT('ns:%s:%s', ns.name, ns.id), ns.id FROM namespace ns WHERE ns.metadata != '{}')
    RETURNING id, description, migrated_from
)

INSERT INTO metadata_service_map(metadata_id, namespace_id, description, stacking_order) (SELECT id, migrated_from, description, 300 FROM ns_ins);

WITH srv_ins AS(
    INSERT INTO metadata (value, description, migrated_from) (SELECT s.metadata, FORMAT('eve-bot:%s:%s', s.name, n.name), s.id FROM service s left join namespace n on s.namespace_id = n.id where s.metadata != '{}')
        RETURNING id, description, migrated_from
)

INSERT INTO metadata_service_map(metadata_id, service_id, description, stacking_order) (SELECT id, migrated_from, description, 400 FROM srv_ins);


