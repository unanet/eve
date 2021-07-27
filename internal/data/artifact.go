package data

import (
	"context"
	"database/sql"
	goErrors "errors"
	"github.com/unanet/go/pkg/errors"
)

type Artifact struct {
	ID            int    `db:"id"`
	Name          string `db:"name"`
	FeedType      string `db:"feed_type"`
	ProviderGroup string `db:"provider_group"`
	ImageTag      string `db:"image_tag"`
	ServicePort   int    `db:"service_port"`
	MetricsPort   int    `db:"metrics_port"`
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
		       a.image_tag,
		       a.service_port,
		       a.metrics_port
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

func (r *Repo) Artifact(ctx context.Context) ([]Artifact, error) {
	rows, err := r.db.QueryxContext(ctx, `
		select
			id,
			name,
			feed_type,
			provider_group,
			image_tag,
			service_port,
			metrics_port
		from artifact`)

	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var aa []Artifact
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}

		var a Artifact
		err = rows.StructScan(&a)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		aa = append(aa, a)
	}

	return aa, nil
}

func (r *Repo) CreateArtifact(ctx context.Context, art *Artifact) error {

	err := r.db.QueryRowxContext(ctx, `
	INSERT INTO artifact (
		 id, 
		 name, 
		 feed_type, 
		 provider_group,
		 image_tag, 
		 service_port, 
		 metrics_port)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`,
		art.ID,
		art.Name,
		art.FeedType,
		art.ProviderGroup,
		art.ImageTag,
		art.ServicePort,
		art.MetricsPort).
		StructScan(art)

	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (r *Repo) UpdateArtifact(ctx context.Context, model *Artifact) error {

	result, err := r.db.ExecContext(ctx, `
		update artifact set 
			name = $2,
			feed_type = $3,
			provider_group = $4,
			image_tag = $5,
			service_port = $6,
			metrics_port = $7
		where id = $1
	`,
		model.ID,
		model.Name,
		model.FeedType,
		model.ProviderGroup,
		model.ImageTag,
		model.ServicePort,
		model.MetricsPort)
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.NotFoundf("artifact id: %d not found", model.ID)
	}
	return nil
}

func (r *Repo) DeleteArtifact(ctx context.Context, id int) error {
	return r.deleteByID(ctx, "artifact", id)
}
