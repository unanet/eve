package data

import (
	"context"
	"database/sql"

	"gitlab.unanet.io/devops/eve/pkg/errors"
)

type RequestArtifact struct {
	ArtifactID       int      `db:"artifact_id"`
	ArtifactName     string   `db:"artifact_name"`
	ProviderGroup    string   `db:"provider_group"`
	FeedName         string   `db:"feed_name"`
	ArtifactMetadata JSONText `db:"artifact_metadata"`
	ServerMetadata   JSONText `db:"server_metadata"`
	RequestedVersion string   `db:"requested_version"`
}

type RequestArtifacts []RequestArtifact

type DeployedArtifact struct {
	ID               int            `db:"id"`
	NamespaceID      int            `db:"namespace_id"`
	NamespaceName    string         `db:"namespace_name"`
	ArtifactID       int            `db:"artifact_id"`
	ArtifactName     string         `db:"artifact_name"`
	RequestedVersion string         `db:"requested_version"`
	DeployedVersion  sql.NullString `db:"deployed_version"`
	Metadata         JSONText       `db:"metadata"`
	CustomerName     sql.NullString `db:"customer_name"`
}

type DeployedArtifacts []DeployedArtifact

func (r *Repo) RequestArtifactByEnvironment(ctx context.Context, artifactName string, environmentID int) (*RequestArtifact, error) {
	var requestedArtifact RequestArtifact

	row := r.db.QueryRowxContext(ctx, `
		select a.id as artifact_id,
		       a.name as artifact_name,
		       a.provider_group as provider_group,
		       f.name as feed_name
		from artifact as a
		    left join environment e on e.id = $1
		    left join environment_feed_map efm on e.id = efm.environment_id
			left join feed f on efm.feed_id = f.id
		where a.name = $2
	`, environmentID, artifactName)

	err := row.StructScan(&requestedArtifact)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, NotFoundErrorf("artifact with name: %s not found", artifactName)
		}
		return nil, errors.Wrap(err)
	}

	return &requestedArtifact, nil
}
