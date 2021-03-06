package data

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/unanet/go/pkg/errors"
)

type RequestArtifact struct {
	ArtifactID       int    `db:"artifact_id"`
	ArtifactName     string `db:"artifact_name"`
	ProviderGroup    string `db:"provider_group"`
	FeedName         string `db:"feed_name"`
	FeedType         string `db:"feed_type"`
	RequestedVersion string `db:"requested_version"`
}

func (ra *RequestArtifact) Path() string {
	return fmt.Sprintf("%s/%s", ra.ProviderGroup, ra.ArtifactName)
}

type RequestArtifacts []RequestArtifact

func (r *Repo) RequestServiceArtifactByEnvironment(ctx context.Context, serviceName, artifactName string, environmentID int, ns []int) (RequestArtifacts, error) {
	esql, args, err := sqlx.In(`
		select distinct a.id as artifact_id,
		       			a.name as artifact_name,
		       			a.feed_type as feed_type,
		       			a.provider_group as provider_group,
		       			f.name as feed_name,
		       			CASE WHEN ? THEN COALESCE(s.override_version, ns.requested_version)
		       			     ELSE ''
						END as requested_version
		from service as s
		    left join namespace as ns on s.namespace_id = ns.id
		    left join artifact as a on s.artifact_id = a.id
		    left join environment e on e.id = ?
		    left join environment_feed_map efm on e.id = efm.environment_id
			left join feed f on efm.feed_id = f.id and f.feed_type = a.feed_type
		where f.name is not null and (a.name = ? or s.name = ?) and s.namespace_id in (?)
	`, len(ns) == 1, environmentID, artifactName, serviceName, ns)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	esql = r.db.Rebind(esql)
	rows, err := r.db.QueryxContext(ctx, esql, args...)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	var services []RequestArtifact
	for rows.Next() {
		var service RequestArtifact
		err = rows.StructScan(&service)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		services = append(services, service)
	}

	if len(services) == 0 {
		return nil, NotFoundErrorf("service with name: %s not found", serviceName)
	}
	return services, nil
}

func (r *Repo) RequestJobArtifactByEnvironment(ctx context.Context, jobName, artifactName string, environmentID int, ns []int) (RequestArtifacts, error) {
	esql, args, err := sqlx.In(`
		select distinct a.id as artifact_id,
		       			a.name as artifact_name,
		       			a.feed_type as feed_type,
		       			a.provider_group as provider_group,
		       			f.name as feed_name,
		       			CASE WHEN ? THEN COALESCE(j.override_version, ns.requested_version)
		       			     ELSE ''
						END as requested_version
		from job as j
		    left join namespace as ns on j.namespace_id = ns.id
		    left join artifact as a on j.artifact_id = a.id
		    left join environment e on e.id = ?
		    left join environment_feed_map efm on e.id = efm.environment_id
			left join feed f on efm.feed_id = f.id and f.feed_type = a.feed_type
		where f.name is not null and (a.name = ? or j.name = ?) and j.namespace_id in (?)
	`, len(ns) == 1, environmentID, artifactName, jobName, ns)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	esql = r.db.Rebind(esql)
	rows, err := r.db.QueryxContext(ctx, esql, args...)
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

	if len(jobs) == 0 {
		return nil, NotFoundErrorf("job with name: %s not found", jobName)
	}
	return jobs, nil
}
