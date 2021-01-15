package data

import (
	"context"
	"database/sql"
	goErrors "errors"
	"time"

	"gitlab.unanet.io/devops/go/pkg/errors"
	"gitlab.unanet.io/devops/go/pkg/json"
)

type Label struct {
	ID          int          `db:"id"`
	Description string       `db:"description"`
	Data        json.Object  `db:"data"`
	CreatedAt   sql.NullTime `db:"created_at"`
	UpdatedAt   sql.NullTime `db:"updated_at"`
}

type LabelServiceMap struct {
	Description   string        `db:"description"`
	LabelID       int           `db:"label_id"`
	EnvironmentID sql.NullInt32 `db:"environment_id"`
	ArtifactID    sql.NullInt32 `db:"artifact_id"`
	NamespaceID   sql.NullInt32 `db:"namespace_id"`
	ServiceID     sql.NullInt32 `db:"service_id"`
	StackingOrder int           `db:"stacking_order"`
	CreatedAt     sql.NullTime  `db:"created_at"`
	UpdatedAt     sql.NullTime  `db:"updated_at"`
}

type LabelJobMap struct {
	Description   string        `db:"description"`
	LabelID       int           `db:"label_id"`
	EnvironmentID sql.NullInt32 `db:"environment_id"`
	ArtifactID    sql.NullInt32 `db:"artifact_id"`
	NamespaceID   sql.NullInt32 `db:"namespace_id"`
	JobID         sql.NullInt32 `db:"job_id"`
	StackingOrder int           `db:"stacking_order"`
	CreatedAt     sql.NullTime  `db:"created_at"`
	UpdatedAt     sql.NullTime  `db:"updated_at"`
}

type LabelService struct {
	LabelID          int           `db:"label_id"`
	Data             json.Object   `db:"data"`
	LabelDescription string        `db:"label_description"`
	MapDescription   string        `db:"map_description"`
	MapEnvironmentID sql.NullInt32 `db:"map_environment_id"`
	MapArtifactID    sql.NullInt32 `db:"map_artifact_id"`
	MapNamespaceID   sql.NullInt32 `db:"map_namespace_id"`
	MapServiceID     sql.NullInt32 `db:"map_service_id"`
	StackingOrder    int           `db:"stacking_order"`
	CreatedAt        sql.NullTime  `db:"created_at"`
	UpdatedAt        sql.NullTime  `db:"updated_at"`
}

type LabelJob struct {
	LabelID          int           `db:"label_id"`
	Data             json.Object   `db:"data"`
	LabelDescription string        `db:"label_description"`
	MapDescription   string        `db:"map_description"`
	MapEnvironmentID sql.NullInt32 `db:"map_environment_id"`
	MapArtifactID    sql.NullInt32 `db:"map_artifact_id"`
	MapNamespaceID   sql.NullInt32 `db:"map_namespace_id"`
	MapJobID         sql.NullInt32 `db:"map_job_id"`
	StackingOrder    int           `db:"stacking_order"`
	CreatedAt        sql.NullTime  `db:"created_at"`
	UpdatedAt        sql.NullTime  `db:"updated_at"`
}

