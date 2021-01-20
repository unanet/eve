package data

import (
	"context"
	"database/sql"
	goErrors "errors"

	"gitlab.unanet.io/devops/go/pkg/errors"
	"gitlab.unanet.io/devops/go/pkg/json"
)

type Artifact struct {
	ID              int            `db:"id"`
	Name            string         `db:"name"`
	FeedType        string         `db:"feed_type"`
	ProviderGroup   string         `db:"provider_group"`
	FunctionPointer sql.NullString `db:"function_pointer"`
	ImageTag        string         `db:"image_tag"`
	ServicePort     int            `db:"service_port"`
	MetricsPort     int            `db:"metrics_port"`
	ServiceAccount  string         `db:"service_account"`
	RunAs           int            `db:"run_as"`
	LivelinessProbe json.Object    `db:"liveliness_probe"`
	ReadinessProbe  json.Object    `db:"readiness_probe"`
}

type Artifacts []Artifact

func (r *Repo) ArtifactByName(ctx context.Context, name string) (*Artifact, error) {
	var artifact Artifact

	row := r.db.QueryRowxContext(ctx, "select * from artifact where name = $1", name)
	err := row.StructScan(&artifact)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return nil, NotFoundErrorf("artifact with name: %s, not found", name)
		}
		return nil, errors.Wrap(err)
	}

	return &artifact, nil
}

func (r *Repo) ArtifactByID(ctx context.Context, id int) (*Artifact, error) {
	var artifact Artifact

	row := r.db.QueryRowxContext(ctx, "select * from artifact where id = $1", id)
	err := row.StructScan(&artifact)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return nil, NotFoundErrorf("artifact with id: %v, not found", id)
		}
		return nil, errors.Wrap(err)
	}

	return &artifact, nil
}

func (r *Repo) ArtifactsByProvider(ctx context.Context, provider string) (Artifacts, error) {

	rows, err := r.db.QueryxContext(ctx, `
		select a.id,
		       a.name,
		       a.feed_type,
		       a.provider_group,
		       a.function_pointer,
		       a.metadata,
		       a.image_tag,
		       a.service_port,
		       a.metrics_port,
		       a.service_account,
		       a.run_as,
		       a.liveliness_probe,
		       a.readiness_probe
		       from artifact a where provider_group = $1`, provider)

	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return nil, NotFoundErrorf("artifacts with provider: %v, not found", provider)
		}
		return nil, errors.Wrap(err)
	}
	var artifacts []Artifact

	for rows.Next() {
		var artifact Artifact
		err = rows.StructScan(&artifact)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		artifacts = append(artifacts, artifact)
	}

	return artifacts, nil
}
