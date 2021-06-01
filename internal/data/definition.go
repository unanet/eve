package data

import (
	"context"
	"database/sql"
	goErrors "errors"
	"gitlab.unanet.io/devops/go/pkg/errors"
	"gitlab.unanet.io/devops/go/pkg/json"
	"time"
)

type Definition struct {
	ID               int          `db:"id"`
	Description      string       `db:"description"`
	DefinitionTypeID int          `db:"definition_type_id"`
	Data             json.Object  `db:"data"`
	CreatedAt        sql.NullTime `db:"created_at"`
	UpdatedAt        sql.NullTime `db:"updated_at"`
}

type DefinitionServiceMap struct {
	Description   string        `db:"description"`
	DefinitionID  int           `db:"definition_id"`
	EnvironmentID sql.NullInt32 `db:"environment_id"`
	ArtifactID    sql.NullInt32 `db:"artifact_id"`
	NamespaceID   sql.NullInt32 `db:"namespace_id"`
	ServiceID     sql.NullInt32 `db:"service_id"`
	ClusterID     sql.NullInt32 `db:"cluster_id"`
	StackingOrder int           `db:"stacking_order"`
	CreatedAt     sql.NullTime  `db:"created_at"`
	UpdatedAt     sql.NullTime  `db:"updated_at"`
}

type DefinitionJobMap struct {
	Description   string        `db:"description"`
	DefinitionID  int           `db:"definition_id"`
	EnvironmentID sql.NullInt32 `db:"environment_id"`
	ArtifactID    sql.NullInt32 `db:"artifact_id"`
	NamespaceID   sql.NullInt32 `db:"namespace_id"`
	JobID         sql.NullInt32 `db:"job_id"`
	ClusterID     sql.NullInt32 `db:"cluster_id"`
	StackingOrder int           `db:"stacking_order"`
	CreatedAt     sql.NullTime  `db:"created_at"`
	UpdatedAt     sql.NullTime  `db:"updated_at"`
}

type DefinitionService struct {
	DefinitionID          int           `db:"definition_id"`
	DefinitionTypeID      int           `db:"definition_type_id"`
	StackingOrder         int           `db:"stacking_order"`
	DefinitionType        string        `db:"definition_type"`
	DefinitionVersion     string        `db:"definition_version"`
	DefinitionClass       string        `db:"definition_class"`
	DefinitionKind        string        `db:"definition_kind"`
	DefinitionOrder       string        `db:"definition_order"`
	DefinitionDescription string        `db:"definition_description"`
	MapDescription        string        `db:"map_description"`
	Data                  json.Object   `db:"data"`
	MapEnvironmentID      sql.NullInt32 `db:"map_environment_id"`
	MapArtifactID         sql.NullInt32 `db:"map_artifact_id"`
	MapNamespaceID        sql.NullInt32 `db:"map_namespace_id"`
	MapServiceID          sql.NullInt32 `db:"map_service_id"`
	MapClusterID          sql.NullInt32 `db:"map_cluster_id"`
	CreatedAt             sql.NullTime  `db:"created_at"`
	UpdatedAt             sql.NullTime  `db:"updated_at"`
}

type DefinitionJob struct {
	DefinitionID          int           `db:"definition_id"`
	DefinitionTypeID      int           `db:"definition_type_id"`
	StackingOrder         int           `db:"stacking_order"`
	DefinitionType        string        `db:"definition_type"`
	DefinitionVersion     string        `db:"definition_version"`
	DefinitionClass       string        `db:"definition_class"`
	DefinitionKind        string        `db:"definition_kind"`
	DefinitionOrder       string        `db:"definition_order"`
	DefinitionDescription string        `db:"definition_description"`
	MapDescription        string        `db:"map_description"`
	Data                  json.Object   `db:"data"`
	MapEnvironmentId      sql.NullInt32 `db:"map_environment_id"`
	MapArtifactId         sql.NullInt32 `db:"map_artifact_id"`
	MapNamespaceId        sql.NullInt32 `db:"map_namespace_id"`
	MapJobId              sql.NullInt32 `db:"map_job_id"`
	MapClusterID          sql.NullInt32 `db:"map_cluster_id"`
	CreatedAt             sql.NullTime  `db:"created_at"`
	UpdatedAt             sql.NullTime  `db:"updated_at"`
}

