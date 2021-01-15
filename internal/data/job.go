package data

import (
	"context"
	"database/sql"
	goErrors "errors"
	"time"

	"github.com/jmoiron/sqlx"

	"gitlab.unanet.io/devops/go/pkg/errors"
)

type DeployJob struct {
	JobID            int            `db:"job_id"`
	JobName          string         `db:"job_name"`
	ArtifactID       int            `db:"artifact_id"`
	ArtifactName     string         `db:"artifact_name"`
	RequestedVersion string         `db:"requested_version"`
	DeployedVersion  sql.NullString `db:"deployed_version"`
	ServiceAccount   string         `db:"service_account"`
	ImageTag         string         `db:"image_tag"`
	RunAs            int            `db:"run_as"`
	NodeGroup        string         `db:"node_group"`
	EnvironmentID    int            `db:"environment_id"`
	NamespaceID      int            `db:"namespace_id"`
	CreatedAt        sql.NullTime   `db:"created_at"`
	UpdatedAt        sql.NullTime   `db:"updated_at"`
	SuccessExitCodes string         `db:"success_exit_codes"`
}

type DeployJobs []DeployJob

type Job struct {
	ID              int            `db:"id"`
	NamespaceID     int            `db:"namespace_id"`
	NamespaceName   string         `db:"namespace_name"`
	ArtifactID      int            `db:"artifact_id"`
	ArtifactName    string         `db:"artifact_name"`
	OverrideVersion sql.NullString `db:"override_version"`
	DeployedVersion sql.NullString `db:"deployed_version"`
	CreatedAt       sql.NullTime   `db:"created_at"`
	UpdatedAt       sql.NullTime   `db:"updated_at"`
	Name            string         `db:"name"`
	NodeGroup       string         `db:"node_group"`
}

func (r *Repo) UpdateDeployedJobVersion(ctx context.Context, id int, version string) error {
	result, err := r.db.ExecContext(ctx, `
		update job
		set deployed_version = $1, 
		    updated_at = $2 
		where id = $3
	`, version, time.Now().UTC(), id)
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.Wrapf("the following id: %d was not found to update in job table", id)
	}
	return nil
}

func (r *Repo) DeployedJobsByNamespaceID(ctx context.Context, namespaceID int) (DeployJobs, error) {
	rows, err := r.db.QueryxContext(ctx, `
		select j.id as job_id,
		       j.name as job_name,
		       j.artifact_id,
		       a.name as artifact_name,
		       COALESCE(j.override_version, n.requested_version) as requested_version,
		       j.deployed_version,
		       a.service_account,
		       a.image_tag,
		       a.run_as as run_as,
		       j.node_group,
		       e.id as environment_id,
		       n.id as namespace_id,
		       j.success_exit_codes,
		       j.created_at,
		       j.updated_at
		from job as j 
		    left join artifact as a on a.id = j.artifact_id
			left join namespace n on j.namespace_id = n.id
			left join environment e on n.environment_id = e.id
		where j.namespace_id = $1

	`, namespaceID)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	var jobs DeployJobs

	for rows.Next() {
		var job DeployJob
		err = rows.StructScan(&job)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (r *Repo) JobByName(ctx context.Context, name string, namespace string) (*Job, error) {
	var job Job

	row := r.db.QueryRowxContext(ctx, `
		select j.id, 
		       j.name, 
		       j.namespace_id, 
		       j.artifact_id, 
		       j.override_version,
		       j.created_at,
		       j.updated_at,
		       j.node_group,
		       n.name as namespace_name, 
		       a.name as artifact_name
		from job j 
		    left join namespace n on j.namespace_id = n.id
			left join artifact a on j.artifact_id = a.id
		where j.name = $1 and n.name = $2
		`, name, namespace)
	err := row.StructScan(&job)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return nil, NotFoundErrorf("job with name: %s, namespace: %s, not found", name, namespace)
		}
		return nil, errors.Wrap(err)
	}

	return &job, nil
}

func (r *Repo) JobArtifacts(ctx context.Context, namespaceIDs []int) (RequestArtifacts, error) {
	s, args, err := sqlx.In(`
		select distinct j.artifact_id, 
		                a.function_pointer as function_pointer,
		                a.name as artifact_name, 
		                a.provider_group as provider_group,
		                a.feed_type as feed_type,
		                f.name as feed_name,
		                COALESCE(j.override_version, ns.requested_version) as requested_version 
		from job as j 
			left join namespace as ns on ns.id = j.namespace_id
			left join artifact as a on a.id = j.artifact_id
			left join environment e on ns.environment_id = e.id
			left join environment_feed_map efm on e.id = efm.environment_id
			left join feed f on efm.feed_id = f.id and f.feed_type = a.feed_type
		where f.name is not null and ns.id in (?)`, namespaceIDs)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	s = r.db.Rebind(s)
	rows, err := r.db.QueryxContext(ctx, s, args...)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	var jobs []RequestArtifact
	for rows.Next() {
		var job RequestArtifact
		err = rows.StructScan(&job)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (r *Repo) JobByID(ctx context.Context, id int) (*Job, error) {
	var job Job

	row := r.db.QueryRowxContext(ctx, `
		select j.id, 
		       j.name, 
		       j.namespace_id, 
		       j.artifact_id, 
		       j.override_version,
		       j.created_at,
		       j.updated_at,
		       j.node_group,
		       n.name as namespace_name, 
		       a.name as artifact_name
		from job j 
		    left join namespace n on j.namespace_id = n.id
			left join artifact a on j.artifact_id = a.id
		where j.id = $1
		`, id)
	err := row.StructScan(&job)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return nil, NotFoundErrorf("job with id: %d, not found", id)
		}
		return nil, errors.Wrap(err)
	}

	return &job, nil
}

func (r *Repo) JobsByNamespaceID(ctx context.Context, namespaceID int) ([]Job, error) {
	return r.jobs(ctx, Where("n.id", namespaceID))
}

func (r *Repo) JobsByNamespaceName(ctx context.Context, namespaceName string) ([]Job, error) {
	return r.jobs(ctx, Where("n.name", namespaceName))
}

func (r *Repo) jobs(ctx context.Context, whereArgs ...WhereArg) ([]Job, error) {
	s, args := CheckWhereArgs(`
		select j.id, 
		       j.namespace_id, 
		       j.artifact_id, 
		       j.override_version, 
		       j.deployed_version, 
		       j.created_at, 
		       j.updated_at, 
		       j.name, 
               j.node_group,
		       n.name as namespace_name,
		       a.name as artifact_name
		from job j 
		    left join namespace n on j.namespace_id = n.id
			left join artifact a on j.artifact_id = a.id
		`, whereArgs)
	rows, err := r.db.QueryxContext(ctx, s+"order by j.name", args...)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	var jobs []Job
	for rows.Next() {
		var job Job
		err = rows.StructScan(&job)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (r *Repo) UpdateJob(ctx context.Context, job *Job) error {
	job.UpdatedAt.Time = time.Now().UTC()
	job.UpdatedAt.Valid = true

	result, err := r.db.ExecContext(ctx, `
		update job set 
		   	name = $1, 
			namespace_id = $2,
			artifact_id = $3,
		   	override_version = $4,
		    deployed_version = $5,
		    updated_at = $6,
		    node_group = $7
		where id = $8
	`,
		job.Name,
		job.NamespaceID,
		job.ArtifactID,
		job.OverrideVersion,
		job.DeployedVersion,
		job.UpdatedAt,
		job.NodeGroup,
		job.ID)
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.NotFoundf("job id: %d not found", job.ID)
	}
	return nil
}
