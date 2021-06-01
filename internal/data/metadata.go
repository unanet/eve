package data

import (
	"context"
	"database/sql"
	goErrors "errors"
	"time"

	"gitlab.unanet.io/devops/go/pkg/errors"
	"gitlab.unanet.io/devops/go/pkg/json"
)

type Metadata struct {
	ID           int          `db:"id"`
	Description  string       `db:"description"`
	Value        json.Object  `db:"value"`
	MigratedFrom int          `db:"migrated_from"`
	CreatedAt    sql.NullTime `db:"created_at"`
	UpdatedAt    sql.NullTime `db:"updated_at"`
}

type MetadataServiceMap struct {
	Description   string        `db:"description"`
	MetadataID    int           `db:"metadata_id"`
	EnvironmentID sql.NullInt32 `db:"environment_id"`
	ArtifactID    sql.NullInt32 `db:"artifact_id"`
	NamespaceID   sql.NullInt32 `db:"namespace_id"`
	ClusterID     sql.NullInt32 `db:"cluster_id"`
	ServiceID     sql.NullInt32 `db:"service_id"`
	StackingOrder int           `db:"stacking_order"`
	CreatedAt     sql.NullTime  `db:"created_at"`
	UpdatedAt     sql.NullTime  `db:"updated_at"`
}

type MetadataJobMap struct {
	Description   string        `db:"description"`
	MetadataID    int           `db:"metadata_id"`
	EnvironmentID sql.NullInt32 `db:"environment_id"`
	ArtifactID    sql.NullInt32 `db:"artifact_id"`
	NamespaceID   sql.NullInt32 `db:"namespace_id"`
	ClusterID     sql.NullInt32 `db:"cluster_id"`
	JobID         sql.NullInt32 `db:"job_id"`
	StackingOrder int           `db:"stacking_order"`
	CreatedAt     sql.NullTime  `db:"created_at"`
	UpdatedAt     sql.NullTime  `db:"updated_at"`
}

type MetadataService struct {
	MetadataID          int           `db:"metadata_id"`
	Metadata            json.Object   `db:"metadata"`
	MetadataDescription string        `db:"metadata_description"`
	MapDescription      string        `db:"map_description"`
	MapEnvironmentID    sql.NullInt32 `db:"map_environment_id"`
	MapArtifactID       sql.NullInt32 `db:"map_artifact_id"`
	MapNamespaceID      sql.NullInt32 `db:"map_namespace_id"`
	MapServiceID        sql.NullInt32 `db:"map_service_id"`
	StackingOrder       int           `db:"stacking_order"`
	CreatedAt           sql.NullTime  `db:"created_at"`
	UpdatedAt           sql.NullTime  `db:"updated_at"`
}

type MetadataJob struct {
	MetadataID          int           `db:"metadata_id"`
	Metadata            json.Object   `db:"metadata"`
	MetadataDescription string        `db:"metadata_description"`
	MapDescription      string        `db:"map_description"`
	MapEnvironmentId    sql.NullInt32 `db:"map_environment_id"`
	MapArtifactId       sql.NullInt32 `db:"map_artifact_id"`
	MapNamespaceId      sql.NullInt32 `db:"map_namespace_id"`
	MapJobId            sql.NullInt32 `db:"map_job_id"`
	StackingOrder       int           `db:"stacking_order"`
	CreatedAt           sql.NullTime  `db:"created_at"`
	UpdatedAt           sql.NullTime  `db:"updated_at"`
}

