package data

import (
	"context"
	"database/sql"
	goErrors "errors"
	"time"

	"gitlab.unanet.io/devops/go/pkg/errors"
)

type Cluster struct {
	ID            string       `db:"id"`
	Name          string       `db:"name"`
	ProviderGroup string       `db:"provider_group"`
	SchQueueUrl   string       `db:"sch_queue_url"`
	CreatedAt     sql.NullTime `db:"created_at"`
	UpdatedAt     sql.NullTime `db:"updated_at"`
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

func (r *Repo) CreateCluster(ctx context.Context, model *Cluster) error {
	now := time.Now().UTC()
	model.CreatedAt = sql.NullTime{
		Time:  now,
		Valid: true,
	}

	err := r.db.QueryRowxContext(ctx, `
	INSERT INTO cluster(id, name, provider_group, sch_queue_url, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`, model.ID, model.Name, model.ProviderGroup, model.SchQueueUrl, model.CreatedAt).
		StructScan(model)

	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (r *Repo) Clusters(ctx context.Context) ([]Cluster, error) {
	rows, err := r.db.QueryxContext(ctx, `
		select
			id,
			name,
			provider_group,
			sch_queue_url,
			created_at,
			updated_at
		from cluster`)

	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var cc []Cluster
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}

		var c Cluster
		err = rows.StructScan(&c)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		cc = append(cc, c)
	}

	return cc, nil
}

func (r *Repo) UpdateCluster(ctx context.Context, model *Cluster) error {
	now := time.Now().UTC()
	model.UpdatedAt = sql.NullTime{
		Time:  now,
		Valid: true,
	}

	result, err := r.db.ExecContext(ctx, `
		update cluster set 
			name = $2,
		    provider_group = $3,
			sch_queue_url = $4,
			updated_at = $5
		where id = $1
		RETURNING created_at
	`,
		model.ID,
		model.Name,
		model.ProviderGroup,
		model.SchQueueUrl,
		model.UpdatedAt)
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.NotFoundf("cluster id: %s not found", model.ID)
	}
	return nil
}

func (r *Repo) DeleteCluster(ctx context.Context, id int) error {
	return r.deleteByID(ctx, "cluster", id)
}
