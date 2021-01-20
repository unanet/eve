package data

import (
	"context"
	"database/sql"
	goErrors "errors"
	"time"

	"gitlab.unanet.io/devops/go/pkg/errors"
)

type Cluster struct {
	ID            string     `db:"id"`
	Name          string     `db:"name"`
	ProviderGroup string     `db:"provider_group"`
	SchQueueUrl   string     `db:"sch_queue_url"`
	CreatedAt     *time.Time `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
}

type Clusters []Cluster

func (r *Repo) ClusterByID(ctx context.Context, id int) (*Cluster, error) {
	var cluster Cluster

	row := r.db.QueryRowxContext(ctx, "select * from cluster where id = $1", id)
	err := row.StructScan(&cluster)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return nil, NotFoundErrorf("cluster with id: %d, not found", id)
		}
		return nil, errors.Wrap(err)
	}

	return &cluster, nil
}

func (r *Repo) ClustersByProvider(ctx context.Context, provider string) (Artifacts, error) {

	rows, err := r.db.QueryxContext(ctx, `
		select c.id,
		       c.name,
		       c.sch_queue_url,
		       c.provider_group,
		       c.created_at,
		       c.updated_at
		       from cluster c where provider_group = $1`, provider)

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
