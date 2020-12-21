INSERT INTO annotation_service_map(description,annotation_id,artifact_id)
SELECT ('unanet path') as description,
       (SELECT id from annotation where description = 'art:unanet') as annotation_id,
       (SELECT id from artifact where name = 'unanet-app') as artifact_id;


INSERT INTO annotation_service_map(description,annotation_id,artifact_id)
SELECT ('unanet analytics') as description,
       (SELECT id from annotation where description = 'art:unanet-analytics') as annotation_id,
       (SELECT id from artifact where name = 'unanet-analytics') as artifact_id;


INSERT INTO annotation_service_map(description,annotation_id,artifact_id)
SELECT ('unanet applinks') as description,
       (SELECT id from annotation where description = 'art:applinks') as annotation_id,
       (SELECT id from artifact where name = 'applinks') as artifact_id;


INSERT INTO label_service_map(description,label_id,environment_id)
SELECT ('env:una-int') as description,
       (SELECT id from label where description = 'unanet-proxy') as label_id,
       (SELECT id from environment where name = 'una-int') as environment_id;


INSERT INTO label_service_map(description,label_id,environment_id)
SELECT ('env:una-qa') as description,
       (SELECT id from label where description = 'unanet-proxy') as label_id,
       (SELECT id from environment where name = 'una-qa') as environment_id;


INSERT INTO label_service_map(description,label_id,environment_id)
SELECT ('env:una-stage') as description,
       (SELECT id from label where description = 'unanet-proxy') as label_id,
       (SELECT id from environment where name = 'una-stage') as environment_id;


INSERT INTO label_service_map(description,label_id,environment_id)
SELECT ('env:una-stage') as description,
       (SELECT id from label where description = 'unanet-proxy') as label_id,
       (SELECT id from environment where name = 'una-stage') as environment_id;


INSERT INTO label_service_map(description,label_id,environment_id)
SELECT ('env:una-prod') as description,
       (SELECT id from label where description = 'unanet-proxy') as label_id,
       (SELECT id from environment where name = 'una-prod') as environment_id;


INSERT INTO label_service_map(description,label_id,environment_id)
SELECT ('env:una-demo') as description,
       (SELECT id from label where description = 'unanet-proxy') as label_id,
       (SELECT id from environment where name = 'una-demo') as environment_id;


INSERT INTO label_service_map(description,label_id,environment_id)
SELECT ('env:una-portal') as description,
       (SELECT id from label where description = 'unanet-proxy') as label_id,
       (SELECT id from environment where name = 'una-portal') as environment_id;
