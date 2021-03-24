INSERT INTO definition(description,definition_type_id,data)
SELECT ('defaultDeployment') as description,
       (SELECT id from definition_type where kind = 'Deployment') as definition_type_id,
       ('{}') as data;

INSERT INTO definition(description,definition_type_id,data)
SELECT ('defaultService') as description,
       (SELECT id from definition_type where kind = 'Service') as definition_type_id,
       ('{}') as data;

INSERT INTO definition(description,definition_type_id,data)
SELECT ('defaultJob') as description,
       (SELECT id from definition_type where kind = 'Job') as definition_type_id,
       ('{}') as data;
