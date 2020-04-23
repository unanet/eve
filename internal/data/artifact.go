package data

import (
	"context"

	"gitlab.unanet.io/devops/eve/pkg/errors"
)

type Artifact struct {
	ID               string   `db:"id"`
	Name             string   `db:"name"`
	ArtifactType     string   `db:"artifact_type"`
	ProviderGroup    string   `db:"provider_group"`
	FunctionPointer  string   `db:"function_pointer"`
	CustomerDeployed bool     `db:"customer_deployed"`
	Metadata         JSONText `db:"metadata"`
}

type Artifacts []Artifact

func (r *Repo) ArtifactByID(ctx context.Context, id int) (*Artifact, error) {
	var artifact Artifact

	row := r.db.QueryRowxContext(ctx, "select * from artifact where id = $1", id)
	err := row.StructScan(&artifact)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, NotFoundErrorf("artifact with id: %d, not found", id)
		}
		return nil, errors.Wrap(err)
	}

	return &artifact, nil
}

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

func (r *Repo) Artifacts(ctx context.Context) (Artifacts, error) {
	return r.artifacts(ctx)
}

func (r *Repo) artifacts(ctx context.Context, whereArgs ...WhereArg) (Artifacts, error) {
	sql, args := CheckWhereArgs("select * from artifact", whereArgs)
	rows, err := r.db.QueryxContext(ctx, sql, args...)
	if err != nil {
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