func (r *Repo) UpsertMergeDefinition(ctx context.Context, def *Definition) error {
	now := time.Now().UTC()
	def.CreatedAt = sql.NullTime{
		Time:  now,
		Valid: true,
	}
	def.UpdatedAt = sql.NullTime{
		Time:  now,
		Valid: true,
	}

	err := r.db.QueryRowxContext(ctx, `
	INSERT INTO definition(description, definition_type_id, data, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (description)
		DO UPDATE SET data = definition.data || $3, updated_at = $5
		RETURNING id, data, created_at
	`, def.Description, def.DefinitionTypeID, def.Data, def.CreatedAt, def.UpdatedAt).
		StructScan(def)

	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (r *Repo) UpsertDefinition(ctx context.Context, def *Definition) error {
	now := time.Now().UTC()
	def.CreatedAt = sql.NullTime{
		Time:  now,
		Valid: true,
	}
	def.UpdatedAt = sql.NullTime{
		Time:  now,
		Valid: true,
	}

	err := r.db.QueryRowxContext(ctx, `
	
	INSERT INTO definition(description, definition_type_id,  data, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (description)
		DO UPDATE SET data = $3, updated_at = $5
		RETURNING id, created_at
	
	`, def.Description, def.DefinitionTypeID, def.Data, def.CreatedAt, def.UpdatedAt).
		StructScan(def)

	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (r *Repo) UpsertDefinitionJobMap(ctx context.Context, djm *DefinitionJobMap) error {
	now := time.Now().UTC()
	djm.CreatedAt = sql.NullTime{
		Time:  now,
		Valid: true,
	}
	djm.UpdatedAt = sql.NullTime{
		Time:  now,
		Valid: true,
	}

	err := r.db.QueryRowxContext(ctx, `

	INSERT INTO definition_job_map(description, definition_id, environment_id, artifact_id, namespace_id, job_id, cluster_id, stacking_order, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	ON CONFLICT (description)
	DO UPDATE SET environment_id = $3, artifact_id = $4, namespace_id = $5, job_id = $6, cluster_id = $7, stacking_order = $8, updated_at = $10
	RETURNING created_at
	
	`, djm.Description, djm.DefinitionID, djm.EnvironmentID, djm.ArtifactID, djm.NamespaceID, djm.JobID, djm.ClusterID, djm.StackingOrder, djm.CreatedAt, djm.UpdatedAt).
		StructScan(djm)

	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (r *Repo) UpsertDefinitionServiceMap(ctx context.Context, dsm *DefinitionServiceMap) error {
	now := time.Now().UTC()
	dsm.CreatedAt = sql.NullTime{
		Time:  now,
		Valid: true,
	}
	dsm.UpdatedAt = sql.NullTime{
		Time:  now,
		Valid: true,
	}

	err := r.db.QueryRowxContext(ctx, `
	
	INSERT INTO definition_service_map(description, definition_id, environment_id, artifact_id, namespace_id, service_id, cluster_id, stacking_order, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	ON CONFLICT (description)
	DO UPDATE SET environment_id = $3, artifact_id = $4, namespace_id = $5, service_id = $6, cluster_id = $7, stacking_order = $8, updated_at = $10
	RETURNING created_at
	
	`, dsm.Description, dsm.DefinitionID, dsm.EnvironmentID, dsm.ArtifactID, dsm.NamespaceID, dsm.ServiceID, dsm.ClusterID, dsm.StackingOrder, dsm.CreatedAt, dsm.UpdatedAt).
		StructScan(dsm)

	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (r *Repo) GetDefinition(ctx context.Context, definitionID int) (*Definition, error) {
	var definition Definition

	row := r.db.QueryRowxContext(ctx, `
		select id, 
		       description, 
		       definition_type_id,
		       data, 
		       created_at, 
		       updated_at
		from definition
		where id = $1
		`, definitionID)
	err := row.StructScan(&definition)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return nil, NotFoundErrorf("definition with id: %d not found", definitionID)
		}
		return nil, errors.Wrap(err)
	}

	return &definition, nil
}

func (r *Repo) GetDefinitionByDescription(ctx context.Context, description string) (*Definition, error) {
	var definition Definition

	row := r.db.QueryRowxContext(ctx, `
		select id, 
		       description, 
		       definition_type_id,
		       data, 
		       created_at, 
		       updated_at
		from definition
		where description = $1
		`, description)
	err := row.StructScan(&definition)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return nil, NotFoundErrorf("definition with description: %s not found", description)
		}
		return nil, errors.Wrap(err)
	}

	return &definition, nil
}

