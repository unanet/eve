package data

import (
	"context"
	"database/sql"
	goErrors "errors"
	"fmt"

	"gitlab.unanet.io/devops/go/pkg/errors"
)

type RequestArtifact struct {
	ArtifactID       int            `db:"artifact_id"`
	ArtifactName     string         `db:"artifact_name"`
	ProviderGroup    string         `db:"provider_group"`
	FeedName         string         `db:"feed_name"`
	FeedType         string         `db:"feed_type"`
	FunctionPointer  sql.NullString `db:"function_pointer"`
	RequestedVersion string         `db:"requested_version"`
}

func (ra *RequestArtifact) Path() string {
	return fmt.Sprintf("%s/%s", ra.ProviderGroup, ra.ArtifactName)
}

type RequestArtifacts []RequestArtifact

func (r *Repo) RequestServiceArtifactByEnvironment(ctx context.Context, serviceName string, environmentID int) (*RequestArtifact, error) {
	var requestedArtifact RequestArtifact

	row := r.db.QueryRowxContext(ctx, `
		select a.id as artifact_id,
		       a.name as artifact_name,
		       a.function_pointer as function_pointer,
		       a.feed_type as feed_type,
		       a.provider_group as provider_group,
		       f.name as feed_name
		from service as s
		    left join artifact as a on s.artifact_id = a.id
		    left join environment e on e.id = $1
		    left join environment_feed_map efm on e.id = efm.environment_id
			left join feed f on efm.feed_id = f.id and f.feed_type = a.feed_type
		where f.name is not null and s.name = $2
	`, environmentID, serviceName)

	err := row.StructScan(&requestedArtifact)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return nil, NotFoundErrorf("service with name: %s not found", serviceName)
		}
		return nil, errors.Wrap(err)
	}

	return &requestedArtifact, nil
}

func (r *Repo) RequestJobArtifactByEnvironment(ctx context.Context, jobName string, environmentID int) (*RequestArtifact, error) {
	var requestedArtifact RequestArtifact

	row := r.db.QueryRowxContext(ctx, `
		select a.id as artifact_id,
		       a.name as artifact_name,
		       a.function_pointer as function_pointer,
		       a.feed_type as feed_type,
		       a.provider_group as provider_group,
		       f.name as feed_name
		from job as j
		    left join artifact as a on j.artifact_id = a.id
		    left join environment e on e.id = $1
		    left join environment_feed_map efm on e.id = efm.environment_id
			left join feed f on efm.feed_id = f.id and f.feed_type = a.feed_type
		where f.name is not null and j.name = $2
	`, environmentID, jobName)

	err := row.StructScan(&requestedArtifact)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return nil, NotFoundErrorf("job with name: %s not found", jobName)
		}
		return nil, errors.Wrap(err)
	}

	return &requestedArtifact, nil
}

func (r *Repo) RequestDatabaseArtifactByEnvironment(ctx context.Context, databaseName string, environmentID int) (*RequestArtifact, error) {
	var requestedArtifact RequestArtifact

	row := r.db.QueryRowxContext(ctx, `
		select a.id as artifact_id,
		       a.name as artifact_name,
		       a.function_pointer as function_pointer,
		       a.feed_type as feed_type,
		       a.provider_group as provider_group,
		       f.name as feed_name
		from eve.public.database_instance as di
		    left join database_type as dt on di.database_type_id = dt.id
		    left join artifact as a on dt.migration_artifact_id = a.id
		    left join environment e on e.id = $1
		    left join environment_feed_map efm on e.id = efm.environment_id
			left join feed f on efm.feed_id = f.id and f.feed_type = a.feed_type
		where f.name is not null and di.name = $2
	`, environmentID, databaseName)

	err := row.StructScan(&requestedArtifact)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return nil, NotFoundErrorf("database instance with name: %s not found", databaseName)
		}
		return nil, errors.Wrap(err)
	}

	return &requestedArtifact, nil
}
