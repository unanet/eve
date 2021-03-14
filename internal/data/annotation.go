package data

import (
	"context"
	"database/sql"
	goErrors "errors"
	"time"

	"gitlab.unanet.io/devops/go/pkg/errors"
	"gitlab.unanet.io/devops/go/pkg/json"
)

type Annotation struct {
	ID          int          `db:"id"`
	Description string       `db:"description"`
	Data        json.Object  `db:"data"`
	CreatedAt   sql.NullTime `db:"created_at"`
	UpdatedAt   sql.NullTime `db:"updated_at"`
}

type AnnotationServiceMap struct {
	Description   string        `db:"description"`
	AnnotationID  int           `db:"annotation_id"`
	EnvironmentID sql.NullInt32 `db:"environment_id"`
	ArtifactID    sql.NullInt32 `db:"artifact_id"`
	NamespaceID   sql.NullInt32 `db:"namespace_id"`
	ServiceID     sql.NullInt32 `db:"service_id"`
	StackingOrder int           `db:"stacking_order"`
	CreatedAt     sql.NullTime  `db:"created_at"`
	UpdatedAt     sql.NullTime  `db:"updated_at"`
}

type AnnotationJobMap struct {
	Description   string        `db:"description"`
	AnnotationID  int           `db:"annotation_id"`
	EnvironmentID sql.NullInt32 `db:"environment_id"`
	ArtifactID    sql.NullInt32 `db:"artifact_id"`
	NamespaceID   sql.NullInt32 `db:"namespace_id"`
	JobID         sql.NullInt32 `db:"job_id"`
	StackingOrder int           `db:"stacking_order"`
	CreatedAt     sql.NullTime  `db:"created_at"`
	UpdatedAt     sql.NullTime  `db:"updated_at"`
}

type AnnotationService struct {
	AnnotationID          int           `db:"annotation_id"`
	Data                  json.Object   `db:"data"`
	AnnotationDescription string        `db:"annotation_description"`
	MapDescription        string        `db:"map_description"`
	MapEnvironmentID      sql.NullInt32 `db:"map_environment_id"`
	MapArtifactID         sql.NullInt32 `db:"map_artifact_id"`
	MapNamespaceID        sql.NullInt32 `db:"map_namespace_id"`
	MapServiceID          sql.NullInt32 `db:"map_service_id"`
	StackingOrder         int           `db:"stacking_order"`
	CreatedAt             sql.NullTime  `db:"created_at"`
	UpdatedAt             sql.NullTime  `db:"updated_at"`
}

type AnnotationJob struct {
	AnnotationID          int           `db:"annotation_id"`
	Data                  json.Object   `db:"data"`
	AnnotationDescription string        `db:"annotation_description"`
	MapDescription        string        `db:"map_description"`
	MapEnvironmentID      sql.NullInt32 `db:"map_environment_id"`
	MapArtifactID         sql.NullInt32 `db:"map_artifact_id"`
	MapNamespaceID        sql.NullInt32 `db:"map_namespace_id"`
	MapJobID              sql.NullInt32 `db:"map_job_id"`
	StackingOrder         int           `db:"stacking_order"`
	CreatedAt             sql.NullTime  `db:"created_at"`
	UpdatedAt             sql.NullTime  `db:"updated_at"`
}

