package data

import (
	"context"
	"fmt"

	"gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/json"
)

type RequestArtifact struct {
	ArtifactID       int       `db:"artifact_id"`
	ArtifactName     string    `db:"artifact_name"`
	ProviderGroup    string    `db:"provider_group"`
	FeedName         string    `db:"feed_name"`
	ArtifactMetadata json.Text `db:"artifact_metadata"`
	ServerMetadata   json.Text `db:"server_metadata"`
	RequestedVersion string    `db:"requested_version"`
}

func (ra *RequestArtifact) Path() string {
	return fmt.Sprintf("%s/%s", ra.ProviderGroup, ra.ArtifactName)
}

type RequestArtifacts []RequestArtifact

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