func (r *Repo) UpsertMergeMetadata(ctx context.Context, m *Metadata) error {
	now := time.Now().UTC()
	m.CreatedAt = sql.NullTime{
		Time:  now,
		Valid: true,
	}
	m.UpdatedAt = sql.NullTime{
		Time:  now,
		Valid: true,
	}

	err := r.db.QueryRowxContext(ctx, `
	INSERT INTO metadata(description, value, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (description)
		DO UPDATE SET value = metadata.value || $2, updated_at = $4
		RETURNING id, value, created_at
	`, m.Description, m.Value, m.CreatedAt, m.UpdatedAt).
		StructScan(m)

	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (r *Repo) UpsertMetadata(ctx context.Context, m *Metadata) error {
	now := time.Now().UTC()
	m.CreatedAt = sql.NullTime{
		Time:  now,
		Valid: true,
	}
	m.UpdatedAt = sql.NullTime{
		Time:  now,
		Valid: true,
	}

	err := r.db.QueryRowxContext(ctx, `
	
	INSERT INTO metadata(description, value, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (description)
		DO UPDATE SET value = $2, updated_at = $4
		RETURNING id, created_at
	
	`, m.Description, m.Value, m.CreatedAt, m.UpdatedAt).
		StructScan(m)

	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (r *Repo) UpsertMetadataJobMap(ctx context.Context, mjm *MetadataJobMap) error {
	now := time.Now().UTC()
	mjm.CreatedAt = sql.NullTime{
		Time:  now,
		Valid: true,
	}
	mjm.UpdatedAt = sql.NullTime{
		Time:  now,
		Valid: true,
	}

	err := r.db.QueryRowxContext(ctx, `
	
	INSERT INTO metadata_job_map(description, metadata_id, environment_id, artifact_id, namespace_id, job_id, stacking_order, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	ON CONFLICT (description)
	DO UPDATE SET environment_id = $3, artifact_id = $4, namespace_id = $5, job_id = $6, stacking_order = $7, updated_at = $9
	RETURNING created_at
	
	`, mjm.Description, mjm.MetadataID, mjm.EnvironmentID, mjm.ArtifactID, mjm.NamespaceID, mjm.JobID, mjm.StackingOrder, mjm.CreatedAt, mjm.UpdatedAt).
		StructScan(mjm)

	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (r *Repo) UpsertMetadataServiceMap(ctx context.Context, msm *MetadataServiceMap) error {
	now := time.Now().UTC()
	msm.CreatedAt = sql.NullTime{
		Time:  now,
		Valid: true,
	}
	msm.UpdatedAt = sql.NullTime{
		Time:  now,
		Valid: true,
	}

	err := r.db.QueryRowxContext(ctx, `
	
	INSERT INTO metadata_service_map(description, metadata_id, environment_id, artifact_id, namespace_id, service_id, stacking_order, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	ON CONFLICT (description)
	DO UPDATE SET environment_id = $3, artifact_id = $4, namespace_id = $5, service_id = $6, stacking_order = $7, updated_at = $9
	RETURNING created_at
	
	`, msm.Description, msm.MetadataID, msm.EnvironmentID, msm.ArtifactID, msm.NamespaceID, msm.ServiceID, msm.StackingOrder, msm.CreatedAt, msm.UpdatedAt).
		StructScan(msm)

	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (r *Repo) GetMetadata(ctx context.Context, metadataID int) (*Metadata, error) {
	var metadata Metadata

	row := r.db.QueryRowxContext(ctx, `
		select id, 
		       description, 
		       value, 
		       created_at, 
		       updated_at
		from metadata
		where id = $1
		`, metadataID)
	err := row.StructScan(&metadata)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return nil, NotFoundErrorf("metadata with id: %d not found", metadataID)
		}
		return nil, errors.Wrap(err)
	}

	return &metadata, nil
}

func (r *Repo) GetMetadataByDescription(ctx context.Context, description string) (*Metadata, error) {
	var metadata Metadata

	row := r.db.QueryRowxContext(ctx, `
		select id, 
		       description, 
		       value, 
		       created_at, 
		       updated_at
		from metadata
		where description = $1
		`, description)
	err := row.StructScan(&metadata)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return nil, NotFoundErrorf("metadata with description: %s not found", description)
		}
		return nil, errors.Wrap(err)
	}

	return &metadata, nil
}

func (r *Repo) Metadata(ctx context.Context) ([]Metadata, error) {
	rows, err := r.db.QueryxContext(ctx, `
		select id, 
		       description, 
		       value, 
		       created_at, 
		       updated_at 
		from metadata
	`)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var ms []Metadata
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}

		var m Metadata
		err = rows.StructScan(&m)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		ms = append(ms, m)
	}

	return ms, nil
}

func (r *Repo) DeleteMetadataKey(ctx context.Context, metadataID int, key string) (*Metadata, error) {
	var metadata Metadata
	err := r.db.QueryRowxContext(ctx, `
		UPDATE metadata SET value = metadata.value - $1 WHERE id = $2
		RETURNING id, value, description, created_at, updated_at
	`, key, metadataID).StructScan(&metadata)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return &metadata, nil
}