func (r *Repo) Definition(ctx context.Context) ([]Definition, error) {
	rows, err := r.db.QueryxContext(ctx, `
		select id, 
		       description, 
		       definition_type_id,
		       data, 
		       created_at, 
		       updated_at 
		from definition
	`)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var ms []Definition
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}

		var m Definition
		err = rows.StructScan(&m)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		ms = append(ms, m)
	}

	return ms, nil
}

func (r *Repo) DeleteDefinitionKey(ctx context.Context, definitionID int, key string) (*Definition, error) {
	var definition Definition
	err := r.db.QueryRowxContext(ctx, `
		UPDATE definition SET data = definition.data - $1 WHERE id = $2
		RETURNING id, data, description,definition_type_id, created_at, updated_at
	`, key, definitionID).StructScan(&definition)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return &definition, nil
}

func (r *Repo) DeleteDefinition(ctx context.Context, definitionID int) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM definition WHERE id = $1
	`, definitionID)
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.NotFoundf("definition id: %d not found", definitionID)
	}

	return nil
}

func (r *Repo) DeleteDefinitionJobMap(ctx context.Context, definitionID int, mapDescription string) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM definition_job_map WHERE definition_id = $1 AND description = $2
	`, definitionID, mapDescription)
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.NotFoundf("definition map with  definition_id: %d and description: %s not found", definitionID, mapDescription)
	}

	return nil
}

func (r *Repo) DeleteDefinitionServiceMap(ctx context.Context, definitionID int, mapDescription string) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM definition_service_map WHERE definition_id = $1 AND description = $2
	`, definitionID, mapDescription)
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.NotFoundf("definition map with  definition_id: %d and description: %s not found", definitionID, mapDescription)
	}

	return nil
}

func (r *Repo) JobDefinitionMapsByJobID(ctx context.Context, jobID int) ([]DefinitionJobMap, error) {
	rows, err := r.db.QueryxContext(ctx, `
		select description, 
		       definition_id, 
		       environment_id, 
		       artifact_id, 
		       namespace_id, 
		       job_id,
		       cluster_id,
		       stacking_order, 
		       created_at, 
		       updated_at
		from definition_job_map
		where job_id = $1
		`, jobID)

	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var mjms []DefinitionJobMap
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}
		var mjm DefinitionJobMap
		err = rows.StructScan(&mjm)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		mjms = append(mjms, mjm)
	}

	return mjms, nil
}

func (r *Repo) DefinitionJobMaps(ctx context.Context) ([]DefinitionJobMap, error) {
	rows, err := r.db.QueryxContext(ctx, `
		select
			description,
			definition_id,
			environment_id,
			artifact_id,
			namespace_id,
			job_id,
			cluster_id,
			stacking_order,
			created_at,
			updated_at
		from definition_job_map
		`)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var ss []DefinitionJobMap
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}

		var s DefinitionJobMap
		err = rows.StructScan(&s)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		ss = append(ss, s)
	}

	return ss, nil
}

func (r *Repo) DefinitionServiceMaps(ctx context.Context) ([]DefinitionServiceMap, error) {
	rows, err := r.db.QueryxContext(ctx, `
		select 
			description,
			definition_id,
			environment_id,
			artifact_id,
			namespace_id,
			service_id,
			cluster_id,
			stacking_order,
			created_at,
			updated_at
		from definition_service_map`)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var ss []DefinitionServiceMap
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}

		var s DefinitionServiceMap
		err = rows.StructScan(&s)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		ss = append(ss, s)
	}

	return ss, nil
}

func (r *Repo) JobDefinitionMapsByDefinitionID(ctx context.Context, definitionID int) ([]DefinitionJobMap, error) {
	rows, err := r.db.QueryxContext(ctx, `
		select description, 
		       definition_id, 
		       environment_id, 
		       artifact_id, 
		       namespace_id, 
		       job_id, 
		       cluster_id, 
		       stacking_order, 
		       created_at, 
		       updated_at
		from definition_job_map
		where definition_id = $1
		`, definitionID)

	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var mjms []DefinitionJobMap
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}
		var mjm DefinitionJobMap
		err = rows.StructScan(&mjm)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		mjms = append(mjms, mjm)
	}

	return mjms, nil
}

func (r *Repo) ServiceDefinitionMapsByDefinitionID(ctx context.Context, definitionID int) ([]DefinitionServiceMap, error) {
	rows, err := r.db.QueryxContext(ctx, `
		select description, 
		       definition_id, 
		       environment_id, 
		       artifact_id, 
		       namespace_id, 
		       service_id, 
		       cluster_id, 
		       stacking_order, 
		       created_at, 
		       updated_at
		from definition_service_map
		where definition_id = $1
		`, definitionID)

	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var msms []DefinitionServiceMap
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}
		var msm DefinitionServiceMap
		err = rows.StructScan(&msm)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		msms = append(msms, msm)
	}

	return msms, nil
}

func (r *Repo) JobDefinition(ctx context.Context, jobID int) ([]DefinitionJob, error) {
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
		
		SELECT d.id as definition_id,
  			   dt.name as definition_type,
		       dt.id as definition_type_id,
		       dt.class as definition_class,
		       dt.definition_order as definition_order,
		       dt.kind as definition_kind,
		       dt.version as definition_version,
		       d.data as data,
		       d.description as definition_description,
		       djm.description as map_description,
		       djm.environment_id as map_environment_id,
		       djm.artifact_id as map_artifact_id,
		       djm.namespace_id as map_namespace_id,
		       djm.job_id as map_job_id,
		       djm.cluster_id as map_cluster_id,
		       djm.stacking_order as stacking_order,
		       d.created_at,
		       d.updated_at
		FROM definition_job_map djm 
		    LEFT JOIN definition d ON djm.definition_id = d.id 
			LEFT JOIN env_data ed on ed.job_id = $1
			LEFT OUTER JOIN definition_type dt on d.definition_type_id = dt.id
		WHERE
			(djm.job_id = $1)
		OR
			(djm.cluster_id = ed.cluster_id AND djm.artifact_id IS NULL)
		OR
		    (djm.environment_id = ed.environment_id AND djm.artifact_id IS NULL) 
		OR
		    (djm.namespace_id = ed.namespace_id AND djm.artifact_id IS NULL)
		OR
		    (djm.artifact_id = ed.artifact_id AND djm.environment_id IS NULL AND djm.namespace_id IS NULL AND djm.cluster_id IS NULL)
		OR
		    (djm.artifact_id = ed.artifact_id AND djm.cluster_id = ed.cluster_id)
		OR
		    (djm.artifact_id = ed.artifact_id AND djm.environment_id = ed.environment_id)
		OR
		    (djm.artifact_id = ed.artifact_id AND djm.namespace_id = ed.namespace_id)
		ORDER BY
			djm.stacking_order
	`, jobID)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var djms []DefinitionJob
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}
		var djm DefinitionJob
		err = rows.StructScan(&djm)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		djms = append(djms, djm)
	}

	return djms, nil
}