func (r *Repo) UpsertMergeLabel(ctx context.Context, l *Label) error {
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
	INSERT INTO label(description, data, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (description)
		DO UPDATE SET data = label.data || $2, updated_at = $4
		RETURNING id, value, created_at
	`, l.Description, l.Data, l.CreatedAt, l.UpdatedAt).
		StructScan(l)

	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (r *Repo) UpsertLabel(ctx context.Context, l *Label) error {
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
	
	INSERT INTO label(description, data, created_at, updated_at)
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

func (r *Repo) UpsertLabelJobMap(ctx context.Context, ljm *LabelJobMap) error {
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
	
	INSERT INTO label_job_map(description, label_id, environment_id, artifact_id, namespace_id, job_id, stacking_order, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	ON CONFLICT (description)
	DO UPDATE SET environment_id = $3, artifact_id = $4, namespace_id = $5, job_id = $6, stacking_order = $7, updated_at = $9
	RETURNING created_at
	
	`, ljm.Description, ljm.LabelID, ljm.EnvironmentID, ljm.ArtifactID, ljm.NamespaceID, ljm.JobID, ljm.StackingOrder, ljm.CreatedAt, ljm.UpdatedAt).
		StructScan(ljm)

	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (r *Repo) UpsertLabelServiceMap(ctx context.Context, lsm *LabelServiceMap) error {
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
	
	INSERT INTO label_service_map(description, label_id, environment_id, artifact_id, namespace_id, service_id, stacking_order, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	ON CONFLICT (description)
	DO UPDATE SET environment_id = $3, artifact_id = $4, namespace_id = $5, service_id = $6, stacking_order = $7, updated_at = $9
	RETURNING created_at
	
	`, lsm.Description, lsm.LabelID, lsm.EnvironmentID, lsm.ArtifactID, lsm.NamespaceID, lsm.ServiceID, lsm.StackingOrder, lsm.CreatedAt, lsm.UpdatedAt).
		StructScan(lsm)

	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (r *Repo) DeleteLabel(ctx context.Context, labelID int) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM label WHERE id = $1
	`, labelID)
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.NotFoundf("label id: %d not found", labelID)
	}

	return nil
}

func (r *Repo) DeleteLabelKey(ctx context.Context, labelID int, key string) (*Label, error) {
	var label Label
	err := r.db.QueryRowxContext(ctx, `
		UPDATE label SET data = label.data - $1 WHERE id = $2
		RETURNING id, value, description, created_at, updated_at
	`, key, labelID).StructScan(&label)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return &label, nil
}

func (r *Repo) GetLabel(ctx context.Context, labelID int) (*Label, error) {
	var label Label

	row := r.db.QueryRowxContext(ctx, `
		select id, 
		       description, 
		       data, 
		       created_at, 
		       updated_at
		from label
		where id = $1
		`, labelID)
	err := row.StructScan(&label)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return nil, NotFoundErrorf("label with id: %d not found", labelID)
		}
		return nil, errors.Wrap(err)
	}

	return &label, nil
}

func (r *Repo) GetLabelByDescription(ctx context.Context, description string) (*Label, error) {
	var label Label

	row := r.db.QueryRowxContext(ctx, `
		select id, 
		       description, 
		       data, 
		       created_at, 
		       updated_at
		from label
		where description = $1
		`, description)
	err := row.StructScan(&label)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return nil, NotFoundErrorf("label with description: %s not found", description)
		}
		return nil, errors.Wrap(err)
	}

	return &label, nil
}

func (r *Repo) Labels(ctx context.Context) ([]Label, error) {
	rows, err := r.db.QueryxContext(ctx, `
		select id, 
		       description, 
		       data, 
		       created_at, 
		       updated_at 
		from label
	`)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var labels []Label
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}

		var label Label
		err = rows.StructScan(&label)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		labels = append(labels, label)
	}

	return labels, nil
}

func (r *Repo) JobLabelMaps(ctx context.Context, jobID int) ([]LabelJobMap, error) {
	rows, err := r.db.QueryxContext(ctx, `
		select description, 
		       label_id, 
		       environment_id, 
		       artifact_id, 
		       namespace_id, 
		       job_id, 
		       stacking_order, 
		       created_at, 
		       updated_at
		from label_job_map
		where job_id = $1
		`, jobID)

	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	// Hydrate a slice of the records to the Data Structure (PodAutoscaleMap)
	var labelJobMaps []LabelJobMap
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}
		var labelJobMap LabelJobMap
		err = rows.StructScan(&labelJobMap)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		labelJobMaps = append(labelJobMaps, labelJobMap)
	}

	return labelJobMaps, nil
}

func (r *Repo) JobLabelMapsByLabelID(ctx context.Context, labelID int) ([]LabelJobMap, error) {
	rows, err := r.db.QueryxContext(ctx, `
		select description, 
		       label_id, 
		       environment_id, 
		       artifact_id, 
		       namespace_id, 
		       job_id, 
		       stacking_order, 
		       created_at, 
		       updated_at
		from label_job_map
		where label_id = $1
		`, labelID)

	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	// Hydrate a slice of the records to the Data Structure (PodAutoscaleMap)
	var labelJobMaps []LabelJobMap
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}
		var labelJobMap LabelJobMap
		err = rows.StructScan(&labelJobMap)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		labelJobMaps = append(labelJobMaps, labelJobMap)
	}

	return labelJobMaps, nil
}