func (r *Repo) DeleteMetadata(ctx context.Context, metadataID int) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM metadata WHERE id = $1
	`, metadataID)
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.NotFoundf("metadata id: %d not found", metadataID)
	}

	return nil
}

func (r *Repo) DeleteMetadataJobMap(ctx context.Context, metadataID int, mapDescription string) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM metadata_job_map WHERE metadata_id = $1 AND description = $2
	`, metadataID, mapDescription)
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.NotFoundf("metadata map with  metadata_id: %d and description: %s not found", metadataID, mapDescription)
	}

	return nil
}

func (r *Repo) DeleteMetadataServiceMap(ctx context.Context, metadataID int, mapDescription string) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM metadata_service_map WHERE metadata_id = $1 AND description = $2
	`, metadataID, mapDescription)
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.NotFoundf("metadata map with  metadata_id: %d and description: %s not found", metadataID, mapDescription)
	}

	return nil
}

func (r *Repo) JobMetadataMapsByJobID(ctx context.Context, jobID int) ([]MetadataJobMap, error) {
	rows, err := r.db.QueryxContext(ctx, `
		select description, 
		       metadata_id, 
		       environment_id, 
		       artifact_id, 
		       namespace_id, 
		       job_id, 
		       stacking_order, 
		       created_at, 
		       updated_at
		from metadata_job_map
		where job_id = $1
		`, jobID)

	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var mjms []MetadataJobMap
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}
		var mjm MetadataJobMap
		err = rows.StructScan(&mjm)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		mjms = append(mjms, mjm)
	}

	return mjms, nil
}

func (r *Repo) MetadataJobMaps(ctx context.Context) ([]MetadataJobMap, error) {
	rows, err := r.db.QueryxContext(ctx, `
		select 
			description,
			metadata_id,
			environment_id,
			artifact_id,
			namespace_id,
			job_id,
			stacking_order,
			created_at,
			updated_at
		from metadata_job_map`)

	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var ss []MetadataJobMap
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}

		var s MetadataJobMap
		err = rows.StructScan(&s)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		ss = append(ss, s)
	}

	return ss, nil
}

func (r *Repo) MetadataServiceMaps(ctx context.Context) ([]MetadataServiceMap, error) {
	rows, err := r.db.QueryxContext(ctx, `
		select 
			description,
			metadata_id,
			environment_id,
			artifact_id,
			namespace_id,
			service_id,
			stacking_order,
			created_at,
			updated_at
		from metadata_service_map`)

	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var ss []MetadataServiceMap
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}

		var s MetadataServiceMap
		err = rows.StructScan(&s)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		ss = append(ss, s)
	}

	return ss, nil
}

func (r *Repo) JobMetadataMapsByMetadataID(ctx context.Context, metadataID int) ([]MetadataJobMap, error) {
	rows, err := r.db.QueryxContext(ctx, `
		select description, 
		       metadata_id, 
		       environment_id, 
		       artifact_id, 
		       namespace_id, 
		       job_id, 
		       stacking_order, 
		       created_at, 
		       updated_at
		from metadata_job_map
		where metadata_id = $1
		`, metadataID)

	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var mjms []MetadataJobMap
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}
		var mjm MetadataJobMap
		err = rows.StructScan(&mjm)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		mjms = append(mjms, mjm)
	}

	return mjms, nil
}

func (r *Repo) ServiceMetadataMapsByMetadataID(ctx context.Context, metadataID int) ([]MetadataServiceMap, error) {
	rows, err := r.db.QueryxContext(ctx, `
		select description, 
		       metadata_id, 
		       environment_id, 
		       artifact_id, 
		       namespace_id, 
		       service_id, 
		       stacking_order, 
		       created_at, 
		       updated_at
		from metadata_service_map
		where metadata_id = $1
		`, metadataID)

	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var msms []MetadataServiceMap
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}
		var msm MetadataServiceMap
		err = rows.StructScan(&msm)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		msms = append(msms, msm)
	}

	return msms, nil
}

func (r *Repo) JobMetadata(ctx context.Context, jobID int) ([]MetadataJob, error) {
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
		
		SELECT m.id as metadata_id,
		       m.value as metadata,
		       m.description as metadata_description,
		       mjm.description as map_description,
		       mjm.environment_id as map_environment_id,
		       mjm.artifact_id as map_artifact_id,
		       mjm.namespace_id as map_namespace_id,
		       mjm.job_id as map_job_id,
		       mjm.stacking_order as stacking_order,
		       m.created_at,
		       m.updated_at
		FROM metadata_job_map mjm 
		    LEFT JOIN metadata m ON mjm.metadata_id = m.id 
			LEFT JOIN env_data ed on ed.job_id = $1
		WHERE
			(mjm.job_id = $1)
		OR
			(mjm.cluster_id = ed.cluster_id AND mjm.artifact_id IS NULL)
		OR
		    (mjm.environment_id = ed.environment_id AND mjm.artifact_id IS NULL) 
		OR
		    (mjm.namespace_id = ed.namespace_id AND mjm.artifact_id IS NULL)
		OR
		    (mjm.artifact_id = ed.artifact_id AND mjm.environment_id IS NULL AND mjm.namespace_id IS NULL AND mjm.cluster_id IS NULL)
		OR
		    (mjm.artifact_id = ed.artifact_id AND mjm.cluster_id = ed.cluster_id)
		OR
		    (mjm.artifact_id = ed.artifact_id AND mjm.environment_id = ed.environment_id)
		OR
		    (mjm.artifact_id = ed.artifact_id AND mjm.namespace_id = ed.namespace_id)
		ORDER BY
			mjm.stacking_order
	`, jobID)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var mjms []MetadataJob
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}
		var mjm MetadataJob
		err = rows.StructScan(&mjm)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		mjms = append(mjms, mjm)
	}

	return mjms, nil
}

