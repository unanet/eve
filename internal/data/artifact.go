package data

import (
	"context"
	"database/sql"

	"gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/json"
)

type Artifact struct {
	ID              int            `db:"id"`
	Name            string         `db:"name"`
	FeedType        string         `db:"feed_type"`
	ProviderGroup   string         `db:"provider_group"`
	FunctionPointer sql.NullString `db:"function_pointer"`
	Metadata        json.Text      `db:"metadata"`
	ImageTag        string         `db:"image_tag"`
	ServicePort     int            `db:"service_port"`
	MetricsPort     int            `db:"metrics_port"`
	ServiceAccount  string         `db:"service_account"`
	RunAs           int            `db:"run_as"`
	LivelinessProbe json.Text      `db:"liveliness_probe"`
	ReadinessProbe  json.Text      `db:"readiness_probe"`
}

type Artifacts []Artifact

func (r *Repo) ArtifactByName(ctx context.Context, name string) (*Artifact, error) {
	var artifact Artifact

	row := r.db.QueryRowxContext(ctx, "select * from artifact where name = $1", name)
	err := row.StructScan(&artifact)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, NotFoundErrorf("artifact with name: %s, not found", name)
		}
		return nil, errors.Wrap(err)
	}

	return &artifact, nil
}