func (r *Repo) ServiceDefinition(ctx context.Context, serviceID int) ([]DefinitionService, error) {
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
		
		SELECT d.id as definition_id,
		       dt.name as definition_type,
		       dt.id as definition_type_id,
		       dt.class as definition_class,
		       dt.definition_order as definition_order,
		       dt.kind as definition_kind,
		       dt.version as definition_version,
		       d.data as data,
		       d.description as definition_description,
		       dsm.description as map_description,
		       dsm.environment_id as map_environment_id,
		       dsm.artifact_id as map_artifact_id,
		       dsm.namespace_id as map_namespace_id,
		       dsm.service_id as map_service_id,
		       dsm.cluster_id as map_cluster_id,
		       dsm.stacking_order as stacking_order,
		       d.created_at,
		       d.updated_at
		FROM definition_service_map dsm 
		    LEFT JOIN definition d ON dsm.definition_id = d.id 
			LEFT JOIN env_data ed on ed.service_id = $1
			LEFT OUTER JOIN definition_type dt on d.definition_type_id = dt.id
		WHERE
			(dsm.service_id = $1)
		OR
			(dsm.cluster_id = ed.cluster_id AND dsm.artifact_id IS NULL)
		OR
		    (dsm.environment_id = ed.environment_id AND dsm.artifact_id IS NULL) 
		OR
		    (dsm.namespace_id = ed.namespace_id AND dsm.artifact_id IS NULL)
		OR
		    (dsm.artifact_id = ed.artifact_id AND dsm.environment_id IS NULL AND dsm.namespace_id IS NULL AND dsm.cluster_id IS NULL)
		OR
		    (dsm.artifact_id = ed.artifact_id AND dsm.cluster_id = ed.cluster_id)
		OR
		    (dsm.artifact_id = ed.artifact_id AND dsm.environment_id = ed.environment_id)
		OR
		    (dsm.artifact_id = ed.artifact_id AND dsm.namespace_id = ed.namespace_id)
		ORDER BY
			dsm.stacking_order
	`, serviceID)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var dsms []DefinitionService
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}
		var dsm DefinitionService
		err = rows.StructScan(&dsm)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		dsms = append(dsms, dsm)
	}

	return dsms, nil
}