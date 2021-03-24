CREATE TYPE definition_order AS ENUM (
    'main',
    'pre',
    'post'
    );

alter table definition_type
    add class varchar(50),
    add version varchar(50),
    add kind varchar(50),
    add definition_order definition_order;

UPDATE definition_type
SET
    class = 'apps',
    version = 'v1',
    kind = 'Deployment',
    definition_order = 'main'
WHERE name='appsv1.Deployment';


UPDATE definition_type
SET
    class = 'batch',
    version = 'v1',
    kind = 'Job',
    definition_order = 'main'
WHERE name='batchv1.Job';


UPDATE definition_type
SET
    class = 'autoscaling',
    version = 'v2beta2',
    kind = 'HorizontalPodAutoscaler',
    definition_order = 'post'
WHERE name='v2beta2.HorizontalPodAutoscaler';


UPDATE definition_type
SET
    class = '',
    version = 'v1',
    kind = 'Service',
    definition_order = 'main'
WHERE name='apiv1.Service';


ALTER TABLE definition_type
    ALTER COLUMN class SET NOT NULL,
    ALTER COLUMN version SET NOT NULL,
    ALTER COLUMN kind SET NOT NULL,
    ALTER COLUMN definition_order SET NOT NULL;