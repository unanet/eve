package data

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/json"
)

type Job struct {
	JobID            int            `db:"job_id"`
	ArtifactID       int            `db:"artifact_id"`
	ArtifactName     string         `db:"artifact_name"`
	RequestedVersion string         `db:"requested_version"`
	DeployedVersion  sql.NullString `db:"deployed_version"`
	ServiceAccount   string         `db:"service_account"`
	ImageTag         string         `db:"image_tag"`
	RunAs            int            `db:"run_as"`
	Metadata         json.Text      `db:"metadata"`
	JobName          string         `db:"job_name"`
}

type Jobs []Job

func (r *Repo) UpdateDeployedJobVersion(ctx context.Context, id int, version string) error {
	result, err := r.db.ExecContext(ctx, "update job set deployed_version = $1, updated_at = $2 where id = $3", version, time.Now().UTC(), id)
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.Wrapf("the following id: %d was not found to update in the job table", id)
	}
	return nil
}

func (r *Repo) DeployedJobsByNamespaceID(ctx context.Context, namespaceID int) (Jobs, error) {
	rows, err := r.db.QueryxContext(ctx, `
		select 
			j.id as job_id,
		    a.id as artifact_id, 
			a.name as artifact_name,
		    j.deployed_version as deployed_version,
		    a.image_tag as image_tag,
		    a.service_account as service_account,
		    a.run_as as run_as,
		    j.name as job_name,
			jsonb_merge(e.metadata, jsonb_merge(a.metadata, jsonb_merge(ns.metadata, j.metadata))) as metadata,
			COALESCE(j.override_version, ns.requested_version) as requested_version
		from job as j 
		    left join artifact a on j.artifact_id = a.id
		    left join namespace ns on j.namespace_id = ns.id
			left join environment e on ns.environment_id = e.id
		where a.id is not null and ns.id = $1
	`, namespaceID)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	var jobs Jobs
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

func (r *Repo) JobArtifacts(ctx context.Context, namespaceIDs []int) (RequestArtifacts, error) {
	esql, args, err := sqlx.In(`
			select distinct
			 	a.id as artifact_id,
				a.name as artifact_name,
				a.provider_group as provider_group,
				a.function_pointer as function_pointer,
				a.feed_type as feed_type,
				f.name as feed_name,
				COALESCE(j.override_version, ns.requested_version) as requested_version
			from job as j
			    left join namespace ns on j.namespace_id = ns.id
				left join artifact a on j.artifact_id = a.id
			    left join environment e on ns.environment_id = e.id
				left join environment_feed_map efm on e.id = efm.environment_id
				left join feed f on efm.feed_id = f.id and f.feed_type = a.feed_type
			where f.name is not null and ns.id in (?)`, namespaceIDs)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	esql = r.db.Rebind(esql)
	rows, err := r.db.QueryxContext(ctx, esql, args...)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	var artifacts RequestArtifacts
	for rows.Next() {
		var artifact RequestArtifact
		err = rows.StructScan(&artifact)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		artifacts = append(artifacts, artifact)
	}
	return artifacts, nil
}