func (r *Repo) ServiceMetadata(ctx context.Context, serviceID int) ([]MetadataService, error) {
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
		
		SELECT m.id as metadata_id,
		       m.value as metadata,
		       m.description as metadata_description,
		       msm.description as map_description,
		       msm.environment_id as map_environment_id,
		       msm.artifact_id as map_artifact_id,
		       msm.namespace_id as map_namespace_id,
		       msm.service_id as map_service_id,
		       msm.stacking_order as stacking_order,
		       m.created_at,
		       m.updated_at
		FROM metadata_service_map msm 
		    LEFT JOIN metadata m ON msm.metadata_id = m.id 
			LEFT JOIN env_data ed on ed.service_id = $1
		WHERE
			(msm.service_id = $1)
		OR
			(msm.cluster_id = ed.cluster_id AND msm.artifact_id IS NULL)
		OR
		    (msm.environment_id = ed.environment_id AND msm.artifact_id IS NULL) 
		OR
		    (msm.namespace_id = ed.namespace_id AND msm.artifact_id IS NULL)
		OR
		    (msm.artifact_id = ed.artifact_id AND msm.environment_id IS NULL AND msm.namespace_id IS NULL AND msm.cluster_id IS NULL)
		OR
		    (msm.artifact_id = ed.artifact_id AND msm.cluster_id = ed.cluster_id)
		OR
		    (msm.artifact_id = ed.artifact_id AND msm.environment_id = ed.environment_id)
		OR
		    (msm.artifact_id = ed.artifact_id AND msm.namespace_id = ed.namespace_id)
		ORDER BY
			msm.stacking_order
	`, serviceID)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var msms []MetadataService
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}
		var msm MetadataService
		err = rows.StructScan(&msm)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		msms = append(msms, msm)
	}

	return msms, nil
}


func (r *Repo) CreateMetadataJobMap(ctx context.Context, model *MetadataJobMap) error {
	model.CreatedAt = sql.NullTime{
		Time:  time.Now().UTC(),
		Valid: true,
	}

	err := r.db.QueryRowxContext(ctx,`
	INSERT INTO metadata_job_map(
					description,
					metadata_id,
					environment_id,
					artifact_id,
					namespace_id,
					 cluster_id,
					job_id,
					stacking_order,
					created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at
	`,
		model.Description,
		model.MetadataID,
		model.EnvironmentID,
		model.ArtifactID,
		model.NamespaceID,
		model.ClusterID,
		model.JobID,
		model.StackingOrder,
		model.CreatedAt).
		StructScan(model)

	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (r *Repo) CreateMetadataServiceMap(ctx context.Context, model *MetadataServiceMap) error {
	model.CreatedAt = sql.NullTime{
		Time:  time.Now().UTC(),
		Valid: true,
	}

	err := r.db.QueryRowxContext(ctx,`
	INSERT INTO metadata_service_map(
					description,
					metadata_id,
					environment_id,
					artifact_id,
					namespace_id,
					service_id,
					cluster_id,
					stacking_order,
					created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at
	`,
		model.Description,
		model.MetadataID,
		model.EnvironmentID,
		model.ArtifactID,
		model.NamespaceID,
		model.ClusterID,
		model.ServiceID,
		model.StackingOrder,
		model.CreatedAt).
		StructScan(model)

	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}