func (r *Repo) UpsertMergeAnnotation(ctx context.Context, l *Annotation) error {
	now := time.Now().UTC()
	l.CreatedAt = sql.NullTime{
		Time:  now,
		Valid: true,
	}
	l.UpdatedAt = sql.NullTime{
		Time:  now,
		Valid: true,
	}

	err := r.db.QueryRowxContext(ctx, `
	INSERT INTO annotation(description, data, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (description)
		DO UPDATE SET data = annotation.data || $2, updated_at = $4
		RETURNING id, value, created_at
	`, l.Description, l.Data, l.CreatedAt, l.UpdatedAt).
		StructScan(l)

	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (r *Repo) UpsertAnnotation(ctx context.Context, l *Annotation) error {
	now := time.Now().UTC()
	l.CreatedAt = sql.NullTime{
		Time:  now,
		Valid: true,
	}
	l.UpdatedAt = sql.NullTime{
		Time:  now,
		Valid: true,
	}

	err := r.db.QueryRowxContext(ctx, `
	
	INSERT INTO annotation(description, data, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (description)
		DO UPDATE SET data = $2, updated_at = $4
		RETURNING id, created_at
	
	`, l.Description, l.Data, l.CreatedAt, l.UpdatedAt).
		StructScan(l)

	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (r *Repo) UpsertAnnotationJobMap(ctx context.Context, ljm *AnnotationJobMap) error {
	now := time.Now().UTC()
	ljm.CreatedAt = sql.NullTime{
		Time:  now,
		Valid: true,
	}
	ljm.UpdatedAt = sql.NullTime{
		Time:  now,
		Valid: true,
	}

	err := r.db.QueryRowxContext(ctx, `
	
	INSERT INTO annotation_job_map(description, annotation_id, environment_id, artifact_id, namespace_id, job_id, stacking_order, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	ON CONFLICT (description)
	DO UPDATE SET environment_id = $3, artifact_id = $4, namespace_id = $5, job_id = $6, stacking_order = $7, updated_at = $9
	RETURNING created_at
	
	`, ljm.Description, ljm.AnnotationID, ljm.EnvironmentID, ljm.ArtifactID, ljm.NamespaceID, ljm.JobID, ljm.StackingOrder, ljm.CreatedAt, ljm.UpdatedAt).
		StructScan(ljm)

	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (r *Repo) UpsertAnnotationServiceMap(ctx context.Context, lsm *AnnotationServiceMap) error {
	now := time.Now().UTC()
	lsm.CreatedAt = sql.NullTime{
		Time:  now,
		Valid: true,
	}
	lsm.UpdatedAt = sql.NullTime{
		Time:  now,
		Valid: true,
	}

	err := r.db.QueryRowxContext(ctx, `
	
	INSERT INTO annotation_service_map(description, annotation_id, environment_id, artifact_id, namespace_id, service_id, stacking_order, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	ON CONFLICT (description)
	DO UPDATE SET environment_id = $3, artifact_id = $4, namespace_id = $5, service_id = $6, stacking_order = $7, updated_at = $9
	RETURNING created_at
	
	`, lsm.Description, lsm.AnnotationID, lsm.EnvironmentID, lsm.ArtifactID, lsm.NamespaceID, lsm.ServiceID, lsm.StackingOrder, lsm.CreatedAt, lsm.UpdatedAt).
		StructScan(lsm)

	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (r *Repo) GetAnnotation(ctx context.Context, annotationID int) (*Annotation, error) {
	var annotation Annotation

	row := r.db.QueryRowxContext(ctx, `
		select id, 
		       description, 
		       data, 
		       created_at, 
		       updated_at
		from annotation
		where id = $1
		`, annotationID)
	err := row.StructScan(&annotation)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return nil, NotFoundErrorf("annotation with id: %d not found", annotationID)
		}
		return nil, errors.Wrap(err)
	}

	return &annotation, nil
}

func (r *Repo) GetAnnotationByDescription(ctx context.Context, description string) (*Annotation, error) {
	var annotation Annotation

	row := r.db.QueryRowxContext(ctx, `
		select id, 
		       description, 
		       data, 
		       created_at, 
		       updated_at
		from annotation
		where description = $1
		`, description)
	err := row.StructScan(&annotation)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return nil, NotFoundErrorf("annotation with description: %s not found", description)
		}
		return nil, errors.Wrap(err)
	}

	return &annotation, nil
}

func (r *Repo) Annotations(ctx context.Context) ([]Annotation, error) {
	rows, err := r.db.QueryxContext(ctx, `
		select id, 
		       description, 
		       data, 
		       created_at, 
		       updated_at 
		from annotation
	`)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var annotations []Annotation
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}

		var annotation Annotation
		err = rows.StructScan(&annotation)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		annotations = append(annotations, annotation)
	}

	return annotations, nil
}