func (r *Repo) DeleteLabelJobMap(ctx context.Context, labelID int, mapDescription string) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM label_job_map WHERE label_id = $1 AND description = $2
	`, labelID, mapDescription)
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.NotFoundf("label map with  label_id: %d and description: %s not found", labelID, mapDescription)
	}

	return nil
}

func (r *Repo) DeleteLabelServiceMap(ctx context.Context, labelID int, mapDescription string) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM label_service_map WHERE label_id = $1 AND description = $2
	`, labelID, mapDescription)
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.NotFoundf("label map with  label_id: %d and description: %s not found", labelID, mapDescription)
	}

	return nil
}

func (r *Repo) ServiceLabelMapsByLabelID(ctx context.Context, labelID int) ([]LabelServiceMap, error) {
	rows, err := r.db.QueryxContext(ctx, `
		select description, 
		       label_id, 
		       environment_id, 
		       artifact_id, 
		       namespace_id, 
		       service_id, 
		       stacking_order, 
		       created_at, 
		       updated_at
		from label_service_map
		where label_id = $1
		`, labelID)

	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var labelServiceMaps []LabelServiceMap
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}
		var labelServiceMap LabelServiceMap
		err = rows.StructScan(&labelServiceMap)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		labelServiceMaps = append(labelServiceMaps, labelServiceMap)
	}

	return labelServiceMaps, nil
}

func (r *Repo) JobLabels(ctx context.Context, jobID int) ([]LabelJob, error) {
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
		
		SELECT l.id as label_id,
		       l.data as data,
		       l.description as label_description,
		       ljm.description as map_description,
		       ljm.environment_id as map_environment_id,
		       ljm.artifact_id as map_artifact_id,
		       ljm.namespace_id as map_namespace_id,
		       ljm.job_id as map_job_id,
		       ljm.stacking_order as stacking_order,
		       l.created_at,
		       l.updated_at
		FROM label_job_map ljm 
		    LEFT JOIN label l ON ljm.label_id = l.id 
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

	var labelJobs []LabelJob
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}
		var labelJob LabelJob
		err = rows.StructScan(&labelJob)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		labelJobs = append(labelJobs, labelJob)
	}

	return labelJobs, nil
}

func (r *Repo) ServiceLabels(ctx context.Context, serviceID int) ([]LabelService, error) {
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
		
		SELECT l.id as label_id,
		       l.data as data,
		       l.description as label_description,
		       lsm.description as map_description,
		       lsm.environment_id as map_environment_id,
		       lsm.artifact_id as map_artifact_id,
		       lsm.namespace_id as map_namespace_id,
		       lsm.service_id as map_service_id,
		       lsm.stacking_order as stacking_order,
		       l.created_at,
		       l.updated_at
		FROM label_service_map lsm 
		    LEFT JOIN label l ON lsm.label_id = l.id 
			LEFT JOIN env_data ed on ed.service_id = $1
		WHERE
			(lsm.service_id = $1)
		OR
			(lsm.cluster_id = ed.cluster_id AND lsm.artifact_id IS NULL)
		OR
		    (lsm.environment_id = ed.environment_id AND lsm.artifact_id IS NULL) 
		OR
		    (lsm.namespace_id = ed.namespace_id AND lsm.artifact_id IS NULL)
		OR
		    (lsm.artifact_id = ed.artifact_id AND lsm.environment_id IS NULL AND lsm.namespace_id IS NULL AND lsm.cluster_id IS NULL)
		OR
		    (lsm.artifact_id = ed.artifact_id AND lsm.cluster_id = ed.cluster_id)
		OR
		    (lsm.artifact_id = ed.artifact_id AND lsm.environment_id = ed.environment_id)
		OR
		    (lsm.artifact_id = ed.artifact_id AND lsm.namespace_id = ed.namespace_id)
		ORDER BY
			lsm.stacking_order
	`, serviceID)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var labelServices []LabelService
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}
		var labelService LabelService
		err = rows.StructScan(&labelService)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		labelServices = append(labelServices, labelService)
	}

	return labelServices, nil
}