func (r *Repo) JobAnnotationMaps(ctx context.Context, jobID int) ([]AnnotationJobMap, error) {
	rows, err := r.db.QueryxContext(ctx, `
		select description, 
		       annotation_id, 
		       environment_id, 
		       artifact_id, 
		       namespace_id, 
		       job_id, 
		       stacking_order, 
		       created_at, 
		       updated_at
		from annotation_job_map
		where job_id = $1
		`, jobID)

	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var annotationJobMaps []AnnotationJobMap
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}
		var annotationJobMap AnnotationJobMap
		err = rows.StructScan(&annotationJobMap)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		annotationJobMaps = append(annotationJobMaps, annotationJobMap)
	}

	return annotationJobMaps, nil
}

func (r *Repo) JobAnnotationMapsByAnnotationID(ctx context.Context, annotationID int) ([]AnnotationJobMap, error) {
	rows, err := r.db.QueryxContext(ctx, `
		select description, 
		       annotation_id, 
		       environment_id, 
		       artifact_id, 
		       namespace_id, 
		       job_id, 
		       stacking_order, 
		       created_at, 
		       updated_at
		from annotation_job_map
		where annotation_id = $1
		`, annotationID)

	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var annotationJobMaps []AnnotationJobMap
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}
		var annotationJobMap AnnotationJobMap
		err = rows.StructScan(&annotationJobMap)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		annotationJobMaps = append(annotationJobMaps, annotationJobMap)
	}

	return annotationJobMaps, nil
}

func (r *Repo) ServiceAnnotationMapsByAnnotationID(ctx context.Context, annotationID int) ([]AnnotationServiceMap, error) {
	rows, err := r.db.QueryxContext(ctx, `
		select description, 
		       annotation_id, 
		       environment_id, 
		       artifact_id, 
		       namespace_id, 
		       service_id, 
		       stacking_order, 
		       created_at, 
		       updated_at
		from annotation_service_map
		where annotation_id = $1
		`, annotationID)

	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var annotationServiceMaps []AnnotationServiceMap
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}
		var annotationServiceMap AnnotationServiceMap
		err = rows.StructScan(&annotationServiceMap)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		annotationServiceMaps = append(annotationServiceMaps, annotationServiceMap)
	}

	return annotationServiceMaps, nil
}

func (r *Repo) JobAnnotations(ctx context.Context, jobID int) ([]AnnotationJob, error) {
	rows, err := r.db.QueryxContext(ctx, `
		WITH env_data AS (
			select j.id as job_id, 
			       environment_id, 
			       namespace_id, 
			       artifact_id,
			       cluster_id
			from job j 
			    left join namespace n on j.namespace_id = n.id 
			    left join environment e on n.environment_id = e.id
			where j.id = $1
		)
		
		SELECT l.id as annotation_id,
		       l.data as data,
		       l.description as annotation_description,
		       ljm.description as map_description,
		       ljm.environment_id as map_environment_id,
		       ljm.artifact_id as map_artifact_id,
		       ljm.namespace_id as map_namespace_id,
		       ljm.job_id as map_job_id,
		       ljm.stacking_order as stacking_order,
		       l.created_at,
		       l.updated_at
		FROM annotation_job_map ljm 
		    LEFT JOIN annotation l ON ljm.annotation_id = l.id 
			LEFT JOIN env_data ed on ljm.job_id = $1
		WHERE
			(ljm.job_id = $1)
		OR
			(ljm.cluster_id = ed.cluster_id AND ljm.artifact_id IS NULL)
		OR
		    (ljm.environment_id = ed.environment_id AND ljm.artifact_id IS NULL) 
		OR
		    (ljm.namespace_id = ed.namespace_id AND ljm.artifact_id IS NULL)
		OR
		    (ljm.artifact_id = ed.artifact_id AND ljm.environment_id IS NULL AND ljm.namespace_id IS NULL AND ljm.cluster_id IS NULL)
		OR
		    (ljm.artifact_id = ed.artifact_id AND ljm.cluster_id = ed.cluster_id)
		OR
		    (ljm.artifact_id = ed.artifact_id AND ljm.environment_id = ed.environment_id)
		OR
		    (ljm.artifact_id = ed.artifact_id AND ljm.namespace_id = ed.namespace_id)
		ORDER BY
			ljm.stacking_order
	`, jobID)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var annotationJobs []AnnotationJob
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}
		var annotationJob AnnotationJob
		err = rows.StructScan(&annotationJob)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		annotationJobs = append(annotationJobs, annotationJob)
	}

	return annotationJobs, nil
}

func (r *Repo) ServiceAnnotations(ctx context.Context, serviceID int) ([]AnnotationService, error) {
	rows, err := r.db.QueryxContext(ctx, `
		WITH env_data AS (
			select s.id as service_id,
				   environment_id,
				   namespace_id,
				   artifact_id,
				   cluster_id
			from service s
					 left join namespace n on s.namespace_id = n.id
					 left join environment e on n.environment_id = e.id
			where s.id = $1
		)
		
		SELECT a.id as annotation_id,
			   a.data as data,
			   a.description as annotation_description,
			   asm.description as map_description,
			   asm.environment_id as map_environment_id,
			   asm.artifact_id as map_artifact_id,
			   asm.namespace_id as map_namespace_id,
			   asm.service_id as map_service_id,
			   asm.stacking_order as stacking_order,
			   a.created_at,
			   a.updated_at
		FROM annotation_service_map asm
				 LEFT JOIN annotation a ON asm.annotation_id = a.id
				 LEFT JOIN env_data ed on ed.service_id = $1
		WHERE
			(asm.service_id = $1)
		   OR
			(asm.cluster_id = ed.cluster_id AND asm.artifact_id IS NULL)
		   OR
			(asm.environment_id = ed.environment_id AND asm.artifact_id IS NULL)
		   OR
			(asm.namespace_id = ed.namespace_id AND asm.artifact_id IS NULL)
		   OR
			(asm.artifact_id = ed.artifact_id AND asm.environment_id IS NULL AND asm.namespace_id IS NULL AND asm.cluster_id IS NULL)
		   OR
			(asm.artifact_id = ed.artifact_id AND asm.cluster_id = ed.cluster_id)
		   OR
			(asm.artifact_id = ed.artifact_id AND asm.environment_id = ed.environment_id)
		   OR
			(asm.artifact_id = ed.artifact_id AND asm.namespace_id = ed.namespace_id)
		ORDER BY
			asm.stacking_order
`, serviceID)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var annotationServices []AnnotationService
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}
		var annotationService AnnotationService
		err = rows.StructScan(&annotationService)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		annotationServices = append(annotationServices, annotationService)
	}

	return annotationServices, nil
}

func (r *Repo) DeleteAnnotation(ctx context.Context, annotationID int) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM annotation WHERE id = $1
	`, annotationID)
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.NotFoundf("annotation id: %d not found", annotationID)
	}

	return nil
}

func (r *Repo) DeleteAnnotationKey(ctx context.Context, annotationID int, key string) (*Annotation, error) {
	var annotation Annotation
	err := r.db.QueryRowxContext(ctx, `
		UPDATE annotation SET data = annotation.data - $1 WHERE id = $2
		RETURNING id, value, description, created_at, updated_at
	`, key, annotationID).StructScan(&annotation)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return &annotation, nil
}

func (r *Repo) DeleteAnnotationJobMap(ctx context.Context, annotationID int, mapDescription string) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM annotation_job_map WHERE annotation_id = $1 AND description = $2
	`, annotationID, mapDescription)
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.NotFoundf("annotation map with  annotation_id: %d and description: %s not found", annotationID, mapDescription)
	}

	return nil
}

func (r *Repo) DeleteAnnotationServiceMap(ctx context.Context, annotationID int, mapDescription string) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM annotation_service_map WHERE annotation_id = $1 AND description = $2
	`, annotationID, mapDescription)
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.NotFoundf("annotation map with  annotation_id: %d and description: %s not found", annotationID, mapDescription)
	}

	return nil
}
